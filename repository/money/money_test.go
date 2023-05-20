package money

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestMoney(t *testing.T) {
	t.Run("NewMoney create new money for valid input INR", func(t *testing.T) {
		amount := decimal.NewFromFloat(10.0)
		currency := INR
		newTenRupeesMoney, err := NewMoney(amount, currency)

		if err != nil {
			t.Errorf("NewMoney() error = %v, want nil", err)
			return
		}

		if !newTenRupeesMoney.Amount.Equal(amount) {
			t.Errorf("NewMoney() Amount = %d, want %d", newTenRupeesMoney.Amount, amount)
		}

		if newTenRupeesMoney.Currency != currency {
			t.Errorf("NewMoney() Currency = %s, want %s", newTenRupeesMoney.Currency, currency)
		}
	})

	t.Run("NewMoney to return error for negetive money input", func(t *testing.T) {
		amount := decimal.NewFromFloat(-100.0)
		currency := INR

		_, err := NewMoney(amount, currency)

		if err == nil {
			t.Errorf("NewMoney() error = nil, want non-nil")
			return
		}
	})

	t.Run("NewMoney to return error for invalid currency", func(t *testing.T) {
		_, err := NewMoney(decimal.NewFromFloat(100.0), Currency("DIR"))

		if err == nil {
			t.Errorf("NewMoney() error = nil, want non-nil")
			return
		}
	})

	t.Run("ToBaseCurrency to convert USD to INR", func(t *testing.T) {
		hundredDollars, _ := NewMoney(decimal.NewFromFloat(100.0), USD)
		conversionToRupees, _ := NewMoney(decimal.NewFromFloat(1.2), INR)
		result1, err := hundredDollars.ToBaseCurrency()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result1.Amount.Equal(conversionToRupees.Amount) {
			t.Errorf("Expected amount: %s, but got: %s", conversionToRupees.Amount.String(), result1.Amount.String())
		}
		if result1.Currency != INR {
			t.Errorf("Expected currency: %s, but got: %s", INR, result1.Currency)
		}
	})
	t.Run("ToBaseCurrency to convert USD to INR", func(t *testing.T) {
		hundredEuros, _ := NewMoney(decimal.NewFromFloat(100.0), EUR)
		conversionToRupees, _ := NewMoney(decimal.NewFromFloat(1.1), INR)
		result2, err := hundredEuros.ToBaseCurrency()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result2.Amount.Equal(conversionToRupees.Amount) {
			t.Errorf("Expected amount: %s, but got: %s", conversionToRupees.Amount.String(), result2.Amount.String())
		}
		if result2.Currency != INR {
			t.Errorf("Expected currency: %s, but got: %s", INR, result2.Currency)
		}
	})

	t.Run("ToBaseCurrency to convert USD to INR", func(t *testing.T) {

		hundredRupees, _ := NewMoney(decimal.NewFromFloat(100.0), INR)
		conversionToRupees, _ := NewMoney(decimal.NewFromFloat(100.0), INR)
		result3, err := hundredRupees.ToBaseCurrency()
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if !result3.Amount.Equal(conversionToRupees.Amount) {
			t.Errorf("Expected amount: %s, but got: %s", conversionToRupees.Amount.String(), result3.Amount.String())
		}
		if result3.Currency != INR {
			t.Errorf("Expected currency: %s, but got: %s", INR, result3.Currency)
		}
	})

	t.Run("Add method to add money for valid inputs", func(t *testing.T) {
		hundredRupees, _ := NewMoney(decimal.NewFromFloat(100.0), INR)
		fiftyRupees, _ := NewMoney(decimal.NewFromFloat(50.0), INR)

		hundredFiftyRupees, err := hundredRupees.Add(fiftyRupees)

		if err != nil {
			t.Fatalf("Money.AddMoney() error = %v", err)
		}

		if !hundredFiftyRupees.Amount.Equal(decimal.NewFromFloat(150.0)) {
			t.Errorf("Money.AddMoney()got = %v, want %v", hundredFiftyRupees.Amount, 150.0)
		}
	})

	t.Run("Add method to add two different currency, USD and EUR and return USD", func(t *testing.T) {
		hundredDollars, _ := NewMoney(decimal.NewFromFloat(100.0), USD)
		fiftyEuros, _ := NewMoney(decimal.NewFromFloat(50.0), EUR)

		expectedAmount, _ := NewMoney(decimal.NewFromFloat(145.83), USD)

		result, err := hundredDollars.Add(fiftyEuros)
		if err != nil {
			t.Fatalf("Money.Add() error: %v", err)
		}

		if !result.Amount.Equal(expectedAmount.Amount) {
			t.Errorf("Money.Add() got = %s, want = %s", result.Amount.String(), expectedAmount.Amount.String())
		}

		if result.Currency != USD {
			t.Errorf("Money.Add() currency got = %s, want = %s", result.Currency, USD)
		}
	})

	t.Run("Subtract method to deduct second money from first", func(t *testing.T) {
		hundredRupees, _ := NewMoney(decimal.NewFromFloat(100.0), INR)
		fiftyRupees, _ := NewMoney(decimal.NewFromFloat(50.0), INR)

		subtractedFiftyRupees, _ := hundredRupees.Subtract(fiftyRupees)

		expectedRupees, _ := NewMoney(decimal.NewFromFloat(50.0), INR)
		if !subtractedFiftyRupees.Amount.Equal(expectedRupees.Amount) || subtractedFiftyRupees.Currency != expectedRupees.Currency {
			t.Errorf("Money.Subtract() got = %v, want %v", subtractedFiftyRupees, expectedRupees)
		}
	})

	t.Run("Subtract method to return error for negetive subtraction result", func(t *testing.T) {
		hundredRupees, _ := NewMoney(decimal.NewFromFloat(100.0), INR)
		fiftyRupees, _ := NewMoney(decimal.NewFromFloat(200.0), INR)

		_, err := hundredRupees.Subtract(fiftyRupees)

		if err == nil {
			t.Errorf("Money.Subtract() error = nil, want non-nil")
		}
	})
}
