package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/database"
	"swimresults-backend/entities"
	"swimresults-backend/regex"
	"swimresults-backend/store"
	updateschedule "swimresults-backend/updateSchedule"
	"time"

	"github.com/gocolly/colly"
)

var collyMsecm *colly.Collector
var collyMyResults *colly.Collector

var meets = store.Meets

const overviewPageSelector = "div.col-xs-12.col-md-10.msecm-no-padding.msecm-no-margin"
const msecmDetailsSelector = "div#custom-content"

func getMeetByMsecmId(msecmId int) entities.Meet {
	m := meets.GetItemList()
	for _, meet := range m {
		if int(meet.MsecmId.Int64) == msecmId {
			return *meet
		}
	}
	panic(fmt.Sprintf("Meet not found with msecmId: %d", msecmId))
}

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func parseDate(s string) entities.MeetDate {
	var meetDate entities.MeetDate
	// 01.-05.08.2020
	r := regexp.MustCompile("((?<firstDay>^\\d{2})\\.)((?<firstMonth>\\d{2})\\.)?(-(?<lastDay>\\d{2})\\.)?((?<lastMonth>\\d{2})\\.)?(?<year>\\d{4}$)")

	m := regex.EvalRegex(r, s)
	l := len(m)

	if l == 4 {
		// 01.-05.08.2020
		meetDate.StartDate = m["year"] + "-" + m["lastMonth"] + "-" + m["firstDay"]
		meetDate.EndDate = m["year"] + "-" + m["lastMonth"] + "-" + m["lastDay"]
	} else if l == 3 {
		// 03.10.2020
		meetDate.StartDate = m["year"] + "-" + m["firstMonth"] + "-" + m["firstDay"]
		meetDate.EndDate = m["year"] + "-" + m["firstMonth"] + "-" + m["firstDay"]
	} else if l == 5 {
		// 29.02.-01.03.2020
		meetDate.StartDate = m["year"] + "-" + m["firstMonth"] + "-" + m["firstDay"]
		meetDate.EndDate = m["year"] + "-" + m["lastMonth"] + "-" + m["lastDay"]
	}
	return meetDate
}

func parseDeadline(s string) string {
	t, _ := time.Parse("02.01.2006 15:04", s)
	return t.Format("2006-01-02 15:04:05")
}

func onMsecmDetails(e *colly.HTMLElement) {
	r := regexp.MustCompile("\\d+$")
	msecmId, err := strconv.Atoi(r.FindString(e.Request.URL.String()))
	if err != nil {
		panic(err)
	}

	meet := getMeetByMsecmId(msecmId)

	googleMapsLink := e.ChildAttr("p.text-right:nth-child(1) > a", "href")
	if googleMapsLink != "" {
		meet.GoogleMapsLink.SetValid(googleMapsLink)
	}

	e.ForEach("a.hover-effect", func(i int, link *colly.HTMLElement) {
		href := link.Attr("href")
		if strings.Contains(href, ".pdf") {
			meet.Invitations = append(meet.Invitations, e.Request.URL.Hostname()+href)
		}
	})
	meets.SetItem(&meet)
}

func onOverview(e *colly.HTMLElement) {
	if !strings.Contains(e.Request.URL.String(), "/Overview") {
		return
	}
	r := regexp.MustCompile("\\d+")
	meetId := entities.StringToUint(r.FindString(e.Request.URL.String()))
	meet := meets.GetItemById(meetId)
	if meet == nil {
		meet = &entities.Meet{}
	}

	image := "https://myresults.eu" + e.ChildAttr("img.img-responsive.center-block", "src")
	dateString := getOnlyChildText(e, "div:nth-child(4) > div")

	meet.Id = meetId
	meet.Name = e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header:nth-child(1)")
	meet.MeetDate = parseDate(dateString)
	meet.Deadline = parseDeadline(getOnlyChildText(e, "div:nth-child(5) > div"))
	meet.Address = strings.Split(getOnlyChildText(e, "div:nth-child(7) > div"), "\t")[0]
	if image != "" {
		meet.Image.SetValid(image)
	}

	msecmLink := e.ChildAttr("div:nth-child(14) > div > a", "href")

	if strings.Contains(msecmLink, "msecm.at") {
		// Overview on MSECM-Website
		r := regexp.MustCompile("\\d+$")
		match := r.FindString(msecmLink)
		if match == "" {
			meets.SetItem(meet)
			return
		}
		msecmId, err := strconv.Atoi(match)
		if err != nil {
			fmt.Println(meet)
			panic(err)
		}
		meet.MsecmId.SetValid(int64(msecmId))
		meets.SetItem(meet)
		collyMsecm.Visit(msecmLink)
	} else {
		meets.SetItem(meet)
	}
}

func UpdateMeet(meetId uint) {
	collyMyResults = colly.NewCollector()
	collyMsecm = colly.NewCollector()

	collyMyResults.OnHTML(overviewPageSelector, onOverview)
	collyMsecm.OnHTML(msecmDetailsSelector, onMsecmDetails)

	collyMyResults.Visit("https://myresults.eu/de-DE/Meets/Today-Upcoming/" + strconv.FormatUint(uint64(meetId), 10) + "/Overview")
}

func main() {
	supabase, err := database.GetClient()
	if err != nil {
		panic(err)
  }
	UpdateMeet(2032)
	updateschedule.UpdateSchedule(2032, nil)

	supabase.Upsert(store.Meets)
	supabase.Insert(store.Swimmers)
	supabase.Insert(store.Clubs)
	supabase.Insert(store.Sessions)
	supabase.Insert(store.Events)
	supabase.Insert(store.Heats)
	supabase.Insert(store.Results)
	supabase.Insert(store.Starts)
	supabase.Insert(store.Ageclasses)
}
