package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
  "swimresults-backend/regex"

	"github.com/gocolly/colly"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase/postgrest-go"
)

var maxHeatId uint = 0
var maxResultId uint = 0
var clubIdSet []uint
var swimmerIdSet []uint

var heats []Heat
var starts []Start
var ageclasses []AgeClass
var results []Result
var swimmers []Swimmer
var clubs []Club

var m sync.Mutex
var wg sync.WaitGroup

type status string
const (
	EventStatusFinished status = "/images/status_green.png"
	EventStatusNext     status = "/images/status_yellow.png"
	EventStatusReady    status = "/images/status_grey.png"
	EventStatusNoInfo   status = "/images/status_none.png"
)

var client *supabase.Client

func getMaxIds() (uint, uint, uint, uint) {
  var maxIds map[string]uint
  err := json.Unmarshal([]byte(client.Rpc("maxids", "exact", "")), &maxIds)
  if err != nil {
    panic(err)
  }
  return maxIds["maxsessionid"], maxIds["maxeventid"], maxIds["maxheatid"], maxIds["maxresultid"]
}

func parseSessionInfo(row string) (SessionInfo, error) {
	str := strings.TrimSpace(row)
	r, err := regexp.Compile("\\d{2}:\\d{2}|\\d{4}|\\d{2}|\\d+")
	if err != nil {
		return SessionInfo{}, err
	}
	matches := r.FindAllString(str, -1)

	sessionInfo := SessionInfo{}

	displaynr, err := strconv.Atoi(matches[3])
	if err != nil {
		return SessionInfo{}, err
	}

	sessionInfo.day = fmt.Sprintf("%s-%s-%s", matches[2], matches[1], matches[0])
	sessionInfo.displaynr = uint(displaynr)

	if len(matches) == 6 {
		sessionInfo.warmupstart = matches[4] + ":00"
		sessionInfo.sessionstart = matches[5] + ":00"
	}

	return sessionInfo, nil
}

func parseEventInfo(row string) (EventInfo, error) {
	l := strings.Split(row, " - ")
	displaynr, err := strconv.Atoi(l[0])
	if err != nil {
		return EventInfo{}, errors.New("Couldn't convert displaynr to int")
	}
	return EventInfo{
		displaynr: uint(displaynr),
		name:      l[1],
	}, nil
}

func parseTime(tStr string) (string, error) {
	if len(tStr) == 0 {
		return "", errors.New("Time string empty")
	}
	var t time.Time

	if !strings.Contains(tStr, ":") {
		t, _ = time.Parse("05.00", tStr)
	} else if strings.Contains(tStr, "h") {
		t, _ = time.Parse("3h04:05.00", tStr)
	} else {
		t, _ = time.Parse("4:05.00", tStr)
	}
	return t.Format("15:04:05.00"), nil
}

func executeAndParse[T any](f *postgrest.FilterBuilder) (T, int64, error) {
	var err error
	var r T
	data, cnt, err := f.Execute()
	if err != nil {
		return r, 0, err
	}
	err = json.Unmarshal(data, &r)
	return r, cnt, err
}

