package updateschedule

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/repository"
	"sync"
	"time"
	"unicode"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)

// @TODO remove after debugging
var DEBUG_MODE = false

type status string

const (
	EventStatusNoInfo   status = "/images/status_none.png"
	EventStatusReady    status = "/images/status_grey.png"
	EventStatusNext     status = "/images/status_yellow.png"
	EventStatusFinished status = "/images/status_green.png"
)

var repo *repository.Queries
var startResultWg sync.WaitGroup
var ensureSwimmerExistsLock sync.RWMutex
var swimmerIds []int32
var clubIds []int32

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func extractSessionInfo(row string, session *repository.Session) error {
	str := strings.TrimSpace(row)
	r := regexp.MustCompile("\\d{2}\\.\\d{2}\\.\\d{4}|\\d{2}:\\d{2}")
	matches := r.FindAllString(str, -1)

	if len(matches) == 0 {
		return errors.New(fmt.Sprintf("Couldn't find date information: %s", str))
	}

	session.Day = pgtype.Date{Valid: false}
	session.Warmupstart = pgtype.Time{Valid: false}
	session.Sessionstart = pgtype.Time{Valid: false}

	// "02 Jan 06 15:04 MST"
	// Samstag 21.09.2024 - 1. Abschnitt - Einschwimmen 10:00, Beginn 11:10
	day, err := time.Parse("02.01.2006", matches[0])
	if err != nil {
		return err
	}

	session.Day = pgtype.Date{Valid: true, Time: day}

	if len(matches) == 3 {
		warmupTime, err := time.Parse("15:04", matches[1])
		sessionTime, err2 := time.Parse("15:04", matches[2])

		if err != nil || err2 != nil {
			return err
		}

		session.Warmupstart = pgtype.Time{Valid: true, Microseconds: warmupTime.UnixMicro() - day.UnixMicro()}
		session.Sessionstart = pgtype.Time{Valid: true, Microseconds: sessionTime.UnixMicro() - day.UnixMicro()}
	}
	return nil
}

func extractEventInfo(row string, model *repository.Event) error {
	l := strings.Split(row, " - ")
	displaynr, err := strconv.Atoi(l[0])
	if err != nil {
		return errors.New("Couldn't convert displaynr to int")
	}
	if l[1] == "" {
		return errors.New("Event Name empty")
	}
	model.Displaynr = int32(displaynr)
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
	return time.Date(1970, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location()), nil
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

func CreateClubParamsFromStartOrResult(clubId int32, row *colly.HTMLElement) repository.CreateClubParams {
	club := repository.CreateClubParams{
		ID:   clubId,
		Name: row.ChildText("div.hidden-xs.col-sm-4 > a"),
	}
	flagLink := row.ChildAttr("img", "src")
	if flagLink != "" {
		club.Nationality = pgtype.Text{String: "https://myresults.eu/" + flagLink, Valid: true}
	}
	return club
}

func CreateSwimmerParamsFromStartOrResult(swimmerId int32, row *colly.HTMLElement) repository.CreateSwimmerParams {
	var swimmer repository.CreateSwimmerParams
	clubLink := row.ChildAttrs("a", "href")[1]
	r := regexp.MustCompile("\\d+$")
	clubId, _ := strconv.Atoi(r.FindString(clubLink))
	swimmer.ID = swimmerId
	swimmer.Clubid = int32(clubId)
	nameString := getOnlyChildText(row, "div.col-xs-11.col-sm-4 > a")
	swimmer.Lastname, swimmer.Firstname = ParseName(nameString)

	details := row.ChildText("div.col-xs-11.col-sm-4 > span")
	swimmer.Gender = repository.Gender(regexp.MustCompile("[A-Z]").FindString(details))
	birthyear, err := strconv.Atoi(regexp.MustCompile("\\d+").FindString(details))
	swimmer.Birthyear = pgtype.Int4{Int32: int32(birthyear), Valid: err == nil}
	return swimmer
}

func getClubAndSwimmerIdFromRow(row *colly.HTMLElement) (int32, int32) {
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
	return int32(clubid), int32(swimmerid)
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
	exists, _ := repo.CheckSwimmerId(context.Background(), int32(swimmerId))
	if exists {
		return
	}

	// if club doesn't exist, first create club
	clubExists, _ := repo.CheckClubId(context.Background(), int32(clubId))
	createSwimmer := CreateSwimmerParamsFromStartOrResult(int32(swimmerId), row)
	createClub := CreateClubParamsFromStartOrResult(int32(clubId), row)

	if !clubExists {
		err := repo.CreateClub(context.Background(), createClub)

		if err != nil {
			panic(err)
		}
	}

	err = repo.CreateSwimmer(context.Background(), createSwimmer)

	if err != nil {
		fmt.Println(createSwimmer)
		panic(err)
	}
}

func UpdateSchedule(meetId int32, r *repository.Queries) {
	repo = r
	log.Printf("Updating Schedule for: %d", meetId)
	c := colly.NewCollector()

	var err error
	swimmerIds, err = repo.GetSwimmerIds(context.Background())
	if err != nil {
		panic(err)
	}
	clubIds, err = repo.GetClubIds(context.Background())
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

		dbSessionCnt, err := repo.GetSessionCntForMeet(context.Background(), meetId)
		if err != nil {
			panic(err)
		}

		dbEventCnt, err := repo.GetEventCntForMeet(context.Background(), meetId)
		if err != nil {
			panic(err)
		}

		scheduleUpToDate := (eventCnt == int(dbEventCnt) && sessionCnt == int(dbSessionCnt)) && !DEBUG_MODE
		if scheduleUpToDate {
			return
		}

		err = repo.DeleteSessionsForMeet(context.Background(), meetId)
		if err != nil {
			panic(err)
		}

		displayNr := 0
		session := repository.Session{}

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Session-Item
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				displayNr++
				err := extractSessionInfo(row.Text, &session)
				if err != nil {
					panic(err)
				}

				if !scheduleUpToDate {
					session, err = repo.CreateSession(context.Background(), repository.CreateSessionParams{
						Meetid:    meetId,
						Displaynr: int32(displayNr),
						Day:       session.Day,
					})
				} else {
					session, err = repo.GetSessionByPk(context.Background(), repository.GetSessionByPkParams{
						Meetid:    meetId,
						Displaynr: int32(displayNr),
					})
				}
				if err != nil {
					panic(err)
				}
				// Event-Item
			} else {
				event := repository.Event{}
				err := extractEventInfo(strings.Split(row.ChildText(".col-xs-6"), "\t")[0], &event)
				if err != nil {
					fmt.Println("Meetid: ", session.Meetid)
					return
					// panic(err)
				}

				if !scheduleUpToDate {
					event, err = repo.CreateEvent(context.Background(), repository.CreateEventParams{
						Sessionid: session.ID,
						Displaynr: event.Displaynr,
						Name:      event.Name,
					})
				} else {
					event, err = repo.GetEventByPk(context.Background(), repository.GetEventByPkParams{
						Sessionid: event.Sessionid,
						Displaynr: event.Displaynr,
						Name:      event.Name,
					})
				}
				if err != nil {
					fmt.Println(event)
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
				go populateStarts(meetId, startResultId, event.ID)
				if status == EventStatusFinished {
					startResultWg.Add(1)
					go populateNewResults(meetId, startResultId, event.ID)
				}
			}
		})
		startResultWg.Wait()
	})

	c.Visit(fmt.Sprint("https://myresults.eu/de-DE/Meets/Recent/", meetId, "/Schedule"))
}
