package dtos

import "swimresults-backend/internal/repository"


type HeatWithStarts struct {
	repository.Heat
	Starts []StartWithSwimmer
}
