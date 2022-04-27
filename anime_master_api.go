package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var router = mux.NewRouter()

var cacheBases = make(map[int][]byte)
var cacheBasesWithOgp = make(map[int][]byte)

const (
	APIKEY_HEADER_NAME = "X-API-KEY"
	APIKEY_ENV_NAME    = "X_ANIME_API_KEY"
	COURSID_MIN        = 1
	COURSID_MAX        = 104 // COURID_IDの理論的最大値 　2014 + COURID_MAX/4 = 年数、2039年までリクエストを許容
)

func init() {
	http.Handle("/", router)

	// http://www.gorillatoolkit.org/pkg/mux
	router.HandleFunc("/anime/v1/master/cours", coursHandler).Methods("GET")

	router.HandleFunc("/anime/v1/master/{year_num:[0-9]{4}}", yearTitleHandler).Methods("GET")

	router.HandleFunc("/anime/v1/master/{year_num:[0-9]{4}}/{cours:[1-4]}", animeAPIReadHandler).Methods("GET")

	// キャッシュ全クリア 環境変数　認証キーあり
	router.HandleFunc("/anime/v1/master/cache/clear", cacheClear).Methods("POST")

	// キャッシュ全再取得 環境変数　認証キーあり
	router.HandleFunc("/anime/v1/master/cache/refresh", cacheRefresh).Methods("POST")

	//TODO
	// キャッシュ指定クリア 環境変数　認証キーあり
}

