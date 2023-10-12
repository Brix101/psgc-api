-- +goose Up
-- +goose StatementBegin
CREATE TABLE city_muni (
	psgc_code TEXT PRIMARY KEY,
	prov_code TEXT,
	name TEXT,
	level TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE city_muni
-- +goose StatementEnd