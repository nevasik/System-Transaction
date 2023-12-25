package postgresql

import (
	"database/sql"
	"fmt"
	"golangProject/internal/entity"
	"log"
	"os"
	"sync"
	"time"
)

const (
	StatusCreated = "Created"
	StatusSuccess = "Success"
	StatusError   = "Error"
)

type Database struct {
	Mu *sync.Mutex
	Db *sql.DB
}

func NewDatabase(db *sql.DB) *Database {
	return &Database{
		Mu: &sync.Mutex{},
		Db: db,
	}
}

func ConnectionString() string {
	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbname := os.Getenv("db_name")

	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
}

func (db *Database) CreateTransaction(currency string, amount float64, walletOrCardNumber string) (*entity.Transaction, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	tx, err := db.Db.Begin()
	if err != nil {
		return nil, err
	}

	status := StatusCreated
	createdAt := time.Now()

	defer func() {
		if err != nil {
			status = StatusError
			log.Println("An error occurred, transaction status is ERROR")
			tx.Rollback()
		} else {
			status = StatusSuccess
			if err := tx.Commit(); err != nil {
				log.Println("Error committing transaction:", err)
			}
		}
	}()

	return &entity.Transaction{
		Currency:     currency,
		Amount:       amount,
		WalletOrCard: walletOrCardNumber,
		Status:       status,
		CreatedAt:    &createdAt,
	}, nil
}

func (db *Database) GetBalance(walletOrCardNumber string) (float64, error) {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	row := db.Db.QueryRow(`
		SELECT balance
		FROM balance
		WHERE wallet_or_card_number = $1
	`, walletOrCardNumber)

	var balance float64
	if err := row.Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}

func (db *Database) WithdrawFunds(walletOrCardNumber string, amount float64) error {
	db.Mu.Lock()
	defer db.Mu.Unlock()

	result, err := db.Db.Exec(`
		UPDATE balance
		SET balance = balance - $1
		WHERE wallet_or_card_number = $2
		AND balance >= $1
	`, amount, walletOrCardNumber)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("insufficient funds for withdrawal")
	}

	return nil
}
