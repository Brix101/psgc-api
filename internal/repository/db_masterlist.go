package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Brix101/psgc-tool/internal/domain"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type dbMasterlistRepository struct {
	conn   Connection
	tracer trace.Tracer
}

func NewDBMasterlist(conn *sql.DB) domain.MasterlistRepository {
	tracer := otel.Tracer("db:sqlite3:masterlist")

	return &dbMasterlistRepository{conn: conn, tracer: tracer}
}

func (p *dbMasterlistRepository) fetch(
	ctx context.Context,
	query string,
	args ...interface{},
) ([]domain.Masterlist, error) {
	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	rows, err := p.conn.QueryContext(ctx, query, args...)
	if err != nil {
		span.SetStatus(codes.Error, "failed querying masterlist")
		span.RecordError(err)
		return nil, err
	}
	defer rows.Close()

	var mLst []domain.Masterlist
	for rows.Next() {
		var lst domain.Masterlist
		if err := rows.Scan(
			&lst.PsgcCode,
			&lst.Name,
			&lst.Code,
			&lst.Level,
		); err != nil {
			return nil, err
		}
		mLst = append(mLst, lst)
	}
	return mLst, nil
}

func (p *dbMasterlistRepository) paginatedQuery(
	ctx context.Context,
	level string,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	queryParams := []interface{}{}
	query := `SELECT * FROM masterlist`
	countQuery := `SELECT COUNT(*) FROM masterlist`

	if level != "" {
		query += fmt.Sprintf(" WHERE level = '%s'", level)
		countQuery += fmt.Sprintf(" WHERE level = '%s'", level)
	}

	if params.Filter != "" {
		query += `
            AND (
                LOWER(psgc_code) LIKE '%' || LOWER($1) || '%' OR
                LOWER(name) LIKE '%' || LOWER($1) || '%' OR
                LOWER(code) LIKE '%' || LOWER($1) || '%'
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
		return domain.PaginatedMasterlist{}, err
	}

	totalItems := 0
	p.conn.QueryRowContext(ctx, countQuery).
		Scan(&totalItems)

	totalPages := (totalItems + params.PerPage - 1) / params.PerPage

	if len(lst) == 0 {
		lst = []domain.Masterlist{}
	}

	metaData := domain.MetaData{
		Page:       params.Page,
		TotalPages: totalPages,
		PerPage:    params.PerPage,
		TotalItems: totalItems,
		ItemCount:  len(lst),
	}

	res := domain.PaginatedMasterlist{
		MetaData: metaData,
		Data:     lst,
	}

	return res, nil
}

func (p *dbMasterlistRepository) GetAll(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	res, err := p.paginatedQuery(ctx, "Bgy", params)

	return res, err
}

func (p *dbMasterlistRepository) Create(
	ctx context.Context,
	data *domain.Masterlist,
) error {
	query := `
		INSERT OR REPLACE INTO masterlist (psgc_code, name, code, level)
		VALUES (?, ?, ?, ?);`

	ctx, span := spanWithQuery(ctx, p.tracer, query)
	defer span.End()

	_, err := p.conn.ExecContext(
		ctx,
		query,
		data.PsgcCode,
		data.Name,
		data.Code,
		data.Level,
	)
	if err != nil {
		span.SetStatus(codes.Error, "failed inserting masterlist")
		span.RecordError(err)
		return err
	}

	return nil
}

func (p *dbMasterlistRepository) CreateBatch(
    ctx context.Context,
    data []*domain.Masterlist,
) error {
    query := `
        INSERT OR REPLACE INTO masterlist (psgc_code, name, code, level)
        VALUES (?, ?, ?, ?);`

    ctx, span := spanWithQuery(ctx, p.tracer, query)
    defer span.End()

    // Start a transaction
    tx, err := p.conn.BeginTx(ctx, nil)
    if err != nil {
        span.SetStatus(codes.Error, "failed to start transaction")
        span.RecordError(err)
        return err
    }
    defer tx.Rollback() // Ensure the transaction is rolled back in case of an error

    // Prepare the insert statement
    stmt, err := tx.PrepareContext(ctx, query)
    if err != nil {
        span.SetStatus(codes.Error, "failed to prepare statement")
        span.RecordError(err)
        return err
    }
    defer stmt.Close()

    for _, d := range data {
        _, err := stmt.ExecContext(ctx, d.PsgcCode, d.Name, d.Code, d.Level)
        if err != nil {
            span.SetStatus(codes.Error, "failed inserting masterlist")
            span.RecordError(err)
            return err
        }
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        span.SetStatus(codes.Error, "failed to commit transaction")
        span.RecordError(err)
        return err
    }

    return nil
}

func (p *dbMasterlistRepository) GetBarangayList(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	res, err := p.paginatedQuery(ctx, "Bgy", params)

	return res, err
}

func (p *dbMasterlistRepository) GetBarangayById(
	ctx context.Context,
	psgcCode string,
) (domain.Masterlist, error) {
	query := `SELECT * FROM masterlist WHERE level = 'Bgy' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Masterlist{}, err
	}

	if len(accs) == 0 {
		return domain.Masterlist{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbMasterlistRepository) GetCityList(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	res, err := p.paginatedQuery(ctx, "City", params)

	return res, err
}

func (p *dbMasterlistRepository) GetCityById(
	ctx context.Context,
	psgcCode string,
) (domain.Masterlist, error) {
	query := `SELECT * FROM masterlist WHERE level = 'City' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Masterlist{}, err
	}

	if len(accs) == 0 {
		return domain.Masterlist{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbMasterlistRepository) GetProvinceList(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	res, err := p.paginatedQuery(ctx, "Prov", params)

	return res, err
}

func (p *dbMasterlistRepository) GetProvinceById(
	ctx context.Context,
	psgcCode string,
) (domain.Masterlist, error) {
	query := `SELECT * FROM masterlist WHERE level = 'Prov' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Masterlist{}, err
	}

	if len(accs) == 0 {
		return domain.Masterlist{}, domain.ErrNotFound
	}
	return accs[0], nil
}

func (p *dbMasterlistRepository) GetRegionList(
	ctx context.Context,
	params domain.PaginationParams,
) (domain.PaginatedMasterlist, error) {
	res, err := p.paginatedQuery(ctx, "Reg", params)

	return res, err
}

func (p *dbMasterlistRepository) GetRegionById(
	ctx context.Context,
	psgcCode string,
) (domain.Masterlist, error) {
	query := `SELECT * FROM masterlist WHERE level = 'Reg' AND psgc_code = $1`

	accs, err := p.fetch(ctx, query, psgcCode)
	if err != nil {
		return domain.Masterlist{}, err
	}

	if len(accs) == 0 {
		return domain.Masterlist{}, domain.ErrNotFound
	}
	return accs[0], nil
}
