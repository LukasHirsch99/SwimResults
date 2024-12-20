package updateschedule

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/repository"
	"swimresults-backend/regex"
	"time"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)

func extractResultInfo(row *colly.HTMLElement) (repository.CreateResultParams, error) {
	result := repository.CreateResultParams{}

	r := regexp.MustCompile("\\d+$")
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerId, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		return result, err
	}

	result.Swimmerid = int32(swimmerId)

	// Points Holds, Finapoints, Penalty and Additionalinfo when result is invalid
	resultInfoString := row.ChildText("div.myresults_content_divtable_points")
	r = regexp.MustCompile("(?<penalty>RG)|(?<finaPoints>\\d+)|(?<additionalInfo>[\\S]+$)")
	resultInfoMap := regex.EvalRegex(r, resultInfoString)

	// Finapoints
	finaPointsString, ok := resultInfoMap["finaPoints"]
	if ok {
		finapoints, err := strconv.Atoi(finaPointsString)
		result.Finapoints = pgtype.Int4{Int32: int32(finapoints), Valid: err == nil}
	}

	// Invalid Result Info
	additionalInfo, ok := resultInfoMap["additionalInfo"]
	result.Additionalinfo = pgtype.Text{String: additionalInfo, Valid: ok}

	// Penalty Logic
	_, ok = resultInfoMap["penalty"]
	if ok && row.ChildText("div.myresults_content_divtable_points > span") == "" {
		result.Penalty = pgtype.Bool{Bool: true, Valid: true}
	} else {
		result.Penalty = pgtype.Bool{Bool: false, Valid: true}
	}

	// Time
	timeStr, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
	result.Time = pgtype.Time{Microseconds: timeStr.UnixMicro(), Valid: err == nil}

	// Splits
	splitsStr := row.ChildText("span.myresults_content_divtable_details:nth-child(1)")
	r = regexp.MustCompile("^RT \\+(\\d+.\\d+)(.*)")
	matches := r.FindAllStringSubmatch(splitsStr, -1)

	if len(matches) > 0 {
		reactionTime, err := strconv.ParseFloat(matches[0][1], 64)
		result.Reactiontime = pgtype.Float4{Float32: float32(reactionTime), Valid: err == nil}

		r = regexp.MustCompile("(\\d*)m: (\\d{2}:\\d{2},\\d{2})")
		splitMatches := r.FindAllStringSubmatch(matches[0][2], -1)
		if len(splitMatches) == 0 {
			return result, nil
		}

		splitsMap := make(map[int]time.Time)
		for _, split := range splitMatches {
			splitDist, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			splitsMap[splitDist], err = time.Parse("04:05,00", split[2])
			if err != nil {
				panic(err)
			}
		}

		result.Splits, err = json.Marshal(splitsMap)
		if err != nil {
			panic(err)
		}
	}

	return result, nil
}

func populateNewResults(meetId int32, resultId int, eventId int32) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#starts_content", func(e *colly.HTMLElement) {
		defer startResultWg.Done()
		ageclassCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if !strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				ageclassCnt++
			}
		})

		dbAgeclassCnt, err := repo.GetAgeclassCntByEvent(context.Background(), int32(eventId))
		if err != nil {
			panic(err)
		}

		if ageclassCnt == int(dbAgeclassCnt) {
			return
		}

		err = repo.DeleteAgeclassesByEvent(context.Background(), int32(eventId))
		err = repo.DeleteResultsByEvent(context.Background(), int32(eventId))
		err = repo.DeleteAgeclass_to_Results_ByEvent(context.Background(), int32(eventId))
		if err != nil {
			panic(err)
		}

		swimmerToResultid := make(map[int32]int32)
		var ageclassToResults []repository.CreateAgeclassToResultsParams
		var currentAgeclass repository.Ageclass

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Ageclass-Element
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				currentAgeclass, err = repo.GetAgeclassByName(context.Background(), strings.TrimSpace(row.Text))
				if err != nil {
					currentAgeclass, err = repo.CreateAgeclass(context.Background(), strings.TrimSpace(row.Text))
					if err != nil {
						panic(err)
					}
				}

				//
				// Result-Element
			} else {
				ensureSwimmerExists(row)
				result, err := extractResultInfo(row)
				if err != nil {
					panic(err)
				}
				resultId, resultExists := swimmerToResultid[result.Swimmerid]

				if !resultExists {
					resultId, err = repo.CreateResult(context.Background(), result)
					if err != nil {
						panic(err)
					}
					swimmerToResultid[result.Swimmerid] = resultId
				}

				ageclassToResults = append(ageclassToResults, repository.CreateAgeclassToResultsParams{
					Eventid:    int32(eventId),
					Ageclassid: currentAgeclass.ID,
					Resultid:   resultId,
				})
			}
		})

		_, err = repo.CreateAgeclassToResults(context.Background(), ageclassToResults)
		if err != nil {
			panic(err)
		}
	})

	c.Visit(fmt.Sprint("https://myresults.eu/de-DE/Meets/Recent/", meetId, "/Results/", resultId))
}
