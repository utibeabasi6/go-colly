package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

type Product struct {
	Name     string
	Image    string
	Price    string
	Url      string
	Discount string
}

func main() {
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	products := make([]Product, 0)
	// Find and visit all links
	c.OnHTML("a.core", func(e *colly.HTMLElement) {
		e.ForEach("div.name", func(i int, h *colly.HTMLElement) {
			item := Product{}
			item.Name = h.Text
			item.Image = e.ChildAttr("img", "data-src")
			item.Price = e.Attr("data-price")
			item.Url = "https://jumia.com.ng" + e.Attr("href")
			item.Discount = e.ChildText("div.tag._dsct")
			products = append(products, item)
		})

	})

	c.OnError(func(r *colly.Response, e error) {
		log.Fatalln(fmt.Sprintf("An error occured: %v", e))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(products, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("products.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	c.Visit("https://jumia.com.ng/")
}
