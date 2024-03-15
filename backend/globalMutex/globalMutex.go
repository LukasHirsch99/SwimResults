package globalMutex

import (
	"fmt"
	"slices"
	"swimresults-backend/database"
	"sync"

	"github.com/gocolly/colly"
)

var lock = &sync.Mutex{}
var supabase *database.Supabase

var swimmerIds []uint
var clubIds []uint
var maxIds *database.MaxIds

var createdSwimmers []database.Swimmer
var createdClubs []database.Club
var createdSessions []database.Session
var createdEvents []database.Event
var createdHeats []database.Heat
var createdStarts []database.Start
var createdResults []database.Result
var createdAgeclasses []database.AgeClass

func Get() *sync.Mutex {
	return lock
}

func AddSession(session database.Session) uint {
	lock.Lock()
	defer lock.Unlock()
	if maxIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		maxIds, err = supabase.GetMaxIds()
		if err != nil {
			panic("GetMaxIds failed")
		}
	}

	maxIds.MaxSessionId++
	session.Id = maxIds.MaxSessionId
	createdSessions = append(createdSessions, session)
	return session.Id
}

func AddEvent(event database.Event) uint {
	lock.Lock()
	defer lock.Unlock()

	if maxIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		maxIds, err = supabase.GetMaxIds()
		if err != nil {
			panic("GetMaxIds failed")
		}
	}
	maxIds.MaxEventId++
	event.Id = maxIds.MaxEventId
	createdEvents = append(createdEvents, event)
	return event.Id
}

func AddHeat(heat database.Heat) uint {
	lock.Lock()
	defer lock.Unlock()

	if maxIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		maxIds, err = supabase.GetMaxIds()
		if err != nil {
			panic("GetMaxIds failed")
		}
	}
	maxIds.MaxHeatId++
	heat.Id = maxIds.MaxHeatId
	createdHeats = append(createdHeats, heat)
	return heat.Id
}

func AddStart(start database.Start) {
	lock.Lock()
	defer lock.Unlock()

	createdStarts = append(createdStarts, start)
}

func AddResult(result database.Result) uint {
	lock.Lock()
	defer lock.Unlock()

	if maxIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		maxIds, err = supabase.GetMaxIds()
		if err != nil {
			panic("GetMaxIds failed")
		}
	}
	maxIds.MaxResultId++
	result.Id = maxIds.MaxResultId
	createdResults = append(createdResults, result)
	return result.Id
}

func AddAgeclass(ageclass database.AgeClass) {
	lock.Lock()
	defer lock.Unlock()

	createdAgeclasses = append(createdAgeclasses, ageclass)
}

func EnsureClubExists(clubId uint, row *colly.HTMLElement, creator func(uint, *colly.HTMLElement) database.Club) {
	lock.Lock()
	defer lock.Unlock()

	if clubIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		clubIds, err = supabase.GetClubIds()
		if err != nil {
			panic("GetClubIds failed")
		}
	}

	if slices.Contains(clubIds, clubId) {
		return
	}

	club := creator(clubId, row)
	clubIds = append(clubIds, clubId)
	createdClubs = append(createdClubs, club)
}

func EnsureSwimmerExists(swimmerId uint, row *colly.HTMLElement, creator func(uint, *colly.HTMLElement) database.Swimmer) {
	lock.Lock()
	defer lock.Unlock()

	if swimmerIds == nil {
		var err error
		supabase, err = database.GetClient()
		if err != nil {
			panic("GetClient failed")
		}
		swimmerIds, err = supabase.GetSwimmerIds()
		if err != nil {
			panic("GetSwimmerIds failed")
		}
	}

	if slices.Contains(swimmerIds, swimmerId) {
		return
	}

	swimmer := creator(swimmerId, row)
	swimmerIds = append(swimmerIds, swimmerId)
	createdSwimmers = append(createdSwimmers, swimmer)
}

func SaveAll() {
  SaveCreatedClubs()
  SaveCreatedSwimmers()
  SaveCreatedSessions()
  SaveCreatedEvents()
  SaveCreatedHeats()
  SaveCreatedStarts()
  SaveCreatedResults()
  SaveCreatedAgeclasses()
}

func SaveCreatedClubs() {
	if len(createdClubs) == 0 {
		return
	}
	err := supabase.InsertInto("club", createdClubs)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d clubs\n", len(createdClubs))
	clear(createdClubs)
}

func SaveCreatedSwimmers() {
	if len(createdSwimmers) == 0 {
		return
	}
	err := supabase.InsertInto("swimmer", createdSwimmers)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d swimmers\n", len(createdSwimmers))
	clear(createdSwimmers)
}

func SaveCreatedSessions() {
	if len(createdSessions) == 0 {
		return
	}
	err := supabase.InsertInto("session", createdSessions)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d sessions\n", len(createdSessions))
	clear(createdSessions)
}

func SaveCreatedEvents() {
	if len(createdEvents) == 0 {
		return
	}
	err := supabase.InsertInto("event", createdEvents)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d events\n", len(createdEvents))
	clear(createdEvents)
}

func SaveCreatedHeats() {
	if len(createdHeats) == 0 {
		return
	}
	err := supabase.InsertInto("heat", createdHeats)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d heats\n", len(createdHeats))
	clear(createdHeats)
}

func SaveCreatedStarts() {
	if len(createdStarts) == 0 {
		return
	}
	err := supabase.InsertInto("start", createdStarts)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d starts\n", len(createdStarts))
	clear(createdStarts)
}

func SaveCreatedResults() {
	if len(createdResults) == 0 {
		return
	}
	err := supabase.InsertInto("result", createdResults)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d results\n", len(createdResults))
	clear(createdResults)
}

func SaveCreatedAgeclasses() {
	if len(createdAgeclasses) == 0 {
		return
	}
	err := supabase.InsertInto("ageclass", createdAgeclasses)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created %d ageclass\n", len(createdAgeclasses))
	clear(createdAgeclasses)
}
