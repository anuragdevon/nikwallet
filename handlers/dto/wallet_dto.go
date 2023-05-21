package dto

import "nikwallet/repository/money"

type MoneyTransferDTO struct {
	Amount         *money.Money `json:"amount"`
	RecipientEmail string       `json:"recipient_email"`
}

type Response struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
