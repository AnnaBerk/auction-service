package repo

import (
	"auction/internal/domain"
	"context"
	"github.com/go-pg/pg/v10"
	"time"
)

type AuctionRepository interface {
	Create(ctx context.Context, auction domain.Auction) (int, error)
	GetCompletedAuctionsWithoutWinner(ctx context.Context) ([]domain.Auction, error)
	CloseAuction(ctx context.Context, auctionID, winnerID int) error
	GetNewAuctions(ctx context.Context) ([]domain.Auction, error)
}

type AuctionRepo struct {
	db *pg.DB
}

func NewAuctionRepository(db *pg.DB) *AuctionRepo {
	return &AuctionRepo{db: db}
}

func (r *AuctionRepo) Create(ctx context.Context, auction domain.Auction) (int, error) {
	dbAuction := NewDatabaseAuction(auction)
	_, err := r.db.Model(dbAuction).Insert()
	if err != nil {
		return 0, err
	}
	return dbAuction.ID, nil
}

func (r *AuctionRepo) GetCompletedAuctionsWithoutWinner(ctx context.Context) ([]domain.Auction, error) {
	var dbAuctions []*Auction
	err := r.db.Model(&dbAuctions).Where("closed_at <= ? AND winner_id IS NULL", time.Now()).Select()
	if err != nil {
		return nil, err
	}

	return NewDomainAuctions(dbAuctions), nil
}

func (r *AuctionRepo) CloseAuction(ctx context.Context, auctionID, winnerID int) error {
	_, err := r.db.Model(&Auction{}).Where("id = ?", auctionID).Set("winner_id = ?", winnerID).Set("closed_at = ?", time.Now()).Update()
	return err
}

func (r *AuctionRepo) GetNewAuctions(ctx context.Context) ([]domain.Auction, error) {
	var dbAuctions []*Auction
	err := r.db.Model(&dbAuctions).Where("created_at > ?", time.Now().Add(-24*time.Hour)).Select() //время установлено для упрощения тестирования
	if err != nil {
		return nil, err
	}

	return NewDomainAuctions(dbAuctions), nil
}
