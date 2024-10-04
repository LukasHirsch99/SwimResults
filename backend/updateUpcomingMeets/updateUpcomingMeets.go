package main

import (
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"swimresults-backend/internal/config"
	"swimresults-backend/internal/database"
	"swimresults-backend/internal/database/repositories"
	updatemeet "swimresults-backend/updateMeet"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

// @TODO only for debugging purposes
const ONLY_FIRST_EVENT = false

var collyMyResults *colly.Collector
var repos *repositories.Repositories

const upcomingMeetsPageSelector = "div.col-xs-12.col-md-12.myresults_content_divtable"
const overviewPageSelector = "div.col-xs-12.col-md-10.msecm-no-padding.msecm-no-margin"
const msecmDetailsSelector = "div#custom-content"

func onUpcomingMeetsPage(e *colly.HTMLElement) {
	collyMyResults.OnHTMLDetach(upcomingMeetsPageSelector)
	wg := sync.WaitGroup{}

	e.ForEach(".myresults_content_divtablerow", func(i int, row *colly.HTMLElement) {
		country := row.ChildAttr("div.col-xs-1.text-right.myresults_content_divtable_right.myresults_padding_top_5 > img", "src")
		// Insert only meets which are in austria
		if country != "/images/flags/at.png" || ONLY_FIRST_EVENT && i > 0 {
			return
		}
		r := regexp.MustCompile("\\d+")
		meetId, err := strconv.Atoi(r.FindString(row.ChildAttr("a", "href")))
		if err != nil {
			panic(err)
		}
		wg.Add(1)
		go updatemeet.UpdateMeet(meetId, repos, &wg)
	})
	wg.Wait()
}

func updateUpcomingMeets() {
	collyMyResults = colly.NewCollector(colly.Async(true))
	collyMyResults.Limit(&colly.LimitRule{
		Delay:       5 * time.Second,
		Parallelism: 20,
	})

	collyMyResults.OnHTML(upcomingMeetsPageSelector, onUpcomingMeetsPage)

	collyMyResults.Visit("https://myresults.eu/de-DE/Meets/Today-Upcoming")
	collyMyResults.Wait()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig()

	err = cfg.ParseFlags()
	if err != nil {
		panic(err)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	repos = cfg.InitializeRepositories(db)

	log.Println("Updating Meets")
	updateUpcomingMeets()
}
