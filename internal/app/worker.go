package app

import (
	"auction/internal/domain"
	"context"
	"log"
	"time"
)

type AuctionWorker struct {
	service domain.AuctionService
	logger  *log.Logger
	stopCh  chan struct{}
}

func NewAuctionWorker(service domain.AuctionService, logger *log.Logger) *AuctionWorker {
	return &AuctionWorker{
		service: service,
		logger:  logger,
		stopCh:  make(chan struct{}),
	}
}

func (w *AuctionWorker) Start() {
	w.logger.Println("Auction worker started")
	go w.run()
}

func (w *AuctionWorker) Stop() {
	close(w.stopCh)
	w.logger.Println("Auction worker stopped")
}

func (w *AuctionWorker) run() {
	for {
		select {
		case <-time.After(time.Minute):
			w.processCompletedAuctions()
			w.processNewAuctions()
		case <-w.stopCh:
			return
		}
	}
}

func (w *AuctionWorker) processCompletedAuctions() {
	ctx := context.Background()

	// 1. Выбор завершенных аукционов без победителя
	completedAuctions, err := w.service.GetCompletedAuctionsWithoutWinner(ctx)
	if err != nil {
		w.logger.Printf("Error retrieving completed auctions: %v", err)
		return
	}

	for _, auction := range completedAuctions {
		// 2. Получение всех участников аукциона
		bids, err := w.service.GetBidsByAuctionID(ctx, auction.AuctionID)
		if err != nil {
			w.logger.Printf("Error retrieving bids for auction %d: %v", auction.AuctionID, err)
			continue
		}

		if len(bids) == 0 {
			w.logger.Printf("No bids found for auction %d", auction.AuctionID)
			continue
		}

		// 3. Определение победителя
		winnerID, losers, err := w.service.DetermineWinner(ctx, bids)
		if err != nil {
			w.logger.Printf("Error determining winner for auction %d: %v", auction.AuctionID, err)
			continue
		}

		// 4. Обработка транзакций
		err = w.service.ProcessTransactions(ctx, auction.AuctionID, winnerID, losers)
		if err != nil {
			w.logger.Printf("Error processing transactions for auction %d: %v", auction.AuctionID, err)
			continue
		}

		// 5. Уведомление победителя и проигравших
		err = w.service.NotifyAuctionResults(ctx, auction.AuctionID, winnerID, losers)
		if err != nil {
			w.logger.Printf("Error notifying auction results for auction %d: %v", auction.AuctionID, err)
			continue
		}
	}
}

func (w *AuctionWorker) processNewAuctions() {
	ctx := context.Background()

	err := w.service.NotifyUsersAboutNewAuctions(ctx)
	if err != nil {
		w.logger.Printf("Error notifying users about new auctions: %v", err)
		return
	}

	w.logger.Println("Successfully notified users about new auctions")
}
