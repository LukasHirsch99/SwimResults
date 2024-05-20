package store

import (
	"swimresults-backend/database"
	"swimresults-backend/entities"
	sharedlist "swimresults-backend/sharedList"
)

var db = database.GetClient()

var Meets = sharedlist.GetSharedListWithUniqueIdWithItems[*entities.Meet]("meet", db.GetUpcomingMeets())

var Sessions = sharedlist.GetSharedListWithMaxId[*entities.Session]("session", db.GetMaxId("session"))
var Events = sharedlist.GetSharedListWithMaxId[*entities.Event]("event", db.GetMaxId("event"))
var Heats = sharedlist.GetSharedListWithMaxId[*entities.Heat]("heat", db.GetMaxId("heat"))
var Results = sharedlist.GetSharedListWithMaxId[*entities.Result]("result", db.GetMaxId("result"))

var Clubs = sharedlist.GetSharedListWithUniqueId[*entities.Club]("club", db.GetIds("club"))
var Swimmers = sharedlist.GetSharedListWithUniqueId[*entities.Swimmer]("swimmer", db.GetIds("swimmer"))

var Starts = sharedlist.GetSharedList[entities.Start]("start")
var Ageclasses = sharedlist.GetSharedList[entities.AgeClass]("ageclass")
