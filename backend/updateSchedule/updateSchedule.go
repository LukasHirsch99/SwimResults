package updateschedule

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/database/models"
	"swimresults-backend/internal/database/repositories"
	"sync"
	"time"
	"unicode"

	"github.com/gocolly/colly"
)

// @TODO remove after debugging
var DEBUG_MODE = true

type status string

const (
	EventStatusNoInfo   status = "/images/status_none.png"
	EventStatusReady    status = "/images/status_grey.png"
	EventStatusNext     status = "/images/status_yellow.png"
	EventStatusFinished status = "/images/status_green.png"
)

var repos *repositories.Repositories
var startResultWg sync.WaitGroup
var ensureSwimmerExistsLock sync.RWMutex
var swimmerIds []int
var clubIds []int

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func extractSessionInfo(row string, session *models.Session) error {
	str := strings.TrimSpace(row)
	r := regexp.MustCompile("\\d{2}\\.\\d{2}\\.\\d{4}|\\d{2}:\\d{2}")
	matches := r.FindAllString(str, -1)

	if len(matches) == 0 {
		return errors.New(fmt.Sprintf("Couldn't find date information: %s", str))
	}

	// "02 Jan 06 15:04 MST"
	// Samstag 21.09.2024 - 1. Abschnitt - Einschwimmen 10:00, Beginn 11:10
	day, err := time.Parse("02.01.2006", matches[0])
	if err != nil {
		session.Warmupstart = sql.NullTime{Valid: false}
		session.Sessionstart = sql.NullTime{Valid: false}
		return err
	}

	session.Warmupstart = sql.NullTime{Time: day, Valid: true}
	session.Sessionstart = sql.NullTime{Time: day, Valid: true}

	if len(matches) == 3 {
		warmupTime, err := time.Parse("15:04", matches[1])
		sessionTime, err2 := time.Parse("15:04", matches[2])

		if err != nil || err2 != nil {
			session.Warmupstart = sql.NullTime{Valid: false}
			session.Sessionstart = sql.NullTime{Valid: false}
			return err
		}

		session.Warmupstart.Time = day.Add(
			time.Hour*time.Duration(warmupTime.Hour()) +
				time.Minute*time.Duration(warmupTime.Minute()) +
				time.Second*time.Duration(warmupTime.Second()),
		)
		session.Sessionstart.Time = day.Add(
			time.Hour*time.Duration(sessionTime.Hour()) +
				time.Minute*time.Duration(sessionTime.Minute()) +
				time.Second*time.Duration(sessionTime.Second()),
		)
	}
	return nil
}

func extractEventInfo(row string, model *models.Event) error {
	l := strings.Split(row, " - ")
	displaynr, err := strconv.Atoi(l[0])
	if err != nil {
		return errors.New("Couldn't convert displaynr to int")
	}
	if l[1] == "" {
		return errors.New("Event Name empty")
	}
	model.Displaynr = displaynr
	model.Name = l[1]
	return nil
}

func parseTime(tStr string) (time.Time, error) {
	if len(tStr) == 0 {
		return time.Time{}, errors.New("Time string empty")
	}
	var t time.Time

	if !strings.Contains(tStr, ":") {
		t, _ = time.Parse("05.00", tStr)
	} else if strings.Contains(tStr, "h") {
		t, _ = time.Parse("3h04:05.00", tStr)
	} else {
		t, _ = time.Parse("4:05.00", tStr)
	}
	return t, nil
}

func ParseName(name string) (string, string) {
	nameArray := strings.SplitN(name, " ", 2)
	var nextCap bool
	lastnameRunes := []rune(nameArray[0])
	firstnameRunes := []rune(nameArray[1])
	for i, c := range lastnameRunes {
		if c == '-' {
			nextCap = true
		} else if nextCap || i == 0 {
			lastnameRunes[i] = unicode.ToUpper(c)
			nextCap = false
		} else {
			lastnameRunes[i] = unicode.ToLower(c)
		}
	}
	for i, c := range firstnameRunes {
		if c == '-' {
			nextCap = true
		} else if nextCap || i == 0 {
			firstnameRunes[i] = unicode.ToUpper(c)
			nextCap = false
		} else {
			firstnameRunes[i] = unicode.ToLower(c)
		}
	}
	return string(lastnameRunes), string(firstnameRunes)
}

func extractClubFromStartOrResult(clubId int, row *colly.HTMLElement) *models.Club {
	club := models.Club{
		Id:   clubId,
		Name: row.ChildText("div.hidden-xs.col-sm-4 > a"),
	}
	flagLink := row.ChildAttr("img", "src")
	if flagLink != "" {
		club.Nationality = sql.NullString{String: "https://myresults.eu/" + flagLink, Valid: true}
	}
	return &club
}