func UintToString(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

func StringToUint(s string) uint {
	u, _ := strconv.Atoi(s)
	return uint(u)
}

func getSwimmerFromStartOrResult(swimmerId uint, row *colly.HTMLElement) {
	var swimmer Swimmer
	var club Club
	clubLink := row.ChildAttrs("a", "href")[1]
	r, _ := regexp.Compile("\\d+$")
	c, _ := strconv.Atoi(r.FindString(clubLink))
	clubId := uint(c)
	swimmer.Id = swimmerId
	swimmer.ClubId = clubId
	nameArray := strings.SplitN(row.ChildText("a"), " ", 1)
	swimmer.Lastname = nameArray[0]
	swimmer.Firstname = nameArray[1]

	details := row.ChildText(".hidden-xs.myresults_content_divtable_details")
	r, _ = regexp.Compile("\\d+|[A-Z]")
	birthAndGender := r.FindAllString(details, -1)
	if len(birthAndGender) == 1 {
		swimmer.Gender = birthAndGender[0]
		swimmer.IsRelay = true
	} else if len(birthAndGender) == 2 {
		swimmer.BirthYear = StringToUint(birthAndGender[0])
		swimmer.Gender = birthAndGender[1]
		swimmer.IsRelay = false
	}

	if !slices.Contains(clubIdSet, clubId) {
		club.Id = clubId
		club.Name = row.ChildText("div.hidden-xs.col-sm-4 > a")
		flagLink := row.ChildAttr("img", "src")
		if flagLink != "" {
			club.Nationality = "https://myresults.eu/" + flagLink
		}
		m.Lock()
		clubIdSet = append(clubIdSet, clubId)
		clubs = append(clubs, club)
		m.Unlock()
	}
	m.Lock()
	swimmerIdSet = append(swimmerIdSet, swimmerId)
	swimmers = append(swimmers, swimmer)
	m.Unlock()
	return
}

func populateStarts(meetId uint, startId uint, eventId uint) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#starts_content", func(e *colly.HTMLElement) {
		heatCnt := 0
		startCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				heatCnt++
			} else {
				startCnt++
			}
		})

		heatsWithStarts, dbHeatCnt, err := executeAndParse[[]HeatWithStarts](client.
			From("heat").
			Select("*, start!inner(*)", "exact", false).
			Eq("eventid", UintToString(eventId)))
		if err != nil {
			panic(err)
		}

		dbStartCnt := 0
		for _, h := range heatsWithStarts {
			dbStartCnt += len(h.Starts)
		}

		if startCnt == dbStartCnt && heatCnt == int(dbHeatCnt) {
				wg.Done()
				return
		}

		var heatNr uint = 0
		var heatId uint
		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				// Heat-Element
				heatNr++
				m.Lock()
				maxHeatId++
				heatId = maxHeatId
				heats = append(heats, Heat{
					Id:      heatId,
					EventId: eventId,
					HeatNr:  heatNr,
				})
				m.Unlock()
			} else {
				// Start-Element
				swimmerLink := row.ChildAttr("a", "href")
				r, _ := regexp.Compile("\\d+$")
				swimmerId := StringToUint(r.FindString(swimmerLink))
				lane := StringToUint(row.ChildText("div.col-xs-1"))
				startTime, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
				s := Start{
					HeatId:    heatId,
					SwimmerId: swimmerId,
					Lane:      lane,
				}
				if err == nil {
					s.Time.Set(startTime)
				} else {
					s.Time.SetNull()
				}

				m.Lock()
				starts = append(starts, s)
				m.Unlock()

				if !slices.Contains(swimmerIdSet, swimmerId) {
					getSwimmerFromStartOrResult(swimmerId, row)
				}
			}
		})
		wg.Done()
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Starts/" + UintToString(startId))
	return
}

