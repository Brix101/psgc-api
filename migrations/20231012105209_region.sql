-- +goose Up
-- +goose StatementBegin
CREATE TABLE region (
	psgc_code TEXT PRIMARY KEY,
	name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE region
-- +goose StatementEnd
