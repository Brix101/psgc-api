-- +goose Up
-- +goose StatementBegin
CREATE TABLE barangay (
	psgc_code TEXT PRIMARY KEY,
	citmun_code TEXT,
	name TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE barangay
-- +goose StatementEnd