func populateNewResults(meetId uint, resultId uint, eventId uint) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#starts_content", func(e *colly.HTMLElement) {
		resultCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if !strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				resultCnt++
			}
		})

		_, dbResultCnt, err := client.From("ageclass").Select("*, result!inner(*)", "exact", false).Eq("result.eventid", UintToString(eventId)).Execute()
		if err != nil {
			panic(err)
		}

		if resultCnt == int(dbResultCnt) {
			wg.Done()
			return
		}

		var ageClassName string
		swimmerIdToDBResultId := make(map[uint]uint)
		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				// Ageclass-Element
				ageClassName = strings.TrimSpace(row.Text)
			} else {
				// Result-Element
				swimmerLink := row.ChildAttr("a", "href")
				r, _ := regexp.Compile("\\d+$")
				swimmerId := StringToUint(r.FindString(swimmerLink))

				resultInfoString := row.ChildText("div.myresults_content_divtable_points")
				r = regexp.MustCompile("(?<timeToFirst>\\+\\d+\\.\\d+)|(?<reugeld>RG)|(?<finaPoints>\\d+)|(?<additionalInfo>[\\S]+$)")
				resultInfoMap := regex.EvalRegex(r, resultInfoString)

				dbResultId, ok := swimmerIdToDBResultId[swimmerId]
				if !ok {
					m.Lock()
					maxResultId++
					m.Unlock()
					dbResultId = maxResultId
					swimmerIdToDBResultId[swimmerId] = dbResultId
					result := Result{Id: dbResultId, EventId: eventId, SwimmerId: swimmerId}

					finaPointsString, ok := resultInfoMap["finaPoints"]
					if ok {
						result.FinaPoints.Set(StringToUint(finaPointsString))
					} else {
						result.FinaPoints.SetNull()
					}

					additionalInfo, ok := resultInfoMap["additionalInfo"]
					if ok {
						result.AdditionalInfo.Set(additionalInfo)
					} else {
						result.AdditionalInfo.SetNull()
					}

					_, ok = resultInfoMap["reugeld"]
					if ok && row.ChildText("div.myresults_content_divtable_points > span") == "" {
						result.Penalty = true
					} else {
						result.Penalty = false
					}

					splits := row.ChildText("span.myresults_content_divtable_details:nth-child(1)")
					if splits != "" {
						result.Splits.Set(splits)
					} else {
						result.Splits.SetNull()
					}
					time, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
					if err == nil {
						result.Time.Set(time)
					} else {
						result.Time.SetNull()
					}
					m.Lock()
					results = append(results, result)
					m.Unlock()
				}

				var ageclass AgeClass
				timeToFirst, ok := resultInfoMap["timeToFirst"]
				if ok {
					ageclass.TimeToFirst.Set(timeToFirst)
				} else {
					ageclass.TimeToFirst.SetNull()
				}

				position := row.ChildText("span.msecm-place")
				if position != "" {
					ageclass.Position.Set(StringToUint(strings.Replace(position, ".", "", 1)))
				} else {
					ageclass.Position.SetNull()
				}
				ageclass.Name = ageClassName
				ageclass.ResultId = dbResultId

				m.Lock()
				ageclasses = append(ageclasses, ageclass)
				m.Unlock()

				if !slices.Contains(swimmerIdSet, swimmerId) {
					getSwimmerFromStartOrResult(swimmerId, row)
				}
			}
		})
		wg.Done()
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Results/" + UintToString(resultId))
}

