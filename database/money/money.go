package money

import "fmt"

type Money struct {
	Amount   int64
	Currency string
}

func NewMoney(amount int64, currency string) (*Money, error) {
	if amount < 0 {
		return nil, fmt.Errorf("amount cannot be negative")
	}
	if currency == "" {
		return nil, fmt.Errorf("currency cannot be empty")
	}
	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m *Money) Add(money *Money) (*Money, error) {
	if m.Currency != money.Currency {
		return &Money{}, fmt.Errorf("cannot add money with different currency")
	}
	AddedMoney, _ := NewMoney(m.Amount+money.Amount, "INR")
	return AddedMoney, nil
}

func (m *Money) Subtract(money *Money) (*Money, error) {
	if m.Currency != money.Currency {
		return &Money{}, fmt.Errorf("cannot subtract money with different currency")
	}
	if m.Amount < money.Amount {
		return &Money{}, fmt.Errorf("not enough money to deduct")
	}
	SubtractedMoney, _ := NewMoney(m.Amount-money.Amount, "INR")
	return SubtractedMoney, nil
}
