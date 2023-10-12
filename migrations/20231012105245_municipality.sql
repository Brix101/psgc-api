-- +goose Up
-- +goose StatementBegin
CREATE TABLE municipality (
	psgc_code TEXT PRIMARY KEY,
	prov_code TEXT,
	name TEXT,
	code TEXT,
	level TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE municipality
-- +goose StatementEnd