package domain

import "context"

// AuctionService - интерфейс для всех операций аукциона
type AuctionService interface {
	CreateLot(ctx context.Context, lot Lot) (int, error)
	RefillBalance(ctx context.Context, userID int, amount int64) error
	PlaceBid(ctx context.Context, bid Bid) (int, error)
	GetCompletedAuctionsWithoutWinner(ctx context.Context) ([]Auction, error)
	GetBidsByAuctionID(ctx context.Context, auctionID int) ([]Bid, error)
	ProcessTransactions(ctx context.Context, auctionID, winnerID int, losers []int) error
	NotifyAuctionResults(ctx context.Context, auctionID, winnerID int, losers []int) error
	DetermineWinner(ctx context.Context, bids []Bid) (int, []int, error)
	GetNewAuctions(ctx context.Context) ([]Auction, error)
	NotifyUsersAboutNewAuctions(ctx context.Context) error
}
