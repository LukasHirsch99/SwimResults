package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	. "swimresults-backend/database"

	"github.com/gocolly/colly"
)

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func getInfoFromMSECM(msecmLink string) {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

  c.OnHTML("div.col-sm-9.col-xs-12", func(e *colly.HTMLElement) {
    googleMapsLink := e.ChildAttr("p.text-right:nth-child(1) > a", "href")
    fmt.Println(googleMapsLink)
  })

	c.Visit(msecmLink)
}

func getMeetInfo(meetId uint) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div.col-xs-12.myresults_content", func(e *colly.HTMLElement) {
		image := "https://myresults.eu" + e.ChildAttr("img.img-responsive.center-block", "src")
		name := e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header:nth-child(1)")
		dateString := getOnlyChildText(e, "div:nth-child(4) > div")
		deadlineString := getOnlyChildText(e, "div:nth-child(5) > div")
		msecmLink := e.ChildAttr("div:nth-child(14) > div > a", "href")
		if msecmLink != "" {
			getInfoFromMSECM(msecmLink)
		}
		fmt.Println(name, image, dateString, deadlineString, msecmLink)
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Overview")
}

func insertMeets() error {
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div.col-xs-12.col-md-12.myresults_content_divtable", func(e *colly.HTMLElement) {
		e.ForEach(".myresults_content_divtablerow", func(i int, row *colly.HTMLElement) {
			country := row.ChildAttr("div.col-xs-1.text-right.myresults_content_divtable_right.myresults_padding_top_5 > img", "src")
			// Insert only meets which are in austria
			if country == "/images/flags/at.png" && i == 1 {
				r := regexp.MustCompile("\\d+")
				meetId := StringToUint(r.FindString(row.ChildAttr("a", "href")))
				getMeetInfo(meetId)
			}
		})
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Today-Upcoming")
	return nil
}

func main() {
	insertMeets()
}
