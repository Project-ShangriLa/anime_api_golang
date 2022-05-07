# ShangriLa ANIME API (Golang)

ShangriLa Anime APIのGolang実装です

2020年からはこちらが稼働しています

オリジナル版とAPI仕様は[こちら](https://github.com/Project-ShangriLa/sora-playframework-scala) 

## BUILD

```
go build -trimpath -ldflags '-s -w' -o anime_api_server
```

or 

```
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags '-s -w' -o anime_api_server_linux
```


## DB設定

MySQLに下記のDDLを実行してください

https://github.com/Project-ShangriLa/anime_master_db_ddl


## 設定

DBの接続を環境変数で管理しています

````
export ANIME_API_DB_HOST=
export ANIME_API_DB_USER=
export ANIME_API_DB_PASS=
````

管理APIのKEYを環境変数で管理しています

```
export X_ANIME_API_KEY=abcde
```

SANA APIのKEYを環境変数で管理しています(初期検証用)

```
export X_ANIME_CLI_API_KEY=aiueo
```

## 実行

```
./anime_api
```

## 動作確認 Master API

### COURS

```
curl http://localhost:8080/anime/v1/master/cours | jq .
```

### YEAR
```
curl http://localhost:8080/anime/v1/master/2021 | jq .
```

### YEAR-COURS
```
curl http://localhost:8080/anime/v1/master/2021/3 | jq .
```

```
curl "http://localhost:8080/anime/v1/master/2021/3?ogp=1" | jq .
```

### キャッシュクリア

```
curl -XPOST --header 'X-API-KEY:abcde' http://localhost:8080/anime/v1/master/cache/clear
```

### キャッシュ全更新(初期化)

```
curl -XPOST --header 'X-API-KEY:abcde' http://localhost:8080/anime/v1/master/cache/refresh
```

## 動作確認 Sana API

### status by coursId

```
curl --header 'X-CLI-API-KEY:aiueo' http://localhost:8080/anime/v1/twitter/follower/status/bycours
```

### history daily

```
curl --header 'X-CLI-API-KEY:aiueo' http://localhost:8080/anime/v1/twitter/follower/history/daily
```