package app

import (
	"auction/internal/domain"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetermineWinner(t *testing.T) {
	service := &AuctionService{}

	tests := []struct {
		name       string
		bids       []domain.Bid
		wantWinner int
		wantLosers []int
		wantErr    error
	}{
		{
			name:       "No bids available",
			bids:       []domain.Bid{},
			wantWinner: 0,
			wantLosers: nil,
			wantErr:    errors.New("no bids available"),
		},
		{
			name: "Single bid",
			bids: []domain.Bid{
				{UserID: 1, Price: 100},
			},
			wantWinner: 1,
			wantLosers: []int{},
			wantErr:    nil,
		},
		{
			name: "Multiple bids with distinct prices",
			bids: []domain.Bid{
				{UserID: 1, Price: 100},
				{UserID: 2, Price: 200},
				{UserID: 3, Price: 150},
			},
			wantWinner: 2,
			wantLosers: []int{1, 3},
			wantErr:    nil,
		},
		{
			name: "Multiple bids with same highest price",
			bids: []domain.Bid{
				{UserID: 1, Price: 200},
				{UserID: 2, Price: 200},
				{UserID: 3, Price: 150},
			},
			wantWinner: 1, // Assumes first highest bid wins
			wantLosers: []int{2, 3},
			wantErr:    nil,
		},
		{
			name: "All bids from the same user",
			bids: []domain.Bid{
				{UserID: 1, Price: 100},
				{UserID: 1, Price: 200},
				{UserID: 1, Price: 150},
			},
			wantWinner: 1,
			wantLosers: []int{},
			wantErr:    nil,
		},
		{
			name: "Bids with zero and negative prices",
			bids: []domain.Bid{
				{UserID: 1, Price: -50},
				{UserID: 2, Price: 0},
				{UserID: 3, Price: 100},
			},
			wantWinner: 3,
			wantLosers: []int{1, 2},
			wantErr:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winner, losers, err := service.DetermineWinner(context.Background(), tt.bids)

			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantWinner, winner)
				assert.ElementsMatch(t, tt.wantLosers, losers)
			}
		})
	}
}
