package repo

import (
	"auction/internal/domain"
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
)

type LotRepository interface {
	Create(ctx context.Context, lot domain.Lot) (int, error)
	PlaceBid(ctx context.Context, bid domain.Bid) (int, error)
	GetUserBids(ctx context.Context, userID int) ([]domain.Bid, error)
	GetLotByID(ctx context.Context, id int) (domain.Lot, error)
}

type LotRepo struct {
	db *pg.DB
}

func NewLotRepository(db *pg.DB) *LotRepo {
	return &LotRepo{db: db}
}

func (r *LotRepo) Create(ctx context.Context, lot domain.Lot) (int, error) {
	dbLot := NewDatabaseLot(lot)
	_, err := r.db.Model(dbLot).Insert()
	if err != nil {
		return 0, err
	}
	return dbLot.ID, nil
}

func (r *LotRepo) PlaceBid(ctx context.Context, bid domain.Bid) (int, error) {
	dbBid := NewDatabaseBid(bid)
	_, err := r.db.Model(dbBid).Insert()
	if err != nil {
		return 0, err
	}
	return dbBid.ID, nil
}

func (r *LotRepo) GetUserBids(ctx context.Context, userID int) ([]domain.Bid, error) {
	var dbBids []Bid
	err := r.db.Model(&dbBids).Where("user_id = ?", userID).Select()
	if err != nil {
		return nil, err
	}

	var domainBids []domain.Bid
	for _, dbBid := range dbBids {
		domainBids = append(domainBids, NewDomainBid(&dbBid))
	}
	return domainBids, nil
}

func (r *LotRepo) GetLotByID(ctx context.Context, id int) (domain.Lot, error) {
	var dbLot Lot
	err := r.db.Model(&dbLot).Where("id = ?", id).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return domain.Lot{}, domain.ErrLotNotFound
		}
		return domain.Lot{}, err
	}
	return NewDomainLot(dbLot), nil
}
