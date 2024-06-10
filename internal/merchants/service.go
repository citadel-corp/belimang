package merchants

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/id"
	"github.com/citadel-corp/belimang/internal/common/response"
)

type Service interface {
	Create(ctx context.Context, req CreateMerchantPayload) (*MerchantUIDResponse, error)
	List(ctx context.Context, req ListMerchantsPayload) ([]MerchantsResponse, *response.Pagination, error)
	ListByDistance(ctx context.Context, req ListMerchantsByDistancePayload) ([]MerchantWithItemsResponse, *response.Pagination, error)
}

type merchantService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &merchantService{repository: repository}
}

func (s *merchantService) Create(ctx context.Context, req CreateMerchantPayload) (*MerchantUIDResponse, error) {
	merchant := &Merchants{
		UID:      id.GenerateStringID(16),
		Name:     req.Name,
		Category: req.Category,
		ImageURL: req.ImageURL,
		Lat:      *req.Location.Lat,
		Lng:      *req.Location.Lng,
	}
	err := s.repository.Create(ctx, merchant)
	if err != nil {
		return nil, err
	}

	return &MerchantUIDResponse{
		UID: merchant.UID,
	}, nil
}

func (s *merchantService) List(ctx context.Context, req ListMerchantsPayload) ([]MerchantsResponse, *response.Pagination, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	merchants, pagination, err := s.repository.List(ctx, req)
	if err != nil {
		return []MerchantsResponse{}, nil, err
	}

	return CreateMerchantsResponse(merchants), pagination, nil
}

func (s *merchantService) ListByDistance(ctx context.Context, req ListMerchantsByDistancePayload) ([]MerchantWithItemsResponse, *response.Pagination, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	merchantsWithItem, pagination, err := s.repository.ListByDistance(ctx, req)
	if err != nil {
		return []MerchantWithItemsResponse{}, nil, err
	}

	return CreateMerchantsWithItemsResponse(merchantsWithItem), pagination, nil
}
