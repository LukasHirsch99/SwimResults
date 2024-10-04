package updateschedule

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"swimresults-backend/internal/database/models"
	"swimresults-backend/regex"
	"time"

	"github.com/gocolly/colly"
)

func extractResultInfo(ageclassId int, row *colly.HTMLElement) (models.Result, error) {
	result := models.Result{
		Ageclassid: ageclassId,
	}

	r := regexp.MustCompile("\\d+$")
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerId, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		return result, err
	}

	result.Swimmerid = swimmerId

	// Points Holds, Finapoints, Penalty and Additionalinfo when result is invalid
	resultInfoString := row.ChildText("div.myresults_content_divtable_points")
	r = regexp.MustCompile("(?<penalty>RG)|(?<finaPoints>\\d+)|(?<additionalInfo>[\\S]+$)")
	resultInfoMap := regex.EvalRegex(r, resultInfoString)

	// Finapoints
	finaPointsString, ok := resultInfoMap["finaPoints"]
	if ok {
		finapoints, err := strconv.Atoi(finaPointsString)
		result.Finapoints = sql.NullInt16{Int16: int16(finapoints), Valid: err == nil}
	}

	// Invalid Result Info
	additionalInfo, ok := resultInfoMap["additionalInfo"]
	result.Additionalinfo = sql.NullString{String: additionalInfo, Valid: ok}

	// Penalty Logic
	_, ok = resultInfoMap["penalty"]
	if ok && row.ChildText("div.myresults_content_divtable_points > span") == "" {
		result.Penalty = true
	} else {
		result.Penalty = false
	}

	// Time
	timeStr, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))
	result.Time = sql.NullTime{Time: timeStr, Valid: err == nil}

	// Splits
	splitsStr := row.ChildText("span.myresults_content_divtable_details:nth-child(1)")
	r = regexp.MustCompile("^RT \\+(\\d+.\\d+)(.*)")
	matches := r.FindAllStringSubmatch(splitsStr, -1)

	if len(matches) > 0 {
		reactionTime, err := strconv.ParseFloat(matches[0][1], 64)
		result.Reactiontime = sql.NullFloat64{Float64: reactionTime, Valid: err == nil}

		r = regexp.MustCompile("(\\d*)m: (\\d{2}:\\d{2},\\d{2})")
		splitMatches := r.FindAllStringSubmatch(matches[0][2], -1)
		if len(splitMatches) == 0 {
			return result, nil
		}

		result.Splits = make(map[int]time.Time)
		for _, split := range splitMatches {
			splitDist, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			result.Splits[splitDist], err = time.Parse("04:05,00", split[2])
      if err != nil {
        panic(err)
      }
		}
	}

	return result, nil
}

func populateNewResults(meetId int, resultId int, eventId int) {
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

		dbAgeclassCnt, err := repos.AgeclassRepository.CountForEvent(int(eventId))
		if err != nil {
			panic(err)
		}

		if ageclassCnt == dbAgeclassCnt {
			return
		}

		err = repos.AgeclassRepository.DeleteForEvent(int(eventId))
		if err != nil {
			panic(err)
		}

		var ageclass models.Ageclass
		var results []models.Result

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Ageclass-Element
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				ageclass = models.Ageclass{
					Name:    strings.TrimSpace(row.Text),
					Eventid: eventId,
				}
				err = repos.AgeclassRepository.Create(&ageclass)
				if err != nil {
					panic(err)
				}

				//
				// Result-Element
			} else {

				ensureSwimmerExists(row)

				result, err := extractResultInfo(ageclass.Id, row)
				if err != nil {
					panic(err)
				}

				results = append(results, result)
			}
		})
		err = repos.ResultRepository.CreateMany(results)
		if err != nil {
			panic(err)
		}
	})

	c.Visit(fmt.Sprint("https://myresults.eu/de-DE/Meets/Recent/", meetId, "/Results/", resultId))
}
