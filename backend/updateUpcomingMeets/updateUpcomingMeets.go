package updateupcomingmeets

import (
	"log/slog"
	"regexp"
	"strconv"
	"sync"
	"time"

	"swimresults-backend/internal/repository"
	updatemeet "swimresults-backend/updateMeet"

	"github.com/gocolly/colly"
)

// @TODO only for debugging purposes
const ONLY_FIRST_EVENT = false

var collyMyResults *colly.Collector
var repo *repository.Queries
var logger *slog.Logger

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
		go updatemeet.UpdateMeet(int32(meetId), repo, logger, &wg)
	})
	wg.Wait()
}

func UpdateUpcomingMeets(r *repository.Queries, l *slog.Logger) {
	repo = r
	logger = l
	logger.Info("Updating Upcoming Meets")
	collyMyResults = colly.NewCollector(colly.Async(true))
	collyMyResults.Limit(&colly.LimitRule{
		Delay:       5 * time.Second,
		Parallelism: 20,
	})

	collyMyResults.OnHTML(upcomingMeetsPageSelector, onUpcomingMeetsPage)

	collyMyResults.Visit("https://myresults.eu/de-DE/Meets/Today-Upcoming")
	collyMyResults.Wait()
}
