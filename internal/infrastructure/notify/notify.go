package notify

import (
	"auction/internal/domain"
	"auction/internal/infrastructure/repo"
	"context"
	"fmt"
)

type NotifyService interface {
	NotifyUser(ctx context.Context, userID int, message string) error
	NotifyAllUsersAboutNewAuctions(ctx context.Context, auctions []domain.Auction) error
}

type notifyService struct {
	userRepo repo.UserRepository
}

func NewNotifyService(userRepo repo.UserRepository) NotifyService {
	return &notifyService{
		userRepo: userRepo,
	}
}

func (n *notifyService) NotifyUser(ctx context.Context, userID int, message string) error {
	// Логика отправки уведомления пользователю
	// Например, отправка через email или push-уведомления
	fmt.Printf("Sending notification to user %d: %s\n", userID, message)
	return nil
}

func (s *notifyService) NotifyAllUsersAboutNewAuctions(ctx context.Context, auctions []domain.Auction) error {
	users, err := s.userRepo.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	auctionList := ""
	for _, auction := range auctions {
		auctionList += fmt.Sprintf("Auction ID: %d", auction.AuctionID)
	}

	message := fmt.Sprintf("New auctions have started:\n%s", auctionList)

	for _, user := range users {
		err := s.sendNotification(ctx, user.ID, message)
		if err != nil {
			fmt.Printf("Error sending notification to user %d: %v", user.ID, err)
		}
	}

	return nil
}

func (s *notifyService) sendNotification(ctx context.Context, userID int, message string) error {
	fmt.Printf("Sending notification to user %d: %s\n", userID, message)
	return nil
}
