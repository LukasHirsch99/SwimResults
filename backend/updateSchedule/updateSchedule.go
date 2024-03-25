package updateschedule

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	. "swimresults-backend/database"
	"swimresults-backend/entities"
	"swimresults-backend/regex"
	"swimresults-backend/store"
	"sync"
	"time"
	"unicode"

	"github.com/gocolly/colly"
	"github.com/guregu/null/v5"
)

var sessions = store.Sessions
var events = store.Events
var heats = store.Heats
var results = store.Results
var clubs = store.Clubs
var swimmers = store.Swimmers
var starts = store.Starts
var ageclasses = store.Ageclasses

var wg sync.WaitGroup
var supabase *Supabase

type status string

const (
	EventStatusFinished status = "/images/status_green.png"
	EventStatusNext     status = "/images/status_yellow.png"
	EventStatusReady    status = "/images/status_grey.png"
	EventStatusNoInfo   status = "/images/status_none.png"
)

func getOnlyChildText(e *colly.HTMLElement, selector string) string {
	return strings.TrimSpace(e.DOM.Find(selector).First().Clone().Children().Remove().End().Text())
}

func parseSessionInfo(row string) (entities.SessionInfo, error) {
	str := strings.TrimSpace(row)
	r := regexp.MustCompile("\\d{2}:\\d{2}|\\d{4}|\\d{2}|\\d+")
	matches := r.FindAllString(str, -1)

	sessionInfo := entities.SessionInfo{}

	displaynr, err := strconv.Atoi(matches[3])
	if err != nil {
		return entities.SessionInfo{}, err
	}

	sessionInfo.Day = fmt.Sprintf("%s-%s-%s", matches[2], matches[1], matches[0])
	sessionInfo.DisplayNr = uint(displaynr)

	if len(matches) == 6 {
		sessionInfo.WarmupStart.SetValid(matches[4] + ":00")
		sessionInfo.SessionStart.SetValid(matches[5] + ":00")
	}

	return sessionInfo, nil
}

func parseEventInfo(row string) (entities.EventInfo, error) {
	l := strings.Split(row, " - ")
	displaynr, err := strconv.Atoi(l[0])
	if err != nil {
		return entities.EventInfo{}, errors.New("Couldn't convert displaynr to int")
	}
	if l[1] == "" {
		return entities.EventInfo{}, errors.New("Event Name is empty")
	}
	return entities.EventInfo{
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

func parseName(row *colly.HTMLElement) (string, string) {
	nameArray := strings.SplitN(getOnlyChildText(row, "div.col-xs-11.col-sm-4 > a"), " ", 2)
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

func getClubFromStartOrResult(clubId uint, row *colly.HTMLElement) *entities.Club {
	var club entities.Club
	club.Id = clubId
	club.Name = row.ChildText("div.hidden-xs.col-sm-4 > a")
	flagLink := row.ChildAttr("img", "src")
	if flagLink != "" {
		club.Nationality.SetValid("https://myresults.eu/" + flagLink)
	}
	return &club
}

func getSwimmerFromStartOrResult(swimmerId uint, row *colly.HTMLElement) *entities.Swimmer {
	var swimmer entities.Swimmer
	clubLink := row.ChildAttrs("a", "href")[1]
	r := regexp.MustCompile("\\d+$")
	c, _ := strconv.Atoi(r.FindString(clubLink))
	clubId := uint(c)
	swimmer.Id = swimmerId
	swimmer.ClubId = clubId
	swimmer.Lastname, swimmer.Firstname = parseName(row)

	details := row.ChildText("div.col-xs-11.col-sm-4 > span")
	r = regexp.MustCompile("\\d+|[A-Z]")
	birthAndGender := r.FindAllString(details, -1)
	if len(birthAndGender) == 1 {
		swimmer.Gender = birthAndGender[0]
		swimmer.IsRelay = true
	} else if len(birthAndGender) == 2 {
		swimmer.BirthYear = null.IntFrom(entities.StringToInt(birthAndGender[0]))
		swimmer.Gender = birthAndGender[1]
		swimmer.IsRelay = false
	}
	return &swimmer
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
				heatId = heats.Add(&entities.Heat{
					EventId: eventId,
					HeatNr:  heatNr,
				})
			} else {
				// Start-Element
				swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
				clubLink := row.ChildAttr("div.hidden-xs.col-sm-4 > a", "href")
				r := regexp.MustCompile("\\d+$")
				swimmerId := entities.StringToUint(r.FindString(swimmerLink))
				clubId := entities.StringToUint(r.FindString(clubLink))

				lane := entities.StringToUint(row.ChildText("div.col-xs-1"))
				startTime, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
				s := entities.Start{
					HeatId:    heatId,
					SwimmerId: swimmerId,
					Lane:      lane,
				}
				if err == nil {
					s.Time.SetValid(startTime)
				}
				starts.Add(s)

				clubs.EnsureExists(clubId, row, getClubFromStartOrResult)
				swimmers.EnsureExists(swimmerId, row, getSwimmerFromStartOrResult)
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
				swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
				clubLink := row.ChildAttr("div.hidden-xs.col-sm-4 > a", "href")
				r := regexp.MustCompile("\\d+$")
				swimmerId := entities.StringToUint(r.FindString(swimmerLink))
				clubId := entities.StringToUint(r.FindString(clubLink))

				resultInfoString := row.ChildText("div.myresults_content_divtable_points")
				r = regexp.MustCompile("(?<timeToFirst>\\+\\d+\\.\\d+)|(?<reugeld>RG)|(?<finaPoints>\\d+)|(?<additionalInfo>[\\S]+$)")
				resultInfoMap := regex.EvalRegex(r, resultInfoString)

				dbResultId, ok := swimmerIdToDBResultId[swimmerId]
				if !ok {
					result := entities.Result{EventId: eventId, SwimmerId: swimmerId}

					finaPointsString, ok := resultInfoMap["finaPoints"]
					if ok {
						result.FinaPoints.SetValid(entities.StringToInt(finaPointsString))
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
						r = regexp.MustCompile("^RT \\+(\\d+.\\d+)")
						m := r.FindAllStringSubmatch(splits, -1)
						if len(m) > 0 {
							reactionTime, err := strconv.ParseFloat(m[0][1], 64)
							fmt.Println(reactionTime)
							if err == nil {
								result.ReactionTime.SetValid(reactionTime)
								splits = strings.Replace(splits, m[0][0]+" ", "", 1)
							}
						}
						fmt.Println(splits)
						result.Splits.SetValid(splits)
					}
					time, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
					if err == nil {
						result.Time.SetValid(time)
					}
					dbResultId = results.Add(&result)
					swimmerIdToDBResultId[swimmerId] = dbResultId
				}

				var ageclass entities.AgeClass
				timeToFirst, ok := resultInfoMap["timeToFirst"]
				if ok {
					ageclass.TimeToFirst.SetValid(timeToFirst)
				}

				position := row.ChildText("span.msecm-place")
				if position != "" {
					ageclass.Position.SetValid(entities.StringToInt(strings.Replace(position, ".", "", 1)))
				}
				ageclass.Name = ageClassName
				ageclass.ResultId = dbResultId

				ageclasses.Add(ageclass)

				clubs.EnsureExists(clubId, row, getClubFromStartOrResult)
				swimmers.EnsureExists(swimmerId, row, getSwimmerFromStartOrResult)
			}
		})
		wg.Done()
	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Results/" + UintToString(resultId))
}

