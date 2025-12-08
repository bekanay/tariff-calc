package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"tariff-api/internal/model"

	"github.com/jackc/pgconn"
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
		INSERT INTO stan (stan_kod, stan_name, stan_priznak, paragraph)
		VALUES ($1, $2, $3, $4)
		RETURNING id, stan_kod, stan_name, stan_priznak, paragraph
	`, station.Kod, station.Name, station.Priznak, paragraph).Scan(
		&created.ID, &created.Kod, &created.Name, &created.Priznak, &paragraph,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return model.Station{}, ErrDuplicateStation
		}
		return model.Station{}, err
	}

	created.Paragraph = paragraph.String
	return created, nil
}

func (repo *StationRepository) GetStations(filters model.Filters) ([]model.Station, model.Metadata, error) {
	var (
		args      []any
		clauses   []string
		argNumber = 1
	)

	if filters.Name != "" {
		clauses = append(clauses, fmt.Sprintf("LOWER(stan_name) LIKE LOWER($%d)", argNumber))
		args = append(args, "%"+filters.Name+"%")
		argNumber++
	}

	if filters.Paragraph != "" {
		clauses = append(clauses, fmt.Sprintf("LOWER(paragraph) LIKE LOWER($%d)", argNumber))
		args = append(args, "%"+filters.Paragraph+"%")
		argNumber++
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(`
		SELECT count(*) OVER(), id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan`)
	if len(clauses) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(clauses, " OR "))
	}
	queryBuilder.WriteString(fmt.Sprintf(`
		ORDER BY %s %s, id ASC
		LIMIT $%d OFFSET $%d`, filters.SortColumn(), filters.SortDirection(), argNumber, argNumber+1))

	args = append(args, filters.Limit(), filters.Offset())

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

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key value violates unique constraint")
}

func (repo *StationRepository) GetStationByKod(kod string) (model.Station, error) {
	var (
		station   model.Station
		paragraph sql.NullString
	)
	err := repo.db.QueryRow(`
		SELECT id, stan_kod, stan_name, stan_priznak, paragraph
		FROM stan
		WHERE stan_kod = $1
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
		UPDATE stan
		SET stan_kod = $1, stan_name = $2, stan_priznak = $3, paragraph = $4
		WHERE stan_kod = $5
		RETURNING id, stan_kod, stan_name, stan_priznak, paragraph
	`, station.Kod, station.Name, station.Priznak, paragraph, existingKod).Scan(
		&station.ID, &station.Kod, &station.Name, &station.Priznak, &paragraph,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Station{}, ErrStationNotFound
		}
		if isUniqueViolation(err) {
			return model.Station{}, ErrDuplicateStation
		}
		return model.Station{}, err
	}
	station.Paragraph = paragraph.String
	return station, nil
}

func (repo *StationRepository) DeleteStationByKod(kod string) error {
	res, err := repo.db.Exec(`DELETE FROM stan WHERE stan_kod = $1`, kod)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrStationNotFound
	}
	return nil
}
