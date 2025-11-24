package service

import (
	"errors"
	"strings"
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
	station.Kod = strings.TrimSpace(station.Kod)
	station.Name = strings.TrimSpace(station.Name)
	if station.Kod == "" || station.Name == "" {
		return model.Station{}, ErrInvalidStationData
	}

	created, err := s.repo.AddStation(station)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateStation) {
			return model.Station{}, ErrDuplicateStation
		}
		return model.Station{}, err
	}
	return created, nil
}

var (
	ErrDuplicateStation   = errors.New("station already exists")
	ErrInvalidStationData = errors.New("invalid station data")
	ErrStationNotFound    = errors.New("station not found")
)

func (s *StationService) GetStation(kod string) (model.Station, error) {
	kod = strings.TrimSpace(kod)
	if kod == "" {
		return model.Station{}, ErrInvalidStationData
	}
	station, err := s.repo.GetStationByKod(kod)
	if err != nil {
		if errors.Is(err, repository.ErrStationNotFound) {
			return model.Station{}, ErrStationNotFound
		}
		return model.Station{}, err
	}
	return station, nil
}

func (s *StationService) UpdateStation(existingKod string, station model.Station) (model.Station, error) {
	existingKod = strings.TrimSpace(existingKod)
	station.Kod = strings.TrimSpace(station.Kod)
	station.Name = strings.TrimSpace(station.Name)
	if existingKod == "" || station.Kod == "" || station.Name == "" || station.Priznak == 0 {
		return model.Station{}, ErrInvalidStationData
	}

	updated, err := s.repo.UpdateStation(existingKod, station)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateStation):
			return model.Station{}, ErrDuplicateStation
		case errors.Is(err, repository.ErrStationNotFound):
			return model.Station{}, ErrStationNotFound
		default:
			return model.Station{}, err
		}
	}
	return updated, nil
}

func (s *StationService) DeleteStation(kod string) error {
	kod = strings.TrimSpace(kod)
	if kod == "" {
		return ErrInvalidStationData
	}
	if err := s.repo.DeleteStationByKod(kod); err != nil {
		if errors.Is(err, repository.ErrStationNotFound) {
			return ErrStationNotFound
		}
		return err
	}
	return nil
}
