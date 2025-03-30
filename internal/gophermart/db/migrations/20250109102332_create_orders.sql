-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
			user_login varchar(500) REFERENCES users (login),
			id bigint PRIMARY KEY UNIQUE,
			status varchar(200),
			accrual double precision,
			uploaded_at timestamptz DEFAULT NOW()
		);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