func initSupabase() {
	var err error
	supabase, err = GetClient()
	if err != nil {
		fmt.Println("Couldn't initalize client: ", err)
	}
}

func UpdateSchedule(meetId uint, waitGroup *sync.WaitGroup) {
	initSupabase()
	log.Printf("Updating Schedule for: %d", meetId)
	c := colly.NewCollector()

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div.col-xs-12.col-md-12.myresults_content_divtable", func(e *colly.HTMLElement) {
		if waitGroup != nil {
			defer waitGroup.Done()
		}
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

		var sessionId uint
		var sessionIdx int

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				// Session-Item
				sessionInfo, _ := parseSessionInfo(row.Text)
				if scheduleUpToDate {
					sessionIdx = slices.IndexFunc(sessionsWithEvents, func(s entities.SessionWithEvents) bool {
						return s.SessionInfo == sessionInfo
					})
					sessionId = sessionsWithEvents[sessionIdx].Id
				} else {
					sessionId = sessions.Add(&entities.Session{
						Meetid:      meetId,
						SessionInfo: sessionInfo,
					})
				}
			} else {
				// Event-Item
				eventInfo, err := parseEventInfo(strings.Split(row.ChildText(".col-xs-6"), "\t")[0])
				if err != nil {
					return
				}
				var eventId uint
				if scheduleUpToDate {
					eventId = sessionsWithEvents[sessionIdx].Events[slices.IndexFunc(sessionsWithEvents[sessionIdx].Events, func(e entities.Event) bool {
						return e.EventInfo == eventInfo && e.SessionId == sessionId
					})].Id
				} else {
					eventId = events.Add(&entities.Event{
						SessionId: sessionId,
						EventInfo: eventInfo,
					})
				}

				status := status(row.ChildAttr("div.col-xs-1.text-center.myresults_content_divtable_left.hidden-xs > img", "src"))

				if status == EventStatusNoInfo {
					return
				}
				href := row.ChildAttr(".myresults_content_link.myresults_content_divtablecol", "href")
				r := regexp.MustCompile("\\d+$")
				startResultId := entities.StringToUint(r.FindString(href))
				wg.Add(1)
				go populateStarts(meetId, startResultId, eventId)
				if status == EventStatusFinished {
					wg.Add(1)
					go populateNewResults(meetId, startResultId, eventId)
				}
			}
		})

		wg.Wait()

	})

	c.Visit("https://myresults.eu/de-DE/Meets/Recent/" + UintToString(meetId) + "/Schedule")
}
