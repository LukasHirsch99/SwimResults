package updatemeet

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/database/models"
	"swimresults-backend/internal/database/repositories"
	"swimresults-backend/regex"
	updateschedule "swimresults-backend/updateSchedule"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

var collyMsecm *colly.Collector
var collyMyResults *colly.Collector
var repos *repositories.Repositories

const overviewPageSelector = "div.col-xs-12.col-md-10.msecm-no-padding.msecm-no-margin"
const msecmDetailsSelector = "div#custom-content"
const statisticsPageSelector = "div.col-xs-12.col-md-12.myresults_content_divtable"
const clubDetailsPageSelector = "div.col-xs-12.myresults_content_divtable"

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func extractMeetDate(s string, meet *models.Meet) {
	// 01.-05.08.2020
	r := regexp.MustCompile("((?<firstDay>^\\d{2})\\.)((?<firstMonth>\\d{2})\\.)?(-(?<lastDay>\\d{2})\\.)?((?<lastMonth>\\d{2})\\.)?(?<year>\\d{4}$)")

	m := regex.EvalRegex(r, s)
	l := len(m)

	if l == 4 {
		// 01.-05.08.2020
		// "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
		meet.Startdate, _ = time.Parse("02.01.2006", m["firstDay"]+"."+m["lastMonth"]+"."+m["year"])
		meet.Enddate, _ = time.Parse("02.01.2006", m["lastDay"]+"."+m["lastMonth"]+"."+m["year"])
	} else if l == 3 {
		// 03.10.2020
		meet.Startdate, _ = time.Parse("02.01.2006", m["firstDay"]+"-"+m["firstMonth"]+"-"+m["year"])
		meet.Enddate, _ = time.Parse("02.01.2006", m["firstDay"]+"-"+m["firstMonth"]+"-"+m["year"])
	} else if l == 5 {
		// 29.02.-01.03.2020
		meet.Startdate, _ = time.Parse("02.01.2006", m["firstDay"]+"-"+m["firstMonth"]+"-"+m["year"])
		meet.Enddate, _ = time.Parse("02.01.2006", m["lastDay"]+"-"+m["lastMonth"]+"-"+m["year"])
	}
}

func parseDeadline(s string) sql.NullTime {
	t, err := time.Parse("02.01.2006 15:04", s)
	return sql.NullTime{Time: t, Valid: err == nil}
}

func onMsecmDetails(e *colly.HTMLElement) {
	r := regexp.MustCompile("\\d+$")
	msecmId, err := strconv.Atoi(r.FindString(e.Request.URL.String()))
	if err != nil {
		panic(err)
	}

	meet, err := repos.MeetRepository.GetByMsecmId(msecmId)

	if err != nil {
		panic(err)
	}

	googleMapsLink := e.ChildAttr("p.text-right:nth-child(1) > a", "href")
	meet.Googlemapslink = sql.NullString{String: googleMapsLink, Valid: len(googleMapsLink) > 0}
	e.ForEach("a.hover-effect", func(i int, link *colly.HTMLElement) {
		href := link.Attr("href")
		if strings.Contains(href, ".pdf") {
			meet.Invitations = append(meet.Invitations, e.Request.URL.Hostname()+href)
		}
	})
	err = repos.MeetRepository.Upsert(meet)
	if err != nil {
		panic(err)
	}
}

func onOverview(e *colly.HTMLElement) {
	if !strings.HasSuffix(e.Request.URL.String(), "/Overview") {
		return
	}
	r := regexp.MustCompile("\\d+")
	meetId, _ := strconv.Atoi(r.FindString(e.Request.URL.String()))

	meet, err := repos.MeetRepository.GetById(meetId)

	if err != nil {
		meet = &models.Meet{}
	}

	imageLink := e.ChildAttr("img.img-responsive.center-block", "src")
	dateString := getOnlyChildText(e, "div:nth-child(4) > div")
	extractMeetDate(dateString, meet)

	meet.Id = meetId
	meet.Name = e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header:nth-child(1)")
	meet.Deadline = parseDeadline(getOnlyChildText(e, "div:nth-child(5) > div"))
	meet.Address = strings.Split(getOnlyChildText(e, "div:nth-child(7) > div"), "\t")[0]
	meet.Image = sql.NullString{String: "https://myresults.eu" + imageLink, Valid: len(imageLink) > 0}

	msecmLink := e.ChildAttr("div:nth-child(14) > div > a", "href")

	// Maybe defer repos.MeetRepository.Create(&meet) ?
	containsMsecmLink, msecmId := containsMsecmLink(msecmLink)
	meet.Msecmid = sql.NullInt16{Int16: int16(msecmId), Valid: containsMsecmLink}

	err = repos.MeetRepository.Upsert(meet)
	if err != nil {
		panic(err)
	}
	if containsMsecmLink {
		collyMsecm.Visit(msecmLink)
	}
}

func containsMsecmLink(msecmLink string) (bool, int) {
	if !strings.Contains(msecmLink, "msecm.at") {
		return false, -1
	}
	r := regexp.MustCompile("\\d+$")
	match := r.FindString(msecmLink)
	msecmId, err := strconv.Atoi(match)
	if err != nil {
		return false, -1
	}
	return true, msecmId
}

func onClubDetails(e *colly.HTMLElement) {
	if !strings.Contains(e.Request.URL.String(), "/Club/") {
		return
	}

	clubName := e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header td.myresults_personendetails_header")
	clubImage := e.ChildAttr("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header > div > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(2) > td.myresults_personendetails_text2 > img", "src")

	e.ForEach("div.row.tablecard.myresults_content_divtablerow", func(i int, swimmerEl *colly.HTMLElement) {
		name := swimmerEl.ChildText("div:nth-child(1) > div > a")
		swimmer := models.Swimmer{}

		r := regexp.MustCompile("\\d+$")
		swimmer.Clubid, _ = strconv.Atoi(r.FindString(e.Request.URL.String()))
		swimmer.Id, _ = strconv.Atoi(r.FindString(swimmerEl.ChildAttr("div:nth-child(1) > div > a", "href")))
		swimmer.Lastname, swimmer.Firstname = updateschedule.ParseName(name)

		details := swimmerEl.ChildText("div:nth-child(1) > div > span")
		swimmer.Gender = regexp.MustCompile("[A-Z]").FindString(details)
		birthyear, err := strconv.Atoi(regexp.MustCompile("\\d+").FindString(details))
		swimmer.Birthyear = sql.NullInt16{Int16: (int16(birthyear)), Valid: err == nil}
		swimmer.Isrelay = !swimmer.Birthyear.Valid

		err = repos.SwimmerRepository.Create(&swimmer)
		if err != nil {
			if strings.Contains(err.Error(), "insert or update on table \"swimmer\" violates foreign key constraint \"swimmer_clubid_fkey\"") {
				club := models.Club{
					Id:          swimmer.Clubid,
					Name:        clubName,
					Nationality: sql.NullString{String: e.Request.URL.String() + clubImage, Valid: len(clubImage) > 0},
				}
				err = repos.ClubRepository.Create(&club)
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	})
}

func onStatistics(e *colly.HTMLElement) {
	if !strings.HasSuffix(e.Request.URL.String(), "/Statistics") {
		return
	}

	e.ForEach("div.row.myresults_content_divtablerow", func(i int, row *colly.HTMLElement) {
		clubLink := strings.TrimSpace(row.ChildAttr("div:nth-child(2) > a", "href"))
		if clubLink == "" {
			return
		}
		collyMyResults.Visit("https://" + e.Request.URL.Host + clubLink)
	})
}

func UpdateMeet(meetId int, r *repositories.Repositories, wg *sync.WaitGroup) {
  if wg != nil {
    defer wg.Done()
  }
	log.Printf("Updating Meet: %d\n", meetId)
	repos = r
	collyMyResults = colly.NewCollector()
	collyMsecm = colly.NewCollector()

	collyMyResults.OnHTML(overviewPageSelector, onOverview)
	collyMyResults.OnHTML(statisticsPageSelector, onStatistics)
	// collyMyResults.OnHTML(clubDetailsPageSelector, onClubDetails)

	collyMsecm.OnHTML(msecmDetailsSelector, onMsecmDetails)

	collyMyResults.Visit(fmt.Sprintf("https://myresults.eu/de-DE/Meets/Today-Upcoming/%d/Overview", meetId))
	// collyMyResults.Visit(fmt.Sprintf("https://myresults.eu/de-DE/Meets/Today-Upcoming/%d/Overview/Statistics", meetId))

	collyMyResults.Wait()
	collyMyResults.Wait()

	updateschedule.UpdateSchedule(meetId, repos)
}
