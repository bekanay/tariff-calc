package service

import (
	"tariff-api/internal/model"
	"tariff-api/internal/repository"
)

type StationService struct {
	repo *repository.StationRepository
}

func NewStationService(repo *repository.StationRepository) *StationService {
	return &StationService{repo: repo}
}

func (s *StationService) GetStations() ([]model.Station, error) {
	stations, err := s.repo.GetStations()
	if err != nil {
		return nil, err
	}

	return stations, nil
}
