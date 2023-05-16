package wallet

import (
	"fmt"
	"nikwallet/pkg/db"
	"time"

	"nikwallet/pkg/money"
	user "nikwallet/pkg/user"
)

type Wallet struct {
	ID        int
	UserID    int
	Money     money.Money
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateWallet(database *db.DB, userID int) (int, error) {
	_, err := user.GetUserByID(database, userID)
	if err != nil {
		return 0, err
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
	err = database.QueryRow(query, wallet.UserID, wallet.Money.Amount, wallet.Money.Currency, wallet.CreatedAt, wallet.UpdatedAt).
		Scan(&wallet.ID, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet.ID, nil
}

func GetWalletByID(database *db.DB, walletID int) (*Wallet, error) {
	wallet := &Wallet{}
	query := `SELECT id, user_id, amount, currency, created_at, updated_at FROM wallet WHERE id = $1`
	err := database.QueryRow(query, walletID).
		Scan(&wallet.ID, &wallet.UserID, &wallet.Money.Amount, &wallet.Money.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("wallet not found")
	}
	return wallet, nil
}

func AddMoneyToWallet(database *db.DB, walletID int, moneyToAdd money.Money) error {
	wallet, err := GetWalletByID(database, walletID)
	if err != nil {
		return err
	}

	newMoney, err := wallet.Money.Add(&moneyToAdd)
	if err != nil {
		return err
	}

	query := `UPDATE wallet SET amount=$1, updated_at=$2 WHERE id=$3`
	_, err = database.Exec(query, newMoney.Amount, time.Now(), walletID)
	if err != nil {
		return fmt.Errorf("failed to add money to wallet: %w", err)
	}

	return nil
}

func WithdrawMoneyFromWallet(database *db.DB, walletID int, moneyToWithdraw money.Money) (money.Money, error) {
	wallet, err := GetWalletByID(database, walletID)
	if err != nil {
		return money.Money{}, err
	}

	remainedMoney, err := wallet.Money.Subtract(&moneyToWithdraw)
	if err != nil {
		return money.Money{}, err
	}

	query := `UPDATE wallet SET amount=$1, updated_at=$2 WHERE id=$3`
	_, err = database.Exec(query, remainedMoney.Amount, time.Now(), walletID)
	if err != nil {
		return money.Money{}, fmt.Errorf("failed to withdraw money from wallet: %w", err)
	}

	return moneyToWithdraw, nil
}
