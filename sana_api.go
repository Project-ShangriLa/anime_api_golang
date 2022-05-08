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
	rdays := r.FormValue("days")
	const DEFAULT_PDAYS = 30
	var pastdays int
	if rdays == "" {
		pastdays = DEFAULT_PDAYS
	} else {
		pastdays, _ = strconv.Atoi(rdays)
	}

	log.Print(account)

	db := gormConnect()
	defer db.Close()

	var base model.Basis
	db.Where("twitter_account = ?", account).Last(&base)

	log.Print(base.ID)

	var twhs = []model.TwitterStatusHistory{}

	today := time.Now()
	pastday := today.AddDate(0, 0, pastdays*-1)

	// https://gorm.io/docs/query.html
	db.Where("bases_id = ?", base.ID).Where("get_date BETWEEN ? AND ?", pastday, today).Find(&twhs)

	// TODO 日付の重複をMAPで整理
	// TODO time型の表記がおかしいのをなおす
	// TODO JSON形式だけでなくCSV形式で返す

	res, err := json.Marshal(twhs)
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
