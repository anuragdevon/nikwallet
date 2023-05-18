package database

import (
	"fmt"
	"time"

	"nikwallet/database/money"
)

type Wallet struct {
	ID        int
	UserID    int
	Money     money.Money
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (db *PostgreSQL) CreateWallet(userID int) (*Wallet, error) {
	_, err := db.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	initialZeroMoney, _ := money.NewMoney(0, "INR")
	wallet := &Wallet{
		UserID:    userID,
		Money:     *initialZeroMoney,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `INSERT INTO wallet(user_id, amount, currency, created_at, updated_at) 
			  VALUES($1, $2, $3, $4, $5) 
			  RETURNING id, created_at, updated_at`
	err = db.DB.QueryRow(query, wallet.UserID, wallet.Money.Amount, wallet.Money.Currency, wallet.CreatedAt, wallet.UpdatedAt).
		Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}

func (db *PostgreSQL) GetWalletByUserID(userID int) (*Wallet, error) {
	query := `SELECT id, user_id, amount, currency, created_at, updated_at FROM wallet WHERE user_id = $1`
	row := db.DB.QueryRow(query, userID)

	wallet := &Wallet{}
	err := row.Scan(&wallet.ID, &wallet.UserID, &wallet.Money.Amount, &wallet.Money.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("no wallets found for user with ID %d", userID)
	}

	return wallet, nil
}

func (db *PostgreSQL) AddMoneyToWallet(userID int, moneyToAdd money.Money) error {
	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return err
	}

	newMoney, err := wallet.Money.Add(&moneyToAdd)
	if err != nil {
		return err
	}

	query := `UPDATE wallet SET amount=$1, updated_at=$2 WHERE id=$3`
	_, err = db.DB.Exec(query, newMoney.Amount, time.Now(), wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to add money to wallet: %w", err)
	}

	return nil
}

func (db *PostgreSQL) WithdrawMoneyFromWallet(userID int, moneyToWithdraw money.Money) (money.Money, error) {
	wallet, err := db.GetWalletByUserID(userID)
	if err != nil {
		return money.Money{}, err
	}

	remainedMoney, err := wallet.Money.Subtract(&moneyToWithdraw)
	if err != nil {
		return money.Money{}, err
	}

	query := `UPDATE wallet SET amount=$1, updated_at=$2 WHERE id=$3`
	_, err = db.DB.Exec(query, remainedMoney.Amount, time.Now(), wallet.ID)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to withdraw money from wallet: %w", err)
	}

	return moneyToWithdraw, nil
}
