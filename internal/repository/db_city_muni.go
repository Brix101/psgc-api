package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Brix101/psgc-tool/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type dbCityMuniRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewDBCityMuni(conn *sql.DB) domain.CityMuniRepository {
	tracer := otel.Tracer("db:sqlite3:cityMuni")

	return &dbCityMuniRepository{conn: conn, tracer: tracer}
}

func (p *dbCityMuniRepository) fetch(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]domain.CityMuni, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying cityMuni")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var mLst []domain.CityMuni
	for rows.Next() {
		var lst domain.CityMuni
		if err := rows.Scan(
			&lst.PsgcCode,
			&lst.ProvCode,
			&lst.Name,
			&lst.Level,
		); err != nil {
			return nil, err
		}
		mLst = append(mLst, lst)
	}
	return mLst, nil
}

func (p *dbCityMuniRepository) paginatedQuery(
	ctx context.Context,
	level string,
	params domain.PaginationParams,
) (domain.PaginatedCityMuni, error) {
	queryParams := []interface{}{}
	query := `SELECT * FROM city_muni`
	countQuery := `SELECT COUNT(*) FROM city_muni`

	if level != "" {
		query += fmt.Sprintf(" WHERE level = '%s' AND", level)
		countQuery += fmt.Sprintf(" WHERE level = '%s'", level)
	} else {
		query += " WHERE "
	}

	if params.Keyword != "" {
		query += `
			(
                LOWER(psgc_code) LIKE '%' || LOWER($1) || '%' OR
                LOWER(name) LIKE '%' || LOWER($1) || '%' 
            )
        `
		queryParams = append(queryParams, params.Keyword)
	}

	// Add sorting by name in ascending order.
	query += `
        ORDER BY name ASC
        LIMIT $2
        OFFSET $3
    `

	queryParams = append(queryParams, params.PerPage, (params.Page-1)*params.PerPage)

	// Execute the query with appropriate parameters.
	lst, err := p.fetch(ctx, query, queryParams...)
	if err != nil {
		return domain.PaginatedCityMuni{}, err
	}

	totalItems := 0
	if err := p.conn.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
		return domain.PaginatedCityMuni{}, err
	}

	totalPages := (totalItems + params.PerPage - 1) / params.PerPage

	if len(lst) == 0 {
		lst = []domain.CityMuni{}
	}

	metaData := domain.MetaData{
		Page:       params.Page,
		TotalPages: totalPages,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		ItemCount:  len(lst),
	}

	res := domain.PaginatedCityMuni{
		MetaData: metaData,
		Data:     lst,
	}

	return res, nil
}

func (p *dbCityMuniRepository) GetAll(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedCityMuni, error) {
	res, err := p.paginatedQuery(ctx, "", params)

	return res, err
}

func (p *dbCityMuniRepository) GetById(
	ctx context.Context,
	psgcCode string,
) (domain.CityMuni, error) {
	query := `SELECT * FROM city_muni WHERE psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.CityMuni{}, err
	}

	if len(accs) == 0 {
		return domain.CityMuni{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbCityMuniRepository) GetAllCity(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedCityMuni, error) {
	res, err := p.paginatedQuery(ctx, "City", params)

	return res, err
}

func (p *dbCityMuniRepository) GetCityById(
	ctx context.Context,
	psgcCode string,
) (domain.CityMuni, error) {
	query := `SELECT * FROM city_muni WHERE level = 'City' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.CityMuni{}, err
	}

	if len(accs) == 0 {
		return domain.CityMuni{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbCityMuniRepository) GetAllMunicipality(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedCityMuni, error) {
	res, err := p.paginatedQuery(ctx, "Mun", params)

	return res, err
}

func (p *dbCityMuniRepository) GetMunicipalityById(
	ctx context.Context,
	psgcCode string,
) (domain.CityMuni, error) {
	query := `SELECT * FROM city_muni WHERE level = 'Mun' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.CityMuni{}, err
	}

	if len(accs) == 0 {
		return domain.CityMuni{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbCityMuniRepository) Create(
	ctx context.Context,
	data *domain.Masterlist,
) error {
	query := `
		INSERT OR REPLACE INTO city_muni (psgc_code, name, level, prov_code)
		VALUES (?, ?, ?, ?);`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	psgcCode := data.PsgcCode
	provCode := psgcCode[:5] + strings.Repeat("0", len(psgcCode)-5)

	_, err := p.conn.ExecContext(
		ctx,
		query,
		data.PsgcCode,
		data.Name,
		data.Level,
		provCode,
	)
	if err != nil {
		span.SetStatus(codes.Error, "failed inserting cityMuni")
		span.RecordError(err)
		return err
	}

	return nil
}
