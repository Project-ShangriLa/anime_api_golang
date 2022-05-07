package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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

func statusByCoursId(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[statusByCoursId OK]\n"))
}

func historyDaily(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("[historyDaily OK]\n"))
}