func extractSwimmerFromStartOrResult(swimmerId int, row *colly.HTMLElement) *models.Swimmer {
	var swimmer models.Swimmer
	clubLink := row.ChildAttrs("a", "href")[1]
	r := regexp.MustCompile("\\d+$")
	clubId, _ := strconv.Atoi(r.FindString(clubLink))
	swimmer.Id = int(swimmerId)
	swimmer.Clubid = clubId
	nameString := getOnlyChildText(row, "div.col-xs-11.col-sm-4 > a")
	swimmer.Lastname, swimmer.Firstname = ParseName(nameString)

	details := row.ChildText("div.col-xs-11.col-sm-4 > span")
	swimmer.Gender = regexp.MustCompile("[A-Z]").FindString(details)
	birthyear, err := strconv.Atoi(regexp.MustCompile("\\d+").FindString(details))
	swimmer.Birthyear = sql.NullInt16{Int16: (int16(birthyear)), Valid: err == nil}
	return &swimmer
}

func getClubAndSwimmerIdFromRow(row *colly.HTMLElement) (int, int) {
	r := regexp.MustCompile("\\d+$")
	clubLink := row.ChildAttr("div.hidden-xs.col-sm-4 > a", "href")
	clubid, err := strconv.Atoi(r.FindString(clubLink))
	if err != nil {
		panic(err)
	}
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerid, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		panic(err)
	}
	return clubid, swimmerid
}

func ensureSwimmerExists(row *colly.HTMLElement) {
	r := regexp.MustCompile("\\d+$")
	clubLink := row.ChildAttr("div.hidden-xs.col-sm-4 > a", "href")
	clubId, err := strconv.Atoi(r.FindString(clubLink))
	if err != nil {
		panic(err)
	}
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerId, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		panic(err)
	}

	ensureSwimmerExistsLock.Lock()
	defer ensureSwimmerExistsLock.Unlock()

	// if swimmer exists, return
	if repos.SwimmerRepository.CheckId(swimmerId) {
		return
	}

	// if club doesn't exist, first create club
	clubExists := repos.ClubRepository.CheckId(clubId)
	swimmer := extractSwimmerFromStartOrResult(swimmerId, row)
	club := extractClubFromStartOrResult(clubId, row)

	if !clubExists {
		err := repos.ClubRepository.Create(club)
		if err != nil {
			panic(err)
		}
	}

	err = repos.SwimmerRepository.Create(swimmer)
	if err != nil {
		fmt.Println(swimmer)
		panic(err)
	}
}

func UpdateSchedule(meetId int, r *repositories.Repositories) {
	repos = r
	log.Printf("Updating Schedule for: %d", meetId)
	c := colly.NewCollector()

	var err error
	swimmerIds, err = repos.SwimmerRepository.GetIds()
	if err != nil {
		panic(err)
	}
	clubIds, err = repos.ClubRepository.GetIds()
	if err != nil {
		panic(err)
	}

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div.col-xs-12.col-md-12.myresults_content_divtable", func(e *colly.HTMLElement) {
		sessionCnt := 0
		eventCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				sessionCnt++
			} else {
				if getOnlyChildText(row, "div.col-xs-6") != "-" {
					eventCnt++
				}
			}
		})

		dbSessionCnt, err := repos.SessionRepository.CountForMeet(meetId)
		if err != nil {
			panic(err)
		}

		dbEventCnt, err := repos.EventRepository.CountForMeet(meetId)
		if err != nil {
			panic(err)
		}

		scheduleUpToDate := (eventCnt == dbEventCnt && sessionCnt == dbSessionCnt) && !DEBUG_MODE

		if scheduleUpToDate {
			return
		}

		err = repos.SessionRepository.DeleteForMeet(meetId)
		if err != nil {
			panic(err)
		}

		displayNr := 0
		var session models.Session

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Session-Item
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				displayNr++
				session = models.Session{
					Meetid:    meetId,
					Displaynr: displayNr,
				}
				err := extractSessionInfo(row.Text, &session)
				if err != nil {
					panic(err)
				}

				if !scheduleUpToDate {
					err = repos.SessionRepository.Create(&session)
				} else {
					err = repos.SessionRepository.GetByPK(&session)
				}
				if err != nil {
					panic(err)
				}

				// Event-Item
			} else {
				event := models.Event{
					Sessionid: session.Id,
				}
				err := extractEventInfo(strings.Split(row.ChildText(".col-xs-6"), "\t")[0], &event)
				if err != nil {
					panic(err)
				}

				if !scheduleUpToDate {
					err = repos.EventRepository.Create(&event)
				} else {
					err = repos.EventRepository.GetByPK(&event)
				}
				if err != nil {
					panic(err)
				}

				status := status(row.ChildAttr("div.col-xs-1.text-center.myresults_content_divtable_left.hidden-xs > img", "src"))

				if status == EventStatusNoInfo {
					return
				}

				href := row.ChildAttr(".myresults_content_link.myresults_content_divtablecol", "href")
				r := regexp.MustCompile("\\d+$")
				startResultId, err := strconv.Atoi(r.FindString(href))

				startResultWg.Add(1)
				go populateStarts(meetId, startResultId, event.Id)
				if status == EventStatusFinished {
					startResultWg.Add(1)
					go populateNewResults(meetId, startResultId, event.Id)
				}
			}
		})
		startResultWg.Wait()
	})

	c.Visit(fmt.Sprint("https://myresults.eu/de-DE/Meets/Recent/", meetId, "/Schedule"))
}
