package db

import (
	"database/sql"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (ps *PostgresStorage) AddUser(user models.User) error {
	queryAddUser := ps.pSQL.Insert(
		"users",
	).Columns(
		"login",
		"password",
	).Values(
		user.Login,
		user.Password,
	).Suffix("ON CONFLICT (login) DO NOTHING")

	res, err := queryAddUser.Exec()
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil {
		return err
	} else if affected == 0 {
		return &gophermart.UsedLoginError{}
	}
	return nil
}

func (ps *PostgresStorage) UserByLogin(login string) (*models.User, error) {
	var (
		userLogin     string
		userPassword  string
		userBalance   sql.NullFloat64
		userWithdrawn sql.NullFloat64
	)
	query := ps.pSQL.Select(
		"login", "password", "balance", "withdrawn",
	).From(
		"users",
	).Where(sq.Eq{"login": login})

	row := query.QueryRow()
	if err := row.Scan(&userLogin, &userPassword, &userBalance, &userWithdrawn); err != nil {
		return nil, err
	}
	return &models.User{
		Login:     userLogin,
		Password:  userPassword,
		Balance:   userBalance.Float64,
		Withdrawn: userWithdrawn.Float64,
	}, nil
}
