package dtos

import "swimresults-backend/internal/repository"


type StartWithSwimmer struct {
  repository.Start
  Swimmer repository.Swimmer
}
