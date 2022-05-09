package main

import (
	"flag"
	"fmt"
)

func main() {
	// curl --header 'X-CLI-API-KEY:aiueo' "http://localhost:8080/anime/v1/twitter/follower/history/daily?account=paripikoumei_PR&startdate=20220501&enddate=20220506"
	// curl --header 'X-CLI-API-KEY:xxxxx' "https://api.moemoe.tokyo/anime/v1/twitter/follower/history/daily?account=paripikoumei_PR&startdate=20220501&enddate=20220506" | jq .

	var domain string
	var twitterAccount string
	var startdate string
	var enddate string
	var clientApiKey string
	flag.StringVar(&domain, "d", "api.moemoe.tokyo", "Anime API domain")
	flag.StringVar(&twitterAccount, "a", "", "twitter account")
	flag.StringVar(&startdate, "s", "", "start date")
	flag.StringVar(&enddate, "e", "", "end date")
	flag.StringVar(&clientApiKey, "k", "", "Cliant API Key")
	flag.Parse()

	url := fmt.Sprintf("https://%s/v1/twitter/follower/history/daily?account=%s&startdate=%s&enddate=%s",
		domain, twitterAccount, startdate, enddate)
	println(url)
}
