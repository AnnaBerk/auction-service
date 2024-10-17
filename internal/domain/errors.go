package domain

import "errors"

var (
	ErrInvalidLotData    = errors.New("invalid lot data: start price and step must be greater than zero")
	ErrInvalidBidAmount  = errors.New("invalid bid amount: amount must be greater than zero")
	ErrLotNotFound       = errors.New("lot not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)
