package app

import (
	"auction/internal/domain"
	"auction/internal/infrastructure/notify"
	"auction/internal/infrastructure/payment"
	"auction/internal/infrastructure/repo"
	"context"
	"errors"
	"fmt"
	"time"
)

type AuctionService struct {
	lotRepo     repo.LotRepository
	userRepo    repo.UserRepository
	auctionRepo repo.AuctionRepository
	bidRepo     repo.BidRepository
	notify      notify.NotifyService
	balance     payment.BalanceService
}

func NewAuctionService(lotRepo repo.LotRepository,
	userRepo repo.UserRepository,
	auctionRepo repo.AuctionRepository,
	bidRepo repo.BidRepository,
	notify notify.NotifyService,
	balance payment.BalanceService) *AuctionService {
	return &AuctionService{
		lotRepo:     lotRepo,
		userRepo:    userRepo,
		auctionRepo: auctionRepo,
		bidRepo:     bidRepo,
		notify:      notify,
		balance:     balance,
	}
}

func (s *AuctionService) CreateLot(ctx context.Context, lot domain.Lot) (int, error) {
	auction := domain.Auction{
		CreatedAt: time.Now(),
	}
	err := domain.ValidateLot(lot)
	if err != nil {
		return 0, err
	}
	auction.ClosedAt = lot.ClosedAt
	auctionID, err := s.auctionRepo.Create(ctx, auction)
	if err != nil {
		return 0, err
	}

	lot.AuctionID = auctionID

	return s.lotRepo.Create(ctx, lot)
}

func (s *AuctionService) RefillBalance(ctx context.Context, userID int, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than zero")
	}
	return s.userRepo.RefillBalance(ctx, userID, amount)
}

func (s *AuctionService) PlaceBid(ctx context.Context, bid domain.Bid) (int, error) {
	lot, err := s.lotRepo.GetLotByID(ctx, bid.LotID)
	if err != nil {
		if errors.Is(err, domain.ErrLotNotFound) {
			return 0, domain.ErrLotNotFound
		}
		return 0, err
	}

	bid.AuctionID = lot.AuctionID
	userBalance, err := s.userRepo.GetBalance(ctx, bid.UserID)
	if err != nil {
		return 0, err
	}

	currentBids, err := s.lotRepo.GetUserBids(ctx, bid.UserID)
	if err != nil {
		return 0, err
	}

	if err := domain.ValidateBid(bid, *userBalance, currentBids); err != nil {
		return 0, err
	}

	return s.lotRepo.PlaceBid(ctx, bid)
}

func (s *AuctionService) GetCompletedAuctionsWithoutWinner(ctx context.Context) ([]domain.Auction, error) {
	return s.auctionRepo.GetCompletedAuctionsWithoutWinner(ctx)
}

func (s *AuctionService) GetBidsByAuctionID(ctx context.Context, auctionID int) ([]domain.Bid, error) {
	return s.bidRepo.GetBidsByAuctionID(ctx, auctionID)
}

func (s *AuctionService) DetermineWinner(ctx context.Context, bids []domain.Bid) (int, []int, error) {
	if len(bids) == 0 {
		return 0, nil, errors.New("no bids available")
	}

	var winnerID int
	var maxBid int64
	losers := []int{}

	for _, bid := range bids {
		if bid.Price > maxBid {
			maxBid = bid.Price
			winnerID = bid.UserID
		}
	}

	for _, bid := range bids {
		if bid.UserID != winnerID {
			losers = append(losers, bid.UserID)
		}
	}

	return winnerID, losers, nil
}

func (s *AuctionService) ProcessTransactions(ctx context.Context, auctionID, winnerID int, losers []int) error {
	winningBid, err := s.bidRepo.GetWinningBid(ctx, auctionID, winnerID)
	if err != nil {
		return err
	}

	if err := s.balance.DeductBalance(ctx, winnerID, winningBid.Price); err != nil {
		return err
	}

	for _, loserID := range losers {
		losingBid, err := s.bidRepo.GetUserBid(ctx, auctionID, loserID)
		if err != nil {
			return err
		}

		if err := s.balance.RefundBalance(ctx, loserID, losingBid.Price); err != nil {
			return err
		}
	}

	return s.auctionRepo.CloseAuction(ctx, auctionID, winnerID)
}

func (s *AuctionService) NotifyAuctionResults(ctx context.Context, auctionID, winnerID int, losers []int) error {
	err := s.notify.NotifyUser(ctx, winnerID, fmt.Sprintf("Вы победили в аукционе %d", auctionID))
	if err != nil {
		return err
	}

	for _, loserID := range losers {
		err := s.notify.NotifyUser(ctx, loserID, fmt.Sprintf("Вы проиграли в аукционе %d", auctionID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *AuctionService) GetNewAuctions(ctx context.Context) ([]domain.Auction, error) {
	return s.auctionRepo.GetNewAuctions(ctx)
}

func (s *AuctionService) NotifyUsersAboutNewAuctions(ctx context.Context) error {
	newAuctions, err := s.auctionRepo.GetNewAuctions(ctx)
	if err != nil {
		return err
	}

	if len(newAuctions) == 0 {
		return nil
	}

	return s.notify.NotifyAllUsersAboutNewAuctions(ctx, newAuctions)
}
