package payment

import (
	"context"
	"fmt"
)

type BalanceService interface {
	DeductBalance(ctx context.Context, userID int, amount int64) error
	RefundBalance(ctx context.Context, userID int, amount int64) error
}

type balanceService struct {
}

func NewBalanceService() BalanceService {
	return &balanceService{}
}

func (b *balanceService) DeductBalance(ctx context.Context, userID int, amount int64) error {
	// Мок-реализация списания средств с пользователя
	fmt.Printf("Deducting %d from user %d\n", amount, userID)
	return nil
}

func (b *balanceService) RefundBalance(ctx context.Context, userID int, amount int64) error {
	// Мок-реализация возврата средств пользователю
	fmt.Printf("Refunding %d to user %d\n", amount, userID)
	return nil
}
