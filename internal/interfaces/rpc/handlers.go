package rpc

import (
	"auction/internal/domain"
	v1 "auction/internal/interfaces/rpc/pb"

	"context"
	"log"
	"strconv"
)

type AuctionHandler struct {
	v1.UnimplementedAuctionServiceServer
	auctionService domain.AuctionService
}

func NewAuctionHandler(service domain.AuctionService) *AuctionHandler {
	return &AuctionHandler{auctionService: service}
}

func (h *AuctionHandler) CreateLot(ctx context.Context, req *v1.CreateLotRequest) (*v1.CreateLotResponse, error) {
	lot := NewDomainLotFromRequest(req)
	lotID, err := h.auctionService.CreateLot(ctx, lot)
	if err != nil {
		log.Printf("Error creating lot: %v", err)
		return nil, err
	}

	return &v1.CreateLotResponse{LotId: strconv.Itoa(lotID)}, nil
}

func (h *AuctionHandler) RefillBalance(ctx context.Context, req *v1.RefillRequest) (*v1.RefillResponse, error) {
	userID, _ := strconv.Atoi(req.UserId)
	err := h.auctionService.RefillBalance(ctx, userID, req.Amount)
	if err != nil {
		log.Printf("Error refilling balance: %v", err)
		return nil, err
	}

	return &v1.RefillResponse{Message: "Balance refilled successfully"}, nil
}

func (h *AuctionHandler) PlaceBid(ctx context.Context, req *v1.PlaceBidRequest) (*v1.PlaceBidResponse, error) {
	bid := NewDomainBidFromRequest(req)

	_, err := h.auctionService.PlaceBid(ctx, bid)
	if err != nil {
		log.Printf("Error placing bid: %v", err)
		return nil, err
	}

	return &v1.PlaceBidResponse{Message: "bid placed"}, nil
}
