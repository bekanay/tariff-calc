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

func (s *StationService) GetStations(filters model.Filters) ([]model.Station, model.Metadata, error) {
	stations, metadata, err := s.repo.GetStations(filters)
	if err != nil {
		return nil, model.Metadata{}, err
	}

	return stations, metadata, nil
}

func (s *StationService) AddStation(station model.Station) (model.Station, error) {
	created, err := s.repo.AddStation(station)
	if err != nil {
		return model.Station{}, err
	}
	return created, nil
}
