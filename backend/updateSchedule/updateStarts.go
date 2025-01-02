package updateschedule

import (
	"context"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"swimresults-backend/internal/repository"

	"github.com/gocolly/colly"
	"github.com/jackc/pgx/v5/pgtype"
)

func extractStartInfo(heatId int32, row *colly.HTMLElement) (repository.CreateStartsParams, error) {
	r := regexp.MustCompile("\\d+$")
	swimmerLink := row.ChildAttr("div.col-xs-11.col-sm-4 > a", "href")
	swimmerId, err := strconv.Atoi(r.FindString(swimmerLink))
	if err != nil {
		return repository.CreateStartsParams{}, err
	}

	lane, err := strconv.Atoi(row.ChildText("div.col-xs-1"))
	if err != nil {
		return repository.CreateStartsParams{}, err
	}

	startTime, err := parseTime(row.ChildText("div.hidden-xs.col-sm-2.col-md-1.text-right.myresults_content_divtable_right"))

	start := repository.CreateStartsParams{
		Heatid:    heatId,
		Swimmerid: int32(swimmerId),
		Lane:      int32(lane),
		Time:      pgtype.Time{Microseconds: startTime.UnixMicro(), Valid: err == nil},
	}

	return start, nil
}

func populateStarts(meetId int32, startId int, eventId int32) {
	c := colly.NewCollector()
	c.OnError(func(_ *colly.Response, err error) {
    logger.Error("colly error", "error", err)
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

		dbHeatCnt, err := repo.GetHeatCntForEvent(context.Background(), eventId)
		if err != nil {
			panic(err)
		}
		dbStartCnt, err := repo.GetStartCntForEvent(context.Background(), eventId)
		if err != nil {
			panic(err)
		}

		if startCnt == int(dbStartCnt) && heatCnt == int(dbHeatCnt) {
			return
		}

		err = repo.DeleteHeatsForEvent(context.Background(), eventId)
		if err != nil {
			panic(err)
		}

    var currentHeatId int32
    var heatNr int32 = 0
		var starts []repository.CreateStartsParams

		e.ForEach(".myresults_content_divtablerow", func(_ int, row *colly.HTMLElement) {
			// Heat-Element
			if strings.Contains(row.Attr("class"), "myresults_content_divtablerow_header") {
				heatNr++
				currentHeatId, err = repo.CreateHeat(context.Background(), repository.CreateHeatParams{
					Eventid: eventId,
					Heatnr:  heatNr,
				})
				if err != nil {
					panic(err)
				}

				// Start-Element
			} else {
				clubId, swimmerId := getClubAndSwimmerIdFromRow(row)

				if !slices.Contains(swimmerIds, swimmerId) {
					if !slices.Contains(clubIds, clubId) {
						c := CreateClubParamsFromStartOrResult(clubId, row)
						err = repo.CreateClub(context.Background(), c)
						if err != nil {
							panic(err)
						}
						clubIds = append(clubIds, clubId)
					}
					s := CreateSwimmerParamsFromStartOrResult(swimmerId, row)
					err = repo.CreateSwimmer(context.Background(), s)
					swimmerIds = append(swimmerIds, swimmerId)
				}

				start, err := extractStartInfo(currentHeatId, row)
				if err != nil {
					panic(err)
				}

				starts = append(starts, start)
			}
		})

		_, err = repo.CreateStarts(context.Background(), starts)
		if err != nil {
      fmt.Println(starts)
			panic(err)
		}
	})

	err := c.Visit(fmt.Sprintf("https://myresults.eu/de-DE/Meets/Recent/%d/Starts/%d", meetId, startId))
	if err != nil {
		panic(err)
	}
	return
}
