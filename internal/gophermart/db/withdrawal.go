package db

import (
	"context"
	"musthave_tpl/internal/models"

	sq "github.com/Masterminds/squirrel"
)

func (ps *PostgresStorage) AddTransaction(transaction models.Transaction) (*models.Transaction, error) {
	ctx := context.Background()
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	queryWithdraw := ps.pSQL.Insert(
		"withdrawals",
	).Columns(
		"user_login",
		"order_id",
		"amount",
		"processed_at",
	).Values(transaction.User, transaction.ID, transaction.Amount, transaction.Date)

	_, err = queryWithdraw.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	queryUser := ps.pSQL.Update(
		"users",
	).Set(
		"balance",
		sq.Expr("balance - ?", transaction.Amount),
	).Set(
		"withdrawn",
		sq.Expr("withdrawn + ?", transaction.Amount),
	).Where(sq.Eq{"login": transaction.User})
	_, err = queryUser.ExecContext(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	return nil, err
}

func (ps *PostgresStorage) TransactionsByUser(login string) ([]models.Transaction, error) {
	query := ps.pSQL.Select(
		"user_login",
		"order_id",
		"amount",
		"processed_at",
	).From(
		"withdrawals",
	).Where(sq.Eq{"user_login": login})

	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	if rows.Err() != nil {
		return nil, err
	}
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.User, &transaction.ID, &transaction.Amount, &transaction.Date)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil
}
