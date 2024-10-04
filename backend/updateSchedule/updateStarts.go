package updateschedule

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"swimresults-backend/internal/database/models"

	"github.com/gocolly/colly"
)

func extractStartInfo(heatId int, row *colly.HTMLElement) (models.Start, error) {
	r := regexp.MustCompile("\\d+$")
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerId, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		return models.Start{}, err
	}

	lane, err := strconv.Atoi(row.ChildText("div.col-xs-1"))
	if err != nil {
		return models.Start{}, err
	}

	startTime, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))

	start := models.Start{
		Heatid:    heatId,
		Swimmerid: swimmerId,
		Lane:      lane,
		Time:      sql.NullTime{Time: startTime, Valid: err == nil},
	}
	return start, nil
}

func populateStarts(meetId int, startId int, eventId int) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong: ", err)
	})

	c.OnHTML("div#starts_content", func(e *colly.HTMLElement) {
		defer startResultWg.Done()
		heatCnt := 0
		startCnt := 0

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				heatCnt++
			} else {
				startCnt++
			}
		})

		dbHeatCnt, err := repos.HeatRepository.CountForEvent(eventId)
		if err != nil {
			panic(err)
		}
		dbStartCnt, err := repos.StartRepository.CountForEvent(eventId)
		if err != nil {
			panic(err)
		}

		if startCnt == dbStartCnt && heatCnt == dbHeatCnt {
			return
		}

		err = repos.HeatRepository.DeleteForEvent(eventId)
		if err != nil {
			panic(err)
		}

		heatNr := 0
		var heat models.Heat
		var starts []models.Start

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Heat-Element
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				heatNr++
				heat = models.Heat{
					Eventid: eventId,
					Heatnr:  heatNr,
				}
				err := repos.HeatRepository.Create(&heat)
				if err != nil {
					panic(err)
				}

				// Start-Element
			} else {
				clubId, swimmerId := getClubAndSwimmerIdFromRow(row)

				if !slices.Contains(swimmerIds, swimmerId) {
					if !slices.Contains(clubIds, clubId) {
						err = repos.ClubRepository.Create(extractClubFromStartOrResult(clubId, row))
						if err != nil {
							panic(err)
						}
						clubIds = append(clubIds, clubId)
					}
					err = repos.SwimmerRepository.Create(extractSwimmerFromStartOrResult(swimmerId, row))
					if err != nil {
						panic(err)
					}
					swimmerIds = append(swimmerIds, swimmerId)
				}

				start, err := extractStartInfo(heat.Id, row)
				if err != nil {
					panic(err)
				}

				starts = append(starts, start)

				// err = repos.StartRepository.Create(&start)
				// if err != nil {
				// 	panic(err)
				// }
				// ensureSwimmerExists(row)
			}
		})

		err = repos.StartRepository.CreateMany(starts)
		if err != nil {
			panic(err)
		}
	})

	c.Visit(fmt.Sprintf("https://myresults.eu/de-DE/Meets/Recent/%d/Starts/%d", meetId, startId))
	return
}
