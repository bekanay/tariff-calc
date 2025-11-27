package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tariff-api/internal/model"
)

type StationRepository struct {
	db *sql.DB
}

func NewStationRepository(db *sql.DB) *StationRepository {
	return &StationRepository{db: db}
}

var ErrDuplicateStation = errors.New("station already exists")
var ErrStationNotFound = errors.New("station not found")

func (repo *StationRepository) AddStation(station model.Station) (model.Station, error) {
	var (
		paragraph sql.NullString
		created   model.Station
	)
	if station.Paragraph != "" {
		paragraph = sql.NullString{String: station.Paragraph, Valid: true}
	}

	err := repo.db.QueryRow(`
		SELECT id, stan_kod, stan_name, stan_priznak, paragraph
		FROM FINAL TABLE (
			INSERT INTO stan (stan_kod, stan_name, stan_priznak, paragraph)
			VALUES (?, ?, ?, ?)
		)
	`, station.Kod, station.Name, station.Priznak, paragraph).Scan(
		&created.ID, &created.Kod, &created.Name, &created.Priznak, &paragraph,
	)
	if err != nil {
		if isDB2UniqueViolation(err) {
			return model.Station{}, ErrDuplicateStation
		}
		return model.Station{}, err
	}

	created.Paragraph = paragraph.String
	return created, nil
}

func (repo *StationRepository) GetStations(filters model.Filters) ([]model.Station, model.Metadata, error) {
	var (
		args    []any
		clauses []string
	)

	if filters.Name != "" {
		clauses = append(clauses, "LOWER(stan_name) LIKE LOWER(?)")
		args = append(args, "%"+filters.Name+"%")
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT count(*) OVER(), id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan`)
	if len(clauses) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(clauses, " AND "))
	}
	queryBuilder.WriteString(fmt.Sprintf(`
		ORDER BY %s %s, id ASC
		OFFSET ? ROWS FETCH NEXT ? ROWS ONLY`, filters.SortColumn(), filters.SortDirection()))

	args = append(args, filters.Offset(), filters.Limit())

	rows, err := repo.db.Query(queryBuilder.String(), args...)
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

func isDB2UniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToUpper(err.Error())
	// DB2 unique constraint violation typically returns SQLSTATE 23505 or SQL0803N.
	return strings.Contains(msg, "SQLSTATE=23505") || strings.Contains(msg, "SQL0803")
}

func (repo *StationRepository) GetStationByKod(kod string) (model.Station, error) {
	var (
		station   model.Station
		paragraph sql.NullString
	)
	err := repo.db.QueryRow(`
		SELECT id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan
		WHERE stan_kod = ?
	`, kod).Scan(&station.ID, &station.Kod, &station.Name, &station.Priznak, &paragraph)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Station{}, ErrStationNotFound
		}
		return model.Station{}, err
	}
	station.Paragraph = paragraph.String
	return station, nil
}

func (repo *StationRepository) UpdateStation(existingKod string, station model.Station) (model.Station, error) {
	var paragraph sql.NullString
	if station.Paragraph != "" {
		paragraph = sql.NullString{String: station.Paragraph, Valid: true}
	}

	err := repo.db.QueryRow(`
		SELECT id, stan_kod, stan_name, stan_priznak, paragraph
		FROM FINAL TABLE (
			UPDATE stan
			SET stan_kod = ?, stan_name = ?, stan_priznak = ?, paragraph = ?
			WHERE stan_kod = ?
		)
	`, station.Kod, station.Name, station.Priznak, paragraph, existingKod).Scan(
		&station.ID, &station.Kod, &station.Name, &station.Priznak, &paragraph,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Station{}, ErrStationNotFound
		}
		if isDB2UniqueViolation(err) {
			return model.Station{}, ErrDuplicateStation
		}
		return model.Station{}, err
	}
	station.Paragraph = paragraph.String
	return station, nil
}

func (repo *StationRepository) DeleteStationByKod(kod string) error {
	res, err := repo.db.Exec(`DELETE FROM stan WHERE stan_kod = ?`, kod)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrStationNotFound
	}
	return nil
}
