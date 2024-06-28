package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/gocolly/colly/v2"
)

type CSVURL struct {
	URL string `csv:"URL"`
}

func parseUrl() {
	fName := "./dates/catalog_url.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write CSV header
	writer.Write([]string{"URL"})

	// Instantiate default collector
	c := colly.NewCollector()

	// Extract product details
	c.OnHTML("a.a9ki2WaUy", func(e *colly.HTMLElement) {
		writer.Write([]string{
			e.Request.AbsoluteURL(e.Attr("href")),
		})
	})

	c.Visit("https://kdvonline.ru/catalog/igrushki-456")
}

func readCSV(fName string) []CSVURL {
	file, err := os.Open(fName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dates := []CSVURL{}
	if err := gocsv.UnmarshalFile(file, &dates); err != nil {
		fmt.Println(err)
	}
	return dates
}

func main() {

	parseUrl()

	dates := readCSV("./dates/catalog_url.csv")
	for _, date := range dates {
		shortName := strings.Split(date.URL, "/")
		fName := fmt.Sprintf("./dates/kdv_%v.csv", shortName[len(shortName)-1])
		file, err := os.Create(fName)
		if err != nil {
			log.Fatalf("Cannot create file %q: %s\n", fName, err)
			return
		}
		defer file.Close()
		writer := csv.NewWriter(file)
		defer writer.Flush()
		// Write CSV header
		writer.Write([]string{"Name", "Price", "Value", "URL"})

		go func() {
			// Instantiate default collector
			c := colly.NewCollector()

			// Extract product details
			c.OnHTML("div.abU3Ps1qD div.a6auoGJSo", func(e *colly.HTMLElement) {
				writer.Write([]string{
					e.ChildText(".g6auoGJSo"),
					e.ChildText(".BJThHsRzJ"),
					e.ChildText(".mJThHsRzJ"),
					e.Request.AbsoluteURL(e.ChildAttr("a", "href")),
				})
			})

			c.Visit(fmt.Sprintf("%v?page=1&limit=100&sort=default%3Aasc", date.URL))
		}()

		log.Printf("Scraping %s %s", date.URL, time.Now())

		time.Sleep(900 * time.Millisecond)
	}

	log.Println("Scraping finished, check file")
}
