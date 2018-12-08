package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	url = "http://modfin.se/api/puzzles/index-trader.json"
)

func main() {
	data, err := getData(url)
	if err != nil {
		log.Fatalln("cannot parse data ", err)
	}

	d, err := processData(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Best deal:", d)
}

func processData(data []StockData) (d Deal, err error) {
	dataLen := len(data)

	if dataLen < 2 {
		err = errors.New("Not enough data")
		return
	}

	d = Deal{}
	d.Buy(data[dataLen-1])
	d.Sell(data[dataLen-2])

	tmpDeal := d

	for i := len(data) - 3; i >= 0; i-- {
		if data[i+1].Low < tmpDeal.Open.Price {
			if tmpDeal.Profit() > d.Profit() {
				d = tmpDeal
			}
			tmpDeal.Buy(data[i+1])
			tmpDeal.Sell(data[i])
		}
		if data[i].High > tmpDeal.Close.Price {
			tmpDeal.Sell(data[i])
		}
	}

	if tmpDeal.Profit() > d.Profit() {
		d = tmpDeal
	}

	return
}

func getData(url string) (data []StockData, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	p := &PuzzleData{}
	err = json.Unmarshal(raw, p)
	if err != nil {
		return
	}

	return p.Data, nil
}

type PuzzleData struct {
	Puzzle     string      `json:"puzzle"`
	Info       string      `json:"info"`
	Submission string      `json:"submission"`
	Data       []StockData `json:"data"`
}

type StockData struct {
	QuoteDate int32   `json:"quote_date"`
	Paper     string  `json:"paper"`
	Exchange  string  `json:"exch"`
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    float64 `json:"volume"`
	Value     float64 `json:"value"`
}

type Order struct {
	Price float64
	Date  int32
}

func (d *Deal) Buy(s StockData) {
	d.Open.Price = s.Low
	d.Open.Date = s.QuoteDate
}

func (d *Deal) Sell(s StockData) {
	d.Close.Price = s.High
	d.Close.Date = s.QuoteDate
}

type Deal struct {
	Open  Order
	Close Order
}

func (d Deal) String() string {
	return fmt.Sprintf(
		"buy %d at %f, sell %d at %f, profit: %.2f%%",
		d.Open.Date,
		d.Open.Price,
		d.Close.Date,
		d.Close.Price,
		(d.Close.Price/d.Open.Price-1)*100,
	)
}

func (d Deal) Profit() float64 {
	return d.Close.Price - d.Open.Price
}
