-- +goose Up
-- +goose StatementBegin
CREATE TABLE city (
	psgc_code TEXT PRIMARY KEY,
	prov_code TEXT,
	name TEXT,
	code TEXT,
	level TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE city
-- +goose StatementEnd