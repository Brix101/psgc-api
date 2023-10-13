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

type dbProvinceRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewDBProvince(conn *sql.DB) domain.ProvinceRepository {
	tracer := otel.Tracer("db:sqlite3:province")

	return &dbProvinceRepository{conn: conn, tracer: tracer}
}

func (p *dbProvinceRepository) fetch(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]domain.Province, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying province")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var mLst []domain.Province
	for rows.Next() {
		var lst domain.Province
		if err := rows.Scan(
			&lst.PsgcCode,
			&lst.RegCode,
			&lst.Name,
		); err != nil {
			return nil, err
		}
		mLst = append(mLst, lst)
	}
	return mLst, nil
}

func (p *dbProvinceRepository) paginatedQuery(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedProvince, error) {
	queryParams := []interface{}{}
	query := `SELECT * FROM province`
	countQuery := `SELECT COUNT(*) FROM province`

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
		return domain.PaginatedProvince{}, err
	}

	totalItems := 0
	if err := p.conn.QueryRowContext(ctx, countQuery).Scan(&totalItems); err != nil {
		return domain.PaginatedProvince{}, err
	}

	totalPages := (totalItems + params.PerPage - 1) / params.PerPage

	if len(lst) == 0 {
		lst = []domain.Province{}
	}

	metaData := domain.MetaData{
		Page:       params.Page,
		TotalPages: totalPages,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		ItemCount:  len(lst),
	}

	res := domain.PaginatedProvince{
		MetaData: metaData,
		Data:     lst,
	}

	return res, nil
}

func (p *dbProvinceRepository) GetAll(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedProvince, error) {
	res, err := p.paginatedQuery(ctx, params)

	return res, err
}

func (p *dbProvinceRepository) GetById(
	ctx context.Context,
	psgcCode string,
) (domain.Province, error) {
	query := `SELECT * FROM province WHERE psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Province{}, err
	}

	if len(accs) == 0 {
		return domain.Province{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbProvinceRepository) Create(
	ctx context.Context,
	data *domain.Masterlist,
) error {
	query := `
		INSERT OR REPLACE INTO province (psgc_code, name, reg_code)
		VALUES (?, ?, ?);`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	psgcCode := data.PsgcCode
	regCode := psgcCode[:2] + strings.Repeat("0", len(psgcCode)-2)

	_, err := p.conn.ExecContext(
		ctx,
		query,
		data.PsgcCode,
		data.Name,
		regCode,
	)
	if err != nil {
		span.SetStatus(codes.Error, "failed inserting province")
		span.RecordError(err)
		return err
	}

	return nil
}
