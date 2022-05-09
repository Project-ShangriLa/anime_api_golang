package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Project-ShangriLa/anime_api_golang/model"
)

// go build -o sana_client
// ./sana_client -k xxxxx -b 1506 -s 20220401 -e 20220506
// ./sana_client -d http://localhost:8080 -k aiueo -b 1506 -s 20220501 -e 20220506
// ./sana_client -d http://localhost:8080 -k aiueo -a kimetsu_off -s 20220501 -e 20220506
func main() {
	// curl --header 'X-CLI-API-KEY:aiueo' "http://localhost:8080/anime/v1/twitter/follower/history/daily?account=paripikoumei_PR&startdate=20220501&enddate=20220506"

	var domain string
	var account string
	var baseId string
	var startdate string
	var enddate string
	var clientApiKey string
	flag.StringVar(&domain, "d", "https://api.moemoe.tokyo", "Anime API protocol and domain")
	flag.StringVar(&account, "a", "", "Twitter Account")
	flag.StringVar(&baseId, "b", "", "Base Id")
	flag.StringVar(&startdate, "s", "", "start date")
	flag.StringVar(&enddate, "e", "", "end date")
	flag.StringVar(&clientApiKey, "k", "", "Cliant API Key")
	flag.Parse()

	target := "account=" + account
	// -bオプションがあればそちらを優先
	if baseId != "" {
		target = "baseid=" + baseId
	}

	url := fmt.Sprintf("%s/anime/v1/twitter/follower/history/daily?%s&startdate=%s&enddate=%s",
		domain, target, startdate, enddate)
	//println(url)

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		println(err)
	}
	req.Header.Add("X-CLI-API-KEY", clientApiKey)
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		print(err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	twhs := []model.TwitterStatusHistory{}

	json.Unmarshal(body, &twhs)

	//fmt.Println(string(body))

	for _, v := range twhs {
		r := fmt.Sprintf("%s,%d", v.GetDate.Format("2006/01/02"), v.Follower)
		println(r)
	}

}
