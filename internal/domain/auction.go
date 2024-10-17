package domain

import (
	"time"
)

type Lot struct {
	LotID      int
	Title      string
	StartPrice int
	Step       int
	UserID     int
	CreatedAt  time.Time
	AuctionID  int
	ClosedAt   *time.Time
}

type Bid struct {
	BidID     int
	Price     int64
	CreatedAt time.Time
	UserID    int
	LotID     int
	AuctionID int
}

type Auction struct {
	AuctionID int
	CreatedAt time.Time
	ClosedAt  *time.Time
	UserID    *int
	WinnerID  *int
	User      *User
	Winner    *User
}

type User struct {
	UserID  int
	Name    string
	Email   string
	Balance *int64
}

// Проверка валидности ставки и баланса пользователя
func ValidateBid(bid Bid, userBalance int64, currentBids []Bid) error {
	if bid.Price <= 0 {
		return ErrInvalidBidAmount
	}

	totalCommittedAmount := bid.Price
	for _, b := range currentBids {
		totalCommittedAmount += b.Price
	}

	if totalCommittedAmount > userBalance {
		return ErrInsufficientFunds
	}

	return nil
}

// ValidateLot проверяет, что данные лота корректны
func ValidateLot(lot Lot) error {
	if lot.StartPrice <= 0 || lot.Step <= 0 {
		return ErrInvalidLotData
	}
	return nil
}
