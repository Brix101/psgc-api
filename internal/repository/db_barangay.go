package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Brix101/psgc-tool/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type dbBarangayRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewDBBarangay(conn *sql.DB) domain.BarangayRepository {
	tracer := otel.Tracer("db:sqlite3:barangay")

	return &dbBarangayRepository{conn: conn, tracer: tracer}
}

func (p *dbBarangayRepository) fetch(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]domain.Barangay, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying barangay")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var mLst []domain.Barangay
	for rows.Next() {
		var lst domain.Barangay
		if err := rows.Scan(
			&lst.PsgcCode,
			&lst.CityMuniCode,
			&lst.Name,
		); err != nil {
			return nil, err
		}
		mLst = append(mLst, lst)
	}
	return mLst, nil
}

func (p *dbBarangayRepository) paginatedQuery(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedBarangay, error) {
	queryParams := []interface{}{}
	query := `SELECT * FROM barangay`
	countQuery := `SELECT COUNT(*) FROM barangay`

	if params.Keyword != "" {
		query += `
			WHERE (
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
		return domain.PaginatedBarangay{}, err
	}

	totalItems := 0
	if err := p.conn.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
		return domain.PaginatedBarangay{}, err
	}

	totalPages := (totalItems + params.PerPage - 1) / params.PerPage

	if len(lst) == 0 {
		lst = []domain.Barangay{}
	}

	metaData := domain.MetaData{
		Page:       params.Page,
		TotalPages: totalPages,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		ItemCount:  len(lst),
	}

	res := domain.PaginatedBarangay{
		MetaData: metaData,
		Data:     lst,
	}

	return res, nil
}

func (p *dbBarangayRepository) GetAll(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedBarangay, error) {
	res, err := p.paginatedQuery(ctx, params)

	return res, err
}

func (p *dbBarangayRepository) GetById(
	ctx context.Context,
	psgcCode string,
) (domain.Barangay, error) {
	query := `SELECT * FROM barangay WHERE psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Barangay{}, err
	}

	if len(accs) == 0 {
		return domain.Barangay{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbBarangayRepository) Create(
	ctx context.Context,
	data *domain.Masterlist,
) error {
	query := `
		INSERT OR REPLACE INTO barangay (psgc_code, name, citmun_code)
		VALUES (?, ?, ?);`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	psgcCode := data.PsgcCode
	cityMuniCode := psgcCode[:7] + strings.Repeat("0", len(psgcCode)-7)

	_, err := p.conn.ExecContext(
		ctx,
		query,
		data.PsgcCode,
		data.Name,
		cityMuniCode,
	)
	if err != nil {
		span.SetStatus(codes.Error, "failed inserting barangay")
		span.RecordError(err)
		return err
	}

	return nil
}
