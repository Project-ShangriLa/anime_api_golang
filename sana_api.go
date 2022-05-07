package main

import (
	"log"
	"net/http"
	"os"

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
	//days := r.FormValue("days")

	log.Print(account)

	db := gormConnect()
	defer db.Close()

	var base model.Basis
	db.Where("twitter_account = ?", account).Last(&base)

	log.Print(base.ID)

	w.Write([]byte("[historyDaily OK]\n"))
}
