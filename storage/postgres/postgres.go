package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"wallet/internal/config"
	"wallet/storage"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(cfg *config.DBConfig) *Storage {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname)

	ctx := context.Background()
	deadline := time.Now().Add(15 * time.Second)

	for {
		db, err := pgxpool.New(ctx, psqlInfo)
		if err == nil {
			err = db.Ping(ctx)
		}

		if err == nil {
			return &Storage{db: db}
		}

		if time.Now().After(deadline) {
			fmt.Fprintf(os.Stderr, "Unable to connect to database after 15 seconds: %v\n", err)
			os.Exit(1)
		}

		time.Sleep(1 * time.Second)
	}
}

func (s *Storage) ChangeAmount(id string, amount int, operation string) error {

	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.Background())
		} else {
			err = tx.Commit(context.Background())
		}
	}()

	var currentAmount int
	err = tx.QueryRow(context.Background(), `SELECT amount FROM wallets WHERE id = $1 FOR UPDATE`, id).Scan(&currentAmount)
	if err != nil {
		if err == pgx.ErrNoRows {
			return storage.ErrWalletNotFound
		}
		return err
	}

	if operation == "DEPOSIT" {
		newAmount := currentAmount + amount
		if newAmount < 0 {
			return storage.ErrNegativeAmount
		}

		_, err = tx.Exec(context.Background(), `UPDATE wallets SET amount = $1 WHERE id = $2`, newAmount, id)
		if err != nil {
			return err
		}

	} else if operation == "WITHDRAW" {
		newAmount := currentAmount - amount
		if newAmount < 0 {
			return storage.ErrNegativeAmount
		}

		_, err = tx.Exec(context.Background(), `UPDATE wallets SET amount = $1 WHERE id = $2`, newAmount, id)
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Storage) Amount(id string) (int, error) {
	var amount int

	err := s.db.QueryRow(context.Background(), `SELECT amount FROM wallets WHERE id = $1`, id).Scan(&amount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, storage.ErrWalletNotFound
		}
		return 0, err
	}

	return amount, nil
}
