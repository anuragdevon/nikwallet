package money

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Currency string

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
	INR Currency = "INR"
)

var ConversionFactors = map[Currency]decimal.Decimal{
	USD: decimal.NewFromFloat(0.012),
	EUR: decimal.NewFromFloat(0.011),
	INR: decimal.NewFromFloat(1.0),
}

var ZeroAmountValue = decimal.NewFromFloat(0.0)

type Money struct {
	Amount   decimal.Decimal
	Currency Currency
}

func NewMoney(amount decimal.Decimal, currency Currency) (*Money, error) {
	if amount.LessThan(ZeroAmountValue) {
		return nil, fmt.Errorf("amount cannot be negative")
	}
	_, ok := ConversionFactors[currency]
	if !ok {
		return nil, fmt.Errorf("unsupported currency: %s", currency)
	}

	return &Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

func (m *Money) ToBaseCurrency() (*Money, error) {
	baseFactor, err := ConversionFactors[m.Currency]
	if !err {
		return nil, fmt.Errorf("unsupported currency conversion")
	}

	convertedAmount := m.Amount.Mul(baseFactor)

	return &Money{
		Amount:   convertedAmount,
		Currency: INR,
	}, nil
}

func (m *Money) Add(money *Money) (*Money, error) {
	baseCurrencyMoney, err := m.ToBaseCurrency()
	if err != nil {
		return nil, err
	}

	otherBaseCurrencyMoney, err := money.ToBaseCurrency()
	if err != nil {
		return nil, err
	}

	addedAmount := baseCurrencyMoney.Amount.Add(otherBaseCurrencyMoney.Amount)

	conversionFactor, ok := ConversionFactors[m.Currency]
	if !ok {
		return nil, fmt.Errorf("unsupported currency conversion")
	}

	convertedAmount := addedAmount.Div(conversionFactor).Round(2)

	return &Money{
		Amount:   convertedAmount,
		Currency: m.Currency,
	}, nil
}

func (m *Money) Subtract(money *Money) (*Money, error) {
	if m.Currency != money.Currency {
		return nil, fmt.Errorf("cannot subtract money with different currency")
	}

	baseCurrencyMoney, err := m.ToBaseCurrency()
	if err != nil {
		return nil, err
	}

	otherBaseCurrencyMoney, err := money.ToBaseCurrency()
	if err != nil {
		return nil, err
	}

	if baseCurrencyMoney.Amount.LessThan(otherBaseCurrencyMoney.Amount) {
		return nil, fmt.Errorf("not enough money to deduct")
	}

	subtractedAmount := baseCurrencyMoney.Amount.Sub(otherBaseCurrencyMoney.Amount)

	conversionFactor, ok := ConversionFactors[m.Currency]
	if !ok {
		return nil, fmt.Errorf("unsupported currency conversion")
	}

	convertedAmount := subtractedAmount.Div(conversionFactor).Round(2)

	return &Money{
		Amount:   convertedAmount,
		Currency: m.Currency,
	}, nil
}
