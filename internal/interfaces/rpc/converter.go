package rpc

import (
	"auction/internal/domain"
	v1 "auction/internal/interfaces/rpc/pb"
	"strconv"
)

func NewDomainLotFromRequest(req *v1.CreateLotRequest) domain.Lot {
	userID, _ := strconv.Atoi(req.UserId)
	return domain.Lot{
		Title:      req.Title,
		StartPrice: int(req.StartPrice),
		Step:       int(req.Step),
		UserID:     userID,
	}
}

func NewDomainBidFromRequest(req *v1.PlaceBidRequest) domain.Bid {
	userID, _ := strconv.Atoi(req.UserId)
	lotID, _ := strconv.Atoi(req.LotId)
	return domain.Bid{
		LotID:  lotID,
		UserID: userID,
		Price:  req.Amount,
	}
}