func updateSchedule(meetId uint, waitGroup *sync.WaitGroup) error {
	c := colly.NewCollector()

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
				eventCnt++
			}
		})

		sessionsWithEvents, dbSessionCnt, err := executeAndParse[[]SessionWithEvents](client.From("session").Select("*, event!inner(*)", "exact", false).Eq("meetid", UintToString(meetId)))
		if err != nil {
			panic(err)
		}

		dbEventCnt := 0
		for _, s := range sessionsWithEvents {
			dbEventCnt += len(s.Events)
		}
		scheduleUpToDate := (eventCnt == int(dbEventCnt) && sessionCnt == int(dbSessionCnt))

		if !scheduleUpToDate {
			client.From("session").Delete("*", "exact").Eq("meetid", UintToString(meetId)).Execute()
		}

		var newSessions []Session
		var newEvents []Event
		var sessionId uint
		var sessionIdx int

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				// Session-Item
				sessionInfo, _ := parseSessionInfo(row.Text)
				if scheduleUpToDate {
					sessionIdx = slices.IndexFunc(sessionsWithEvents, func(s SessionWithEvents) bool {
						return s.Displaynr == sessionInfo.displaynr && s.Day == sessionInfo.day
					})
					sessionId = sessionsWithEvents[sessionIdx].Id
				} else {
					session := sessionInfo.toSessionIncMaxId(meetId)
					sessionId = session.Id
					newSessions = append(newSessions, session)
				}
			} else {
				// Event-Item
				eventInfo, _ := parseEventInfo(strings.Split(row.ChildText(".col-xs-6"), "\t")[0])
				if eventInfo.name == "" || eventInfo.displaynr == 0 {
					return
				}
				var eventId uint
				if scheduleUpToDate {
					eventId = sessionsWithEvents[sessionIdx].Events[slices.IndexFunc(sessionsWithEvents[sessionIdx].Events, func(e Event) bool {
						return e.SessionId == sessionId && e.Name == eventInfo.name && e.DisplayNr == eventInfo.displaynr
					})].Id
				} else {
					event := eventInfo.toEventIncMaxId(sessionId)
					eventId = event.Id
					newEvents = append(newEvents, event)
				}

				href := row.ChildAttr(".myresults_content_link.myresults_content_divtablecol", "href")
				r, _ := regexp.Compile("\\d+$")
				startResultId := StringToUint(r.FindString(href))
				status := status(row.ChildAttr("div.col-xs-1.text-center.myresults_content_divtable_left.hidden-xs > img", "src"))

				if status == EventStatusNoInfo {
					return
				}
				wg.Add(1)
				go populateStarts(meetId, startResultId, eventId)
				if status == EventStatusFinished {
					// @TODO insert Results
					wg.Add(1)
					go populateNewResults(meetId, startResultId, eventId)
				}
			}
		})

		wg.Wait()

		if len(clubs) != 0 {
			_, _, err = client.From("club").Insert(clubs, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(swimmers) != 0 {
			_, _, err = client.From("swimmer").Insert(swimmers, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(newSessions) != 0 {
			_, _, err = client.From("session").Insert(newSessions, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(newEvents) != 0 {
			_, _, err = client.From("event").Insert(newEvents, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(heats) != 0 {
			_, _, err = client.From("heat").Insert(heats, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(starts) != 0 {
			_, _, err = client.From("start").Insert(starts, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(results) != 0 {
			_, _, err = client.From("result").Insert(results, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if len(ageclasses) != 0 {
			_, _, err = client.From("ageclass").Insert(ageclasses, false, "", "count", "exact").Execute()
			if err != nil {
				panic(err)
			}
		}
		if waitGroup != nil {
			waitGroup.Done()
		}
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Schedule")
	return nil
}

func main() {
	API_URL := "https://qeudknoyuvjztxvgbmou.supabase.co"
	API_KEY := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InFldWRrbm95dXZqenR4dmdibW91Iiwicm9sZSI6ImFub24iLCJpYXQiOjE2Njk0NzU0MjAsImV4cCI6MTk4NTA1MTQyMH0.xa0KNR2EEyJHyfEOJtuNFgbUa4H0e4rBWJ2w4dn49uU"
	var err error
	client, err = supabase.NewClient(API_URL, API_KEY, nil)
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	cSet, _, err := executeAndParse[[]map[string]uint](client.From("club").Select("id", "exact", false))
	if err != nil {
		fmt.Println("error getting clubids", err)
		return
	}
	sSet, _, err := executeAndParse[[]map[string]uint](client.From("swimmer").Select("id", "exact", false))
	if err != nil {
		fmt.Println("error getting swimmerids", err)
		return
	}
	for _, v := range cSet {
		clubIdSet = append(clubIdSet, v["id"])
	}
	for _, v := range sSet {
		swimmerIdSet = append(swimmerIdSet, v["id"])
	}

  maxSessionId, maxEventId, maxHeatId, maxResultId = getMaxIds()

	// wg.Add(1)
	// populateStarts(2088, 74127, 148174)
  startTime := time.Now()
	err = updateSchedule(2088, nil)
  fmt.Println("Took: ", time.Now().Sub(startTime))
	if err != nil {
		fmt.Println("error during update", err)
	}
}
