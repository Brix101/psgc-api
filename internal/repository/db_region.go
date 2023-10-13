package repository

import (
	"context"
	"database/sql"

	"github.com/Brix101/psgc-tool/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type dbRegionRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewDBRegion(conn *sql.DB) domain.RegionRepository {
	tracer := otel.Tracer("db:sqlite3:region")

	return &dbRegionRepository{conn: conn, tracer: tracer}
}

func (p *dbRegionRepository) fetch(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]domain.Region, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying region")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var mLst []domain.Region
	for rows.Next() {
		var lst domain.Region
		if err := rows.Scan(
			&lst.PsgcCode,
			&lst.Name,
		); err != nil {
			return nil, err
		}
		mLst = append(mLst, lst)
	}
	return mLst, nil
}

func (p *dbRegionRepository) paginatedQuery(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedRegion, error) {
	queryParams := []interface{}{}
	query := `SELECT * FROM region`
	countQuery := `SELECT COUNT(*) FROM region`

	if params.Filter != "" {
		query += `
			WHERE (
                LOWER(psgc_code) LIKE '%' || LOWER($1) || '%' OR
                LOWER(name) LIKE '%' || LOWER($1) || '%' 
            )
        `
		queryParams = append(queryParams, params.Filter)
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
		return domain.PaginatedRegion{}, err
	}

	totalItems := 0
	p.conn.QueryRowContext(ctx, countQuery).
		Scan(&totalItems)

	totalPages := (totalItems + params.PerPage - 1) / params.PerPage

	if len(lst) == 0 {
		lst = []domain.Region{}
	}

	metaData := domain.MetaData{
		Page:       params.Page,
		TotalPages: totalPages,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		ItemCount:  len(lst),
	}

	res := domain.PaginatedRegion{
		MetaData: metaData,
		Data:     lst,
	}

	return res, nil
}

func (p *dbRegionRepository) GetAll(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedRegion, error) {
	res, err := p.paginatedQuery(ctx, params)

	return res, err
}

func (p *dbRegionRepository) GetById(
	ctx context.Context,
	psgcCode string,
) (domain.Region, error) {
	query := `SELECT * FROM region WHERE psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Region{}, err
	}

	if len(accs) == 0 {
		return domain.Region{}, domain.ErrNotFound
	}
	return accs[0], nil
}


func (p *dbRegionRepository) Create(
	ctx context.Context,
	data *domain.Masterlist,
) error {
	query := `
		INSERT OR REPLACE INTO region (psgc_code, name)
		VALUES (?, ?);`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	_, err := p.conn.ExecContext(
		ctx,
		query,
		data.PsgcCode,
		data.Name,
	)
	if err != nil {
		span.SetStatus(codes.Error, "failed inserting region")
		span.RecordError(err)
		return err
	}

	return nil
}