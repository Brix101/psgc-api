-- +goose Up
-- +goose StatementBegin
CREATE TABLE masterlist (
	psgc_code TEXT PRIMARY KEY,
	name TEXT,
	code TEXT,
	level TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE masterlist
-- +goose StatementEnd
