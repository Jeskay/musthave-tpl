-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
			login varchar(500) PRIMARY KEY UNIQUE,
			password text,
			balance double precision DEFAULT 0,
			withdrawn double precision DEFAULT 0
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
