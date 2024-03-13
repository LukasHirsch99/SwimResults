package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	. "swimresults-backend/database"
	"swimresults-backend/regex"

	"github.com/gocolly/colly"
	"github.com/guregu/null/v5"
)

var maxHeatId uint = 0
var maxResultId uint = 0
var maxSessionId uint = 0
var maxEventId uint = 0

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

var supabase *Supabase

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

	sessionInfo.Day = fmt.Sprintf("%s-%s-%s", matches[2], matches[1], matches[0])
	sessionInfo.DisplayNr = uint(displaynr)

	if len(matches) == 6 {
		sessionInfo.WarmupStart = matches[4] + ":00"
		sessionInfo.SessionStart = matches[5] + ":00"
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
		DisplayNr: uint(displaynr),
		Name:      l[1],
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
		swimmer.BirthYear = null.IntFrom(StringToInt(birthAndGender[0]))
		swimmer.Gender = birthAndGender[1]
		swimmer.IsRelay = false
	}

	if !slices.Contains(clubIdSet, clubId) {
		club.Id = clubId
		club.Name = row.ChildText("div.hidden-xs.col-sm-4 > a")
		flagLink := row.ChildAttr("img", "src")
		if flagLink != "" {
			club.Nationality.SetValid("https://myresults.eu/" + flagLink)
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

		heatsWithStarts, dbHeatCnt, err := supabase.GetHeatsWithStartsByEventid(eventId)
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

    supabase.DeleteHeatsByEventId(eventId)

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
					s.Time.SetValid(startTime)
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
		ageclassCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if !strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				ageclassCnt++
			}
		})

    dbAgeclassCnt, err := supabase.GetAgeclassCntByEventId(eventId)
		if err != nil {
			panic(err)
		}

		if ageclassCnt == int(dbAgeclassCnt) {
			wg.Done()
			return
		}

    supabase.DeleteResultsByEventId(eventId)

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
						result.FinaPoints.SetValid(StringToInt(finaPointsString))
					}

					additionalInfo, ok := resultInfoMap["additionalInfo"]
					if ok {
						result.AdditionalInfo.SetValid(additionalInfo)
					}

					_, ok = resultInfoMap["reugeld"]
					if ok && row.ChildText("div.myresults_content_divtable_points > span") == "" {
						result.Penalty = true
					} else {
						result.Penalty = false
					}

					splits := row.ChildText("span.myresults_content_divtable_details:nth-child(1)")
					if splits != "" {
						result.Splits.SetValid(splits)
					}
					time, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
					if err == nil {
						result.Time.SetValid(time)
					}
					m.Lock()
					results = append(results, result)
					m.Unlock()
				}

				var ageclass AgeClass
				timeToFirst, ok := resultInfoMap["timeToFirst"]
				if ok {
					ageclass.TimeToFirst.SetValid(timeToFirst)
				}

				position := row.ChildText("span.msecm-place")
				if position != "" {
					ageclass.Position.SetValid(StringToInt(strings.Replace(position, ".", "", 1)))
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

		sessionsWithEvents, dbSessionCnt, err := supabase.GetSessionsWithEventsByMeetId(meetId)
		if err != nil {
			panic(err)
		}

		dbEventCnt := 0
		for _, s := range sessionsWithEvents {
			dbEventCnt += len(s.Events)
		}
		scheduleUpToDate := (eventCnt == int(dbEventCnt) && sessionCnt == int(dbSessionCnt))

		if !scheduleUpToDate {
      supabase.DeleteSessionsByMeetId(meetId)
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
						return s.Displaynr == sessionInfo.DisplayNr && s.Day == sessionInfo.Day
					})
					sessionId = sessionsWithEvents[sessionIdx].Id
				} else {
					maxSessionId++
					sessionId = maxSessionId
					session := sessionInfo.ToSession(meetId, sessionId)
					newSessions = append(newSessions, session)
				}
			} else {
				// Event-Item
				eventInfo, _ := parseEventInfo(strings.Split(row.ChildText(".col-xs-6"), "\t")[0])
				if eventInfo.Name == "" || eventInfo.DisplayNr == 0 {
					return
				}
				var eventId uint
				if scheduleUpToDate {
					eventId = sessionsWithEvents[sessionIdx].Events[slices.IndexFunc(sessionsWithEvents[sessionIdx].Events, func(e Event) bool {
						return e.SessionId == sessionId && e.Name == eventInfo.Name && e.DisplayNr == eventInfo.DisplayNr
					})].Id
				} else {
					maxEventId++
					eventId = maxEventId
					event := eventInfo.ToEvent(sessionId, eventId)
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

		err = supabase.InsertInto("club", clubs)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("swimmer", swimmers)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("session", newSessions)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("event", newEvents)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("heat", heats)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("start", starts)
		if err != nil {
			panic(err)
		}

		err = supabase.InsertInto("result", results)
		if err != nil {
			panic(err)
		}
		err = supabase.InsertInto("ageclass", ageclasses)
		if err != nil {
			panic(err)
		}

		if waitGroup != nil {
			waitGroup.Done()
		}
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Schedule")
	return nil
}

func main() {
  var err error
  supabase, err = InitClient()
	if err != nil {
		fmt.Println("Couldn't initalize client: ", err)
	}

	swimmerIdSet, err = supabase.GetSwimmerIdSet()
	if err != nil {
		fmt.Println("error getting swimmerids: ", err)
		return
	}
	clubIdSet, err = supabase.GetClubIdSet()
	if err != nil {
		fmt.Println("error getting clubids: ", err)
		return
	}
	maxIds, err := supabase.GetMaxIds()
	if err != nil {
		fmt.Println("error getting maxids: ", err)
		return
	}
	maxSessionId, maxEventId, maxHeatId, maxResultId = maxIds.MaxSessionId, maxIds.MaxEventId, maxIds.MaxHeatId, maxIds.MaxResultId

	startTime := time.Now()
	err = updateSchedule(2088, nil)
	fmt.Println("Update took: ", time.Now().Sub(startTime))

	if err != nil {
		fmt.Println("error during update: ", err)
	}
}
