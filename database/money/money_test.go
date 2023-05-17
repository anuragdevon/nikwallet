package money

import (
	"testing"
)

func TestNewMoney(t *testing.T) {
	t.Run("NewMoney create new money for valid input", func(t *testing.T) {
		amount := int64(100)
		currency := "INR"
		m, err := NewMoney(amount, currency)

		if err != nil {
			t.Errorf("NewMoney() error = %v, want nil", err)
			return
		}

		if m.Amount != amount {
			t.Errorf("NewMoney() m.Amount = %d, want %d", m.Amount, amount)
		}

		if m.Currency != currency {
			t.Errorf("NewMoney() m.Currency = %s, want %s", m.Currency, currency)
		}
	})

	t.Run("NewMoney to return error for negetive money input", func(t *testing.T) {
		amount := int64(-100)
		currency := "INR"

		_, err := NewMoney(amount, currency)

		if err == nil {
			t.Errorf("NewMoney() error = nil, want non-nil")
			return
		}
	})

	t.Run("NewMoney to return error for empty currency", func(t *testing.T) {
		amount := int64(100)
		currency := ""

		_, err := NewMoney(amount, currency)

		if err == nil {
			t.Errorf("NewMoney() error = nil, want non-nil")
			return
		}
	})

	t.Run("Add method to add money for valid inputs", func(t *testing.T) {
		hundredRupees, _ := NewMoney(100, "INR")
		hundredFiftyRupees, err := hundredRupees.Add(&Money{Amount: 50, Currency: "INR"})

		if err != nil {
			t.Fatalf("Money.AddMoney() error = %v", err)
		}

		if hundredFiftyRupees.Amount != 150 {
			t.Errorf("Money.AddMoney() = %d, want %d", hundredFiftyRupees.Amount, 150)
		}
	})

	t.Run("Subtract method to deduct second money from first", func(t *testing.T) {
		hundredRupees, _ := NewMoney(100, "INR")
		fiftyRupees, _ := NewMoney(50, "INR")

		subtractedFiftyRupees, _ := hundredRupees.Subtract(fiftyRupees)

		expectedRupees, _ := NewMoney(50, "INR")
		if subtractedFiftyRupees.Amount != expectedRupees.Amount || subtractedFiftyRupees.Currency != expectedRupees.Currency {
			t.Errorf("Money.Subtract() = %v, want %v", subtractedFiftyRupees, expectedRupees)
		}
	})

	t.Run("Subtract method to return error for negetive subtraction result", func(t *testing.T) {
		hundredRupees, _ := NewMoney(100, "INR")
		fiftyRupees, _ := NewMoney(200, "INR")

		_, err := hundredRupees.Subtract(fiftyRupees)

		if err == nil {
			t.Errorf("Money.Subtract() error = nil, want non-nil")
		}
	})
}
