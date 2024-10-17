package repo

import "auction/internal/domain"

func NewDomainAuction(auction *Auction) domain.Auction {
	return domain.Auction{
		AuctionID: auction.ID,
		CreatedAt: auction.CreatedAt,
		ClosedAt:  auction.ClosedAt,
		UserID:    auction.UserID,
		WinnerID:  auction.WinnerID,
		User:      NewDomainUser(auction.User),
		Winner:    NewDomainUser(auction.Winner),
	}
}

func NewDomainAuctions(dbAuctions []*Auction) []domain.Auction {
	auctions := make([]domain.Auction, len(dbAuctions))
	for i, dbAuction := range dbAuctions {
		auctions[i] = NewDomainAuction(dbAuction)
	}
	return auctions
}

func NewDomainBid(bid *Bid) domain.Bid {
	return domain.Bid{
		BidID:     bid.ID,
		Price:     bid.Price,
		CreatedAt: bid.CreatedAt,
		UserID:    bid.UserID,
		LotID:     bid.LotID,
		AuctionID: bid.AuctionID,
	}
}

func NewDomainUser(user *User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		UserID:  user.ID,
		Name:    user.Name,
		Email:   user.Email,
		Balance: user.Balance,
	}
}

func NewDatabaseBid(bid domain.Bid) *Bid {
	return &Bid{
		ID:        bid.BidID,
		Price:     bid.Price,
		CreatedAt: bid.CreatedAt,
		UserID:    bid.UserID,
		LotID:     bid.LotID,
		AuctionID: bid.AuctionID,
	}
}

func NewDatabaseUser(user domain.User) *User {
	return &User{
		ID:      user.UserID,
		Name:    user.Name,
		Email:   user.Email,
		Balance: user.Balance,
	}
}

func NewDatabaseLot(lot domain.Lot) *Lot {
	return &Lot{
		ID:         lot.LotID,
		Title:      lot.Title,
		StartPrice: int64(lot.StartPrice),
		Step:       int64(lot.Step),
		CreatedAt:  lot.CreatedAt,
		UserID:     lot.UserID,
		AuctionID:  lot.AuctionID,
	}
}

func NewDatabaseAuction(auction domain.Auction) *Auction {
	return &Auction{
		ID:        auction.AuctionID,
		CreatedAt: auction.CreatedAt,
		ClosedAt:  auction.ClosedAt,
		WinnerID:  auction.WinnerID,
	}
}

func NewDatabaseAuctions(auctions []domain.Auction) []*Auction {
	dbAuctions := make([]*Auction, len(auctions))
	for i, auction := range auctions {
		dbAuctions[i] = NewDatabaseAuction(auction)
	}
	return dbAuctions
}

func NewDomainLot(dbLot Lot) domain.Lot {
	return domain.Lot{
		LotID:      dbLot.ID,
		Title:      dbLot.Title,
		StartPrice: int(dbLot.StartPrice),
		Step:       int(dbLot.Step),
		CreatedAt:  dbLot.CreatedAt,
		UserID:     dbLot.UserID,
		AuctionID:  dbLot.AuctionID,
	}
}
