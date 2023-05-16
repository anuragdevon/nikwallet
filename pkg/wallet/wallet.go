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
