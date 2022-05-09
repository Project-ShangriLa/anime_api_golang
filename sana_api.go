package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Project-ShangriLa/anime_api_golang/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// リファレンス実装
// https://github.com/Project-ShangriLa/sana_server

func middlewareCliantAuthAPI(next http.Handler) http.Handler {
	// クライアントAPIキーは発行形式でPublicに公開する想定（一般のAPIKEYなどのように）
	const APIKEY_HEADER_NAME = "X-CLI-API-KEY"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rApiKey := r.Header.Get(APIKEY_HEADER_NAME)

		if cApiKey != "" && rApiKey == cApiKey {
			next.ServeHTTP(w, r)
		} else {
			//nolint:errcheck
			w.Write([]byte("authentication error\n"))
		}
	})
}

func gormConnectV2() *gorm.DB {
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

	db, err := gorm.Open(mysql.Open(dbUser + dbPass + "@" + "tcp(" + dbHost + ")/anime_admin_development?parseTime=true"))
	if err != nil {
		panic(err.Error())
	}

	return db
}

func statusByCoursId(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[statusByCoursId OK]\n"))
}

func historyDaily(w http.ResponseWriter, r *http.Request) {
	account := r.FormValue("account")
	rBaseId := r.FormValue("baseid")
	rdays := r.FormValue("days")
	rStartDate := r.FormValue("startdate")
	rEndDate := r.FormValue("enddate")
	const DEFAULT_PDAYS = 30
	var pastdays int
	if rdays == "" {
		pastdays = DEFAULT_PDAYS
	} else {
		pastdays, _ = strconv.Atoi(rdays)
	}
	enddate := time.Now()
	startdate := enddate.AddDate(0, 0, pastdays*-1)

	// startdateパラメタ enddateパラメタ があればそちらを優先させる
	if rStartDate != "" && rEndDate != "" {
		startdate, _ = time.Parse("20060102", rStartDate)
		enddate, _ = time.Parse("20060102", rEndDate)
		// 指定された日の00:00になるので+1する
		enddate = enddate.AddDate(0, 0, 1)
	}

	db := gormConnect()
	defer db.Close()

	var base model.Basis
	var baseId int32

	// baseidが指定されたらそちらを優先する
	if rBaseId != "" {
		i, _ := strconv.Atoi(rBaseId)
		baseId = int32(i)
	} else {
		db.Where("twitter_account = ?", account).Last(&base)
		baseId = base.ID
	}

	log.Print(baseId, account, startdate, enddate)

	twhs := []model.TwitterStatusHistory{}
	resTwhs := []model.TwitterStatusHistory{}

	// https://gorm.io/docs/query.html
	db.Where("bases_id = ?", baseId).Where("get_date BETWEEN ? AND ?", startdate, enddate).Find(&twhs)

	// TODO 日付の重複をMAPで整理

	/*
		 js date parse
		 d = new Date("2022-05-07T08:01:49Z")
		 d.getFullYear() + "/" + (d.getMonth()+1) + "/" + d.getDay()
		->'2022/5/6'
	*/
	var tmpDay string
	for _, v := range twhs {
		getDateSt := v.GetDate.Format("2006-01-02")
		//log.Println(getDateSt)

		if tmpDay != getDateSt {
			resTwhs = append(resTwhs, v)
			tmpDay = getDateSt
			//log.Println(getDateSt)
		}
	}

	res, err := json.Marshal(resTwhs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, error := w.Write(res)
	if error != nil {
		log.Println(error)
	}
}
