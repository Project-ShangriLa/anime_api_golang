# ShangriLa ANIME API (Golang)

ShangriLa Anime APIのGolang実装です

2020年からはこちらが稼働しています

オリジナル版とAPI仕様は[こちら](https://github.com/Project-ShangriLa/sora-playframework-scala) 

## BUILD

```
go build
```

## 設定

DBの接続を環境変数で管理しています

````
export ANIME_API_DB_HOST=
export ANIME_API_DB_USER=
export ANIME_API_DB_PASS=
````

## 実行

```
./anime_master_api
```

## 動作確認

```
curl http://localhost:8080/anime/v1/master/2021 | jq .
```

```
curl http://localhost:8080/anime/v1/master/2021/3 | jq .
```

```
curl "http://localhost:8080/anime/v1/master/2021/3?ogp=1" | jq .
```