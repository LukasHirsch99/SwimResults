package updatemeet

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/repository"
	updateschedule "swimresults-backend/updateSchedule"
	"sync"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)

var collyMsecm *colly.Collector
var collyMyResults *colly.Collector
var repo *repository.Queries

const overviewPageSelector = "div.col-xs-12.col-md-10.msecm-no-padding.msecm-no-margin"
const msecmDetailsSelector = "div#custom-content"
const statisticsPageSelector = "div.col-xs-12.col-md-12.myresults_content_divtable"

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func onMsecmDetails(e *colly.HTMLElement) {
	r := regexp.MustCompile(`\d+$`)
	msecmId, err := strconv.Atoi(r.FindString(e.Request.URL.String()))
	if err != nil {
		panic(err)
	}

	meet, err := repo.GetMeetByMsecmId(context.Background(), pgtype.Int4{Valid: true, Int32: int32(msecmId)})

	if err != nil {
		panic(err)
	}

	googleMapsLink := e.ChildAttr("p.text-right:nth-child(1) > a", "href")
	meet.Googlemapslink = pgtype.Text{String: googleMapsLink, Valid: len(googleMapsLink) > 0}
	e.ForEach("a.hover-effect", func(i int, link *colly.HTMLElement) {
		href := link.Attr("href")
		if strings.Contains(href, ".pdf") {
			meet.Invitations = append(meet.Invitations, e.Request.URL.Hostname()+href)
		}
	})
	err = repo.UpsertMeet(context.Background(), repository.UpsertMeetParams{
		ID:             meet.ID,
		Name:           meet.Name,
		Image:          meet.Image,
		Invitations:    meet.Invitations,
		Deadline:       meet.Deadline,
		Address:        meet.Address,
		Startdate:      meet.Startdate,
		Enddate:        meet.Enddate,
		Googlemapslink: meet.Googlemapslink,
		Msecmid:        meet.Msecmid,
	})
	if err != nil {
		panic(err)
	}
}

func onOverview(e *colly.HTMLElement) {
	if !strings.HasSuffix(e.Request.URL.String(), "/Overview") {
		return
	}
	r := regexp.MustCompile(`\d+`)
	meetId, _ := strconv.Atoi(r.FindString(e.Request.URL.String()))

	meet, err := repo.GetMeetById(context.Background(), int32(meetId))

	if err != nil {
		meet = repository.Meet{}
	}

	imageLink := e.ChildAttr("img.img-responsive.center-block", "src")
	dateString := getOnlyChildText(e, "div:nth-child(4) > div")
	sD, eD, err := extractMeetDate(dateString)
	meet.Startdate = pgtype.Date{Time: sD, Valid: err == nil}
	meet.Enddate = pgtype.Date{Time: eD, Valid: err == nil}

	meet.ID = int32(meetId)
	meet.Name = e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header:nth-child(1)")
	dl, err := parseDeadline(getOnlyChildText(e, "div:nth-child(5) > div"))
	meet.Deadline = pgtype.Timestamp{Time: dl, Valid: err == nil}
	meet.Address = strings.Split(getOnlyChildText(e, "div:nth-child(7) > div"), "\t")[0]
	meet.Image = pgtype.Text{String: "https://myresults.eu" + imageLink, Valid: len(imageLink) > 0}

	msecmLink := e.ChildAttr("div:nth-child(14) > div > a", "href")

	// Maybe defer repos.MeetRepository.Create(&meet) ?
	containsMsecmLink, msecmId := containsMsecmLink(msecmLink)
	meet.Msecmid = pgtype.Int4{Int32: int32(msecmId), Valid: containsMsecmLink}

	err = repo.UpsertMeet(context.Background(), repository.UpsertMeetParams{
		ID:             meet.ID,
		Name:           meet.Name,
		Image:          meet.Image,
		Invitations:    meet.Invitations,
		Deadline:       meet.Deadline,
		Address:        meet.Address,
		Startdate:      meet.Startdate,
		Enddate:        meet.Enddate,
		Googlemapslink: meet.Googlemapslink,
		Msecmid:        meet.Msecmid,
	})
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
	r := regexp.MustCompile(`\d+$`)
	match := r.FindString(msecmLink)
	msecmId, err := strconv.Atoi(match)
	if err != nil {
		return false, -1
	}
	return true, msecmId
}

// func onClubDetails(e *colly.HTMLElement) {
// 	if !strings.Contains(e.Request.URL.String(), "/Club/") {
// 		return
// 	}
//
// 	clubName := e.ChildText("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header td.myresults_personendetails_header")
// 	clubImage := e.ChildAttr("div.row.myresults_content_divtablerow.myresults_content_divtablerow_header > div > table > tbody > tr > td:nth-child(2) > table > tbody > tr:nth-child(2) > td.myresults_personendetails_text2 > img", "src")
//
// 	e.ForEach("div.row.tablecard.myresults_content_divtablerow", func(i int, swimmerEl *colly.HTMLElement) {
// 		name := swimmerEl.ChildText("div:nth-child(1) > div > a")
// 		swimmer := repository.CreateSwimmerParams{}
//
// 		r := regexp.MustCompile(`\d+$`)
// 		clubId, _ := strconv.Atoi(r.FindString(e.Request.URL.String()))
// 		swimmer.Clubid = int32(clubId)
// 		swimmerId, _ := strconv.Atoi(r.FindString(swimmerEl.ChildAttr("div:nth-child(1) > div > a", "href")))
// 		swimmer.ID = int32(swimmerId)
// 		swimmer.Lastname, swimmer.Firstname = updateschedule.ParseName(name)
//
// 		details := swimmerEl.ChildText("div:nth-child(1) > div > span")
// 		swimmer.Gender = repository.Gender(regexp.MustCompile("[A-Z]").FindString(details))
// 		birthyear, err := strconv.Atoi(regexp.MustCompile(`\d+`).FindString(details))
// 		swimmer.Birthyear = pgtype.Int4{Int32: int32(birthyear), Valid: err == nil}
//
// 		err = repo.CreateSwimmer(context.Background(), swimmer)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "insert or update on table \"swimmer\" violates foreign key constraint \"swimmer_clubid_fkey\"") {
// 				club := repository.CreateClubParams{
// 					ID:          swimmer.Clubid,
// 					Name:        clubName,
// 					Nationality: pgtype.Text{String: e.Request.URL.String() + clubImage, Valid: len(clubImage) > 0},
// 				}
// 				err = repo.CreateClub(context.Background(), club)
// 				if err != nil {
// 					panic(err)
// 				}
// 			} else {
// 				panic(err)
// 			}
// 		}
// 	})
// }

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

func UpdateMeet(meetId int32, r *repository.Queries, l *slog.Logger, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	repo = r
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

	updateschedule.UpdateSchedule(meetId, repo, l)
}
