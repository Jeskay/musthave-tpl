package db

import (
	"context"
	"musthave_tpl/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
)

func (ps *PostgresRepository) OrdersByUser(ctx context.Context, login string) ([]models.Order, error) {
	query := ps.pSQL.Select(
		"id", "user_login", "status", "accrual", "uploaded_at",
	).From(
		"orders",
	).Where(
		sq.Eq{"user_login": login},
	).OrderBy("uploaded_at ASC")

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	var orders []models.Order
	for rows.Next() {
		var (
			id         int64
			userLogin  string
			status     string
			accrual    float64
			uploadedAt time.Time
		)
		err := rows.Scan(&id, &userLogin, &status, &accrual, &uploadedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, models.Order{
			User:       models.User{Login: userLogin},
			Number:     id,
			Status:     models.OrderStatus(status),
			Accrual:    accrual,
			UploadedAt: uploadedAt,
		})
	}
	return orders, nil
}

func (ps *PostgresRepository) AddOrder(ctx context.Context, order models.Order) error {

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	orderQuery := ps.pSQL.Insert(
		"orders",
	).Columns(
		"id",
		"user_login",
		"status",
		"accrual",
	).Values(
		order.Number,
		order.User.Login,
		order.Status,
		order.Accrual,
	).Suffix("ON CONFLICT (id) DO NOTHING")

	res, err := orderQuery.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	} else if affected == 0 {
		tx.Rollback()
		login, err := ps.GetLoginByOrderID(ctx, order.Number)
		if err != nil {
			return err
		}
		if login == order.User.Login {
			return &models.OrderExists{}
		}
		return &models.OrderUsed{}
	}

	balanceQuery := ps.pSQL.Update(
		"users",
	).Set(
		"balance",
		sq.Expr("balance + ?", order.Accrual),
	).Where(
		sq.Eq{"login": order.User.Login},
	)

	_, err = balanceQuery.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit()
	return err
}

func (ps *PostgresRepository) GetLoginByOrderID(ctx context.Context, id int64) (string, error) {
	queryUser := ps.pSQL.Select(
		"user_login",
	).From(
		"orders",
	).Where(sq.Eq{"id": id})

	rows, err := queryUser.Query()
	if err != nil {
		return "", err
	}
	if rows.Err() != nil {
		return "", err
	}
	if rows.Next() {
		var login string
		err := rows.Scan(&login)
		if err != nil {
			return "", err
		}
		return login, nil
	}
	return "", nil
}
