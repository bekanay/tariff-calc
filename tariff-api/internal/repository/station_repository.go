package repository

import (
	"database/sql"
	"fmt"
	"tariff-api/internal/model"
)

type StationRepository struct {
	db *sql.DB
}

func NewStationRepository(db *sql.DB) *StationRepository {
	return &StationRepository{db: db}
}

func (repo *StationRepository) GetStations(filters model.Filters) ([]model.Station, model.Metadata, error) {
	query := fmt.Sprintf(`
		SELECT count(*) OVER(), id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan
		ORDER BY %s %s, id ASC
		LIMIT $1 OFFSET $2`, filters.SortColumn(), filters.SortDirection())
	rows, err := repo.db.Query(query, filters.Limit(), filters.Offset())
	if err != nil {
		return nil, model.Metadata{}, err
	}
	defer rows.Close()

	stations := make([]model.Station, 0)
	var totalRecords int
	for rows.Next() {
		var (
			station   model.Station
			paragraph sql.NullString
			total     int
		)
		if err := rows.Scan(&total, &station.ID, &station.Kod, &station.Name, &station.Priznak, &paragraph); err != nil {
			return nil, model.Metadata{}, err
		}
		station.Paragraph = paragraph.String
		stations = append(stations, station)
		totalRecords = total
	}
	if err := rows.Err(); err != nil {
		return nil, model.Metadata{}, err
	}

	metadata := model.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return stations, metadata, nil
}
