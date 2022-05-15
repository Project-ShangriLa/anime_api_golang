package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Project-ShangriLa/anime_api_golang/model"
)

// go build -o sana_client
// ./sana_client -k xxxxx -b 1506 -s 20220401 -e 20220506
// ./sana_client -d http://localhost:8080 -k aiueo -b 1506 -s 20220501 -e 20220506
// ./sana_client -d http://localhost:8080 -k aiueo -a kimetsu_off -s 20220501 -e 20220506
// ./sana_client -k xxxxx -a paripikoumei_PR -s 20220401 -e 20220514 -o ~/tmp/paripi.csv
func main() {
	// curl --header 'X-CLI-API-KEY:aiueo' "http://localhost:8080/anime/v1/twitter/follower/history/daily?account=paripikoumei_PR&startdate=20220501&enddate=20220506"

	var domain string
	var account string
	var baseId string
	var startdate string
	var enddate string
	var clientApiKey string
	var outFileName string
	flag.StringVar(&domain, "d", "https://api.moemoe.tokyo", "Anime API protocol and domain")
	flag.StringVar(&account, "a", "", "Twitter Account")
	flag.StringVar(&baseId, "b", "", "Base Id")
	flag.StringVar(&startdate, "s", "", "start date")
	flag.StringVar(&enddate, "e", "", "end date")
	flag.StringVar(&clientApiKey, "k", "", "Cliant API Key")
	flag.StringVar(&outFileName, "o", "", "output File Name")
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
		fmt.Println(r)
	}

	if outFileName != "" {
		createCsv(outFileName, &twhs)
	}

}

func createCsv(outFileName string, twsh *[]model.TwitterStatusHistory) {
	f, err := os.OpenFile(outFileName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	w.Write([]string{"日付", "フォロワー数"})

	for _, v := range *twsh {
		date := v.GetDate.Format("2006/01/02")
		follower := fmt.Sprintf("%d", v.Follower)

		if err := w.Write([]string{date, follower}); err != nil {
			log.Fatal(err)
		}
	}
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