func gormConnect() *gorm.DB {
	var err error

	dbHost := os.Getenv("ANIME_API_DB_HOST")
	dbUser := os.Getenv("ANIME_API_DB_USER")
	dbPass := os.Getenv("ANIME_API_DB_PASS")

	if len(dbHost) == 0 {
		dbUser = "root"
	}

	if len(dbUser) > 0 {
		dbPass = ":" + dbPass
	}

	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	db, err := gorm.Open("mysql", dbUser+dbPass+"@"+"tcp("+dbHost+")/anime_admin_development?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return db
}

func main() {
	http.ListenAndServe(":8080", nil)
}

type Base struct {
	Id               int
	Title            string
	TitleShort1      sql.NullString
	TitleShort2      sql.NullString
	TitleShort3      sql.NullString
	TitleEn          sql.NullString
	PublicURL        string
	TwitterAccount   string
	TwitterHashTag   string
	CoursID          int
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Sex              sql.NullInt64
	Sequel           sql.NullInt64
	CityCode         sql.NullInt64
	CityName         sql.NullString
	ProductCompanies sql.NullString
	OgTitle          sql.NullString
	OgType           sql.NullString
	OgDescription    sql.NullString
	OgUrl            sql.NullString
	OgImage          sql.NullString
	OgSiteName       sql.NullString
}

type BaseJson struct {
	Id               int       `json:"id"`
	Title            string    `json:"title"`
	TitleShort1      string    `json:"title_short1"`
	TitleShort2      string    `json:"title_short2"`
	TitleShort3      string    `json:"title_short3"`
	TitleEn          string    `json:"title_en"`
	PublicURL        string    `json:"public_url"`
	TwitterAccount   string    `json:"twitter_account"`
	TwitterHashTag   string    `json:"twitter_hash_tag"`
	CoursID          int       `json:"cours_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Sex              int       `json:"sex"`
	Sequel           int       `json:"sequel"`
	CityCode         int       `json:"city_code"`
	CityName         string    `json:"city_name"`
	ProductCompanies string    `json:"product_companies"`
}

type OgpJson struct {
	OgTitle       string `json:"og_title"`
	OgType        string `json:"og_type"`
	OgDescription string `json:"og_description"`
	OgUrl         string `json:"og_url"`
	OgImage       string `json:"og_image"`
	OgSiteName    string `json:"og_site_name"`
}

type BaseJsonWithOgp struct {
	BaseJson
	Ogp OgpJson `json:"ogp"`
}

func animeAPIReadHandler(w http.ResponseWriter, r *http.Request) {

	coursID := year2coursID(r)

	log.Print(coursID)

	var res []byte
	var err error

	if coursID < COURSID_MIN || coursID > COURSID_MAX {
		http.Error(w, "COURSID OVER", http.StatusBadRequest)
		return
	}

	if r.FormValue("ogp") == "1" {
		if cacheBasesWithOgp[coursID] != nil {
			log.Print("Hit cache ogp")
			res = cacheBasesWithOgp[coursID]
		} else {
			res, err = selectBasesWithOgpRdb(coursID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Print("not cache. save cache ogp")
			cacheBasesWithOgp[coursID] = res
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(res)
	} else {
		if cacheBases[coursID] != nil {
			log.Print("Hit cache")
			res = cacheBases[coursID]
		} else {
			res, err = selectBasesRdb(coursID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Print("not cache. save cache")
			cacheBases[coursID] = res
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(res)
	}

}

func selectBasesRdb(coursId int) ([]byte, error) {
	db := gormConnect()
	defer db.Close()

	baseJsonList := []BaseJson{}

	rows, _ := db.Table("bases").Select(`
id,
title,
title_short1,
title_short2,
title_short3,
title_en,
public_url,
twitter_account,
twitter_hash_tag,
cours_id,
created_at,
updated_at,
sex,
sequel,
city_code,
city_name,
product_companies
`).
		Where("cours_id = ?", coursId).Rows()

	for rows.Next() {

		var bs Base
		var bsj BaseJson

		if err := rows.Scan(&bs.Id, &bs.Title, &bs.TitleShort1, &bs.TitleShort2, &bs.TitleShort3, &bs.TitleEn,
			&bs.PublicURL, &bs.TwitterAccount, &bs.TwitterHashTag, &bs.CoursID, &bs.CreatedAt, &bs.UpdatedAt,
			&bs.Sex, &bs.Sequel, &bs.CityCode, &bs.CityName, &bs.ProductCompanies); err != nil {
			log.Fatal(err)
		}

		bsj.Id = bs.Id
		bsj.Title = bs.Title
		bsj.TitleShort1 = bs.TitleShort1.String
		bsj.TitleShort2 = bs.TitleShort2.String
		bsj.TitleShort3 = bs.TitleShort3.String
		bsj.TitleEn = bs.TitleEn.String
		bsj.PublicURL = bs.PublicURL
		bsj.TwitterAccount = bs.TwitterAccount
		bsj.TwitterHashTag = bs.TwitterHashTag
		bsj.CreatedAt = bs.CreatedAt
		bsj.UpdatedAt = bs.UpdatedAt
		bsj.Sex = int(bs.Sex.Int64)
		bsj.Sequel = int(bs.Sequel.Int64)
		bsj.CityCode = int(bs.CityCode.Int64)
		bsj.CityName = bs.CityName.String
		bsj.ProductCompanies = bs.ProductCompanies.String

		baseJsonList = append(baseJsonList, bsj)
	}

	res, err := json.Marshal(baseJsonList)

	return res, err

}

func selectBasesWithOgpRdb(coursId int) ([]byte, error) {
	// 指定した条件を元に複数のレコードを引っ張ってくる
	db := gormConnect()
	defer db.Close()

	baseWithOgp := []BaseJsonWithOgp{}

	//log.Print("With OGP option")

	rows, _ := db.Table("bases").Select(`
	id,
	title,
	title_short1,
	title_short2,
	title_short3,
	title_en,
	bases.public_url as public_url,
	twitter_account,
	twitter_hash_tag,
	cours_id,
	bases.created_at as created_at,
	bases.updated_at as updated_at,
	sex,
	sequel,
	city_code,
	city_name,
	og_title, 
	og_type, 
	og_description, 
	og_url, 
	og_image,
	og_site_name,
	product_companies
	`).
		Joins("join site_meta_data on bases.id = site_meta_data.bases_id and cours_id = ?", coursId).Rows()

	for rows.Next() {

		var bs Base
		var bsj BaseJsonWithOgp
		var ogpj OgpJson

		if err := rows.Scan(&bs.Id, &bs.Title, &bs.TitleShort1, &bs.TitleShort2, &bs.TitleShort3, &bs.TitleEn,
			&bs.PublicURL, &bs.TwitterAccount, &bs.TwitterHashTag, &bs.CoursID, &bs.CreatedAt, &bs.UpdatedAt,
			&bs.Sex, &bs.Sequel, &bs.CityCode, &bs.CityName,
			&bs.OgTitle, &bs.OgType, &bs.OgDescription, &bs.OgUrl, &bs.OgImage, &bs.OgSiteName, &bs.ProductCompanies); err != nil {
			log.Fatal(err)
		}

		bsj.Id = bs.Id
		bsj.Title = bs.Title
		bsj.TitleShort1 = bs.TitleShort1.String
		bsj.TitleShort2 = bs.TitleShort2.String
		bsj.TitleShort3 = bs.TitleShort3.String
		bsj.TitleEn = bs.TitleEn.String
		bsj.PublicURL = bs.PublicURL
		bsj.TwitterAccount = bs.TwitterAccount
		bsj.TwitterHashTag = bs.TwitterHashTag
		bsj.CreatedAt = bs.CreatedAt
		bsj.UpdatedAt = bs.UpdatedAt
		bsj.Sex = int(bs.Sex.Int64)
		bsj.Sequel = int(bs.Sequel.Int64)
		bsj.CityCode = int(bs.CityCode.Int64)
		bsj.CityName = bs.CityName.String
		bsj.ProductCompanies = bs.ProductCompanies.String

		ogpj.OgTitle = bs.OgTitle.String
		ogpj.OgType = bs.OgType.String
		ogpj.OgDescription = bs.OgDescription.String
		ogpj.OgUrl = bs.OgUrl.String
		ogpj.OgSiteName = bs.OgSiteName.String
		ogpj.OgImage = bs.OgImage.String

		bsj.Ogp = ogpj

		baseWithOgp = append(baseWithOgp, bsj)
	}

	res, err := json.Marshal(baseWithOgp)

	return res, err
}

type CoursInfo struct {
	Id    int `json:"id"`
	Year  int `json:"year"`
	Cours int `json:"cours"`
}

func coursHandler(w http.ResponseWriter, r *http.Request) {
	db := gormConnect()
	defer db.Close()

	CoursInfoList := []CoursInfo{}

	coursMap := map[string]CoursInfo{}

	db.Find(&CoursInfoList)

	for _, cil := range CoursInfoList {
		coursMap[strconv.Itoa(cil.Id)] = cil
	}

	res, err := json.Marshal(coursMap)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

func yearTitleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	year, _ := strconv.Atoi(vars["year_num"])

	type animeYearTitle struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
	}

	animeYearTitleList := []animeYearTitle{}

	db := gormConnect()
	defer db.Close()

	coursIdList := []int{}

	db.Table("cours_infos").Where("year = ?", year).Pluck("id", &coursIdList)

	db.Table("bases").Select("id, title").Where("cours_id in (?)", coursIdList).Scan(&animeYearTitleList)

	res, err := json.Marshal(animeYearTitleList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(res)
}

// TODO 暫定的に手動で計算、本来は管理テーブルからcours_idを算出
func year2coursID(r *http.Request) int {
	vars := mux.Vars(r)

	year, _ := strconv.Atoi(vars["year_num"])
	cours, _ := strconv.Atoi(vars["cours"])
	coursID := (year-2014)*4 + cours

	return coursID
}

func cacheClear(w http.ResponseWriter, r *http.Request) {
	rApiKey := r.Header.Get(APIKEY_HEADER_NAME)

	if os.Getenv(APIKEY_ENV_NAME) != "" && rApiKey == os.Getenv(APIKEY_ENV_NAME) {
		cacheBases = make(map[int][]byte)
		cacheBasesWithOgp = make(map[int][]byte)
		w.Write([]byte("[OK] Clear Cache\n"))
	} else {
		w.Write([]byte("[NG] Clear Cache error\n"))
	}
}

func cacheRefresh(w http.ResponseWriter, r *http.Request) {
	rApiKey := r.Header.Get(APIKEY_HEADER_NAME)

	if os.Getenv(APIKEY_ENV_NAME) != "" && rApiKey == os.Getenv(APIKEY_ENV_NAME) {
		cacheBases = make(map[int][]byte)
		cacheBasesWithOgp = make(map[int][]byte)
		db := gormConnect()
		defer db.Close()

		CoursInfoList := []CoursInfo{}
		db.Find(&CoursInfoList)

		for _, cil := range CoursInfoList {
			res, _ := selectBasesRdb(cil.Id)
			cacheBases[cil.Id] = res
			res, _ = selectBasesWithOgpRdb(cil.Id)
			cacheBasesWithOgp[cil.Id] = res
		}

		w.Write([]byte("[OK] Refresh Cache\n"))
	} else {
		w.Write([]byte("[NG] Refresh Cache error\n"))
	}
}
