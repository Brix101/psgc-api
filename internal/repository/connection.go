package repository

import (
	"context"
	"database/sql"

	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

// Connection is an interface that defines common query operations.
type Connection interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func spanWithQuery(
	ctx context.Context,
	tracer trace.Tracer,
	query string,
) (context.Context, trace.Span) {
	ctx, span := tracer.Start(ctx, "db:query")
	span.SetAttributes(semconv.DBStatementKey.String(query))
	return ctx, span
}
