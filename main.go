package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func parseUrl(urls []string) []string {
	urls = []string{}
	// Instantiate default collector
	c := colly.NewCollector()

	// Extract product details
	c.OnHTML("a.a9ki2WaUy", func(e *colly.HTMLElement) {
		urls = append(urls, e.Request.AbsoluteURL(e.Attr("href")))
	})

	c.Visit("https://kdvonline.ru/catalog/igrushki-456")

	return urls
}

func main() {

	urls := []string{}

	urls = parseUrl(urls)

	var wg sync.WaitGroup

	for _, date := range urls {
		wg.Add(1)

		shortName := strings.Split(date, "/")
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
			defer wg.Done()
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

			c.Visit(fmt.Sprintf("%v?page=1&limit=100&sort=default%3Aasc", date))
		}()

		time.Sleep(700 * time.Millisecond)

		log.Printf("Scraping %s", date)

		wg.Wait()
	}

	log.Println("Scraping finished, check file")
}
