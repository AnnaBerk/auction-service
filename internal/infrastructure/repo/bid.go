package repo

import (
	"auction/internal/domain"
	"context"
	"github.com/go-pg/pg/v10"
)

type BidRepository interface {
	GetBidsByAuctionID(ctx context.Context, auctionID int) ([]domain.Bid, error)
	GetWinningBid(ctx context.Context, auctionID, winnerID int) (domain.Bid, error)
	GetUserBid(ctx context.Context, auctionID, userID int) (domain.Bid, error)
}

type bidRepo struct {
	db *pg.DB
}

func NewBidRepository(db *pg.DB) BidRepository {
	return &bidRepo{db: db}
}

func (r *bidRepo) GetBidsByAuctionID(ctx context.Context, auctionID int) ([]domain.Bid, error) {
	var bids []domain.Bid
	err := r.db.Model(&bids).Where("auction_id = ?", auctionID).Select()
	return bids, err
}

func (r *bidRepo) GetWinningBid(ctx context.Context, auctionID, winnerID int) (domain.Bid, error) {
	var bid domain.Bid
	err := r.db.Model(&bid).Where("auction_id = ? AND user_id = ?", auctionID, winnerID).Select()
	return bid, err
}

func (r *bidRepo) GetUserBid(ctx context.Context, auctionID, userID int) (domain.Bid, error) {
	var bid domain.Bid
	err := r.db.Model(&bid).Where("auction_id = ? AND user_id = ?", auctionID, userID).Select()
	return bid, err
}
