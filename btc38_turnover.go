package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
)

func parse(body string) float64 {
	body = body[1 : len(body)-1]
	r, err := regexp.Compile(`\[.+?\]`)
	if err != nil {
		panic(err)
	}
	res := r.FindAllString(body, -1)
	sum := 0.0
	for i := 0; i < 24; i++ {
		cuts := strings.Split(res[len(res)-1-i], ",")
		var turnover, price float64
		fmt.Sscan(cuts[1], &turnover)
		fmt.Sscan(cuts[2], &price)
		sum += turnover * price
	}
	return sum
}

func getCoin(coinName string) float64 {
	url := fmt.Sprintf("http://www.btc38.com/trade/getTradeTimeLine.php?coinname=%s&n=%.20f", coinName, rand.Float64())
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return parse(string(body))
}

type Record struct {
	Name     string
	Turnover float64
}

type byTurnover []Record

func (v byTurnover) Len() int           { return len(v) }
func (v byTurnover) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v byTurnover) Less(i, j int) bool { return v[i].Turnover > v[j].Turnover }

func main() {
	names := []string{"BTC", "LTC", "DOG", "XPM", "BEC", "XRP", "ZCC", "MEC", "ANC", "PPC",
		"SRC", "TAG", "PTS", "WDC", "APC", "DGC", "UNC", "QRK"}

	var records []Record

	for _, name := range names {
		turnover := getCoin(name)
		records = append(records, Record{name, turnover})
	}

	sort.Sort(byTurnover(records))

	for _, r := range records {
		fmt.Printf("%v\t%.2f\n", r.Name, r.Turnover)
	}
}
