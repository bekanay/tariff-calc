package repository

import (
	"database/sql"
	"tariff-api/internal/model"
)

type StationRepository struct {
	db *sql.DB
}

func NewStationRepository(db *sql.DB) *StationRepository {
	return &StationRepository{db: db}
}

func (repo *StationRepository) GetStations() ([]model.Station, error) {
	rows, err := repo.db.Query(`
		SELECT id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stations := make([]model.Station, 0)
	for rows.Next() {
		var (
			station   model.Station
			paragraph sql.NullString
		)
		if err := rows.Scan(&station.ID, &station.Kod, &station.Name, &station.Priznak, &paragraph); err != nil {
			return nil, err
		}
		station.Paragraph = paragraph.String
		stations = append(stations, station)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return stations, nil
}
