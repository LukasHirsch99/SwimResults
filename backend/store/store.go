package store

import (
	"swimresults-backend/database"
	"swimresults-backend/entities"
	sharedlist "swimresults-backend/sharedList"
)

var supabase, _ = database.GetClient()

var Meets = sharedlist.GetSharedListWithUniqueIdWithItems("meet", supabase.GetUpcomingMeets())

var Sessions = sharedlist.GetSharedListWithMaxId[*entities.Session]("session", supabase.GetMaxId("session"))
var Events = sharedlist.GetSharedListWithMaxId[*entities.Event]("event", supabase.GetMaxId("event"))
var Heats = sharedlist.GetSharedListWithMaxId[*entities.Heat]("heat", supabase.GetMaxId("heat"))
var Results = sharedlist.GetSharedListWithMaxId[*entities.Result]("result", supabase.GetMaxId("result"))

var Clubs = sharedlist.GetSharedListWithUniqueId[*entities.Club]("club", supabase.GetIds("club"))
var Swimmers = sharedlist.GetSharedListWithUniqueId[*entities.Swimmer]("swimmer", supabase.GetIds("swimmer"))

var Starts = sharedlist.GetSharedList[entities.Start]("start")
var Ageclasses = sharedlist.GetSharedList[entities.AgeClass]("ageclass")
