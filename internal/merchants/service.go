package merchants

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/id"
)

type Service interface {
	Create(ctx context.Context, req CreateMerchantPayload) (*MerchantUIDResponse, error)
	List(ctx context.Context, req ListMerchantsPayload) ([]MerchantsResponse, error)
	ListByDistance(ctx context.Context, req ListMerchantsByDistancePayload) ([]MerchantsResponse, error)
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
		Lat:      req.Location.Lat,
		Lng:      req.Location.Lng,
	}
	err := s.repository.Create(ctx, merchant)
	if err != nil {
		return nil, err
	}

	return &MerchantUIDResponse{
		UID: merchant.UID,
	}, nil
}

func (s *merchantService) List(ctx context.Context, req ListMerchantsPayload) ([]MerchantsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	merchants, err := s.repository.List(ctx, req)
	if err != nil {
		return []MerchantsResponse{}, err
	}

	return CreateMerchantsResponse(merchants), nil
}

func (s *merchantService) ListByDistance(ctx context.Context, req ListMerchantsByDistancePayload) ([]MerchantsResponse, error) {
	if req.Limit == 0 {
		req.Limit = 5
	}

	merchants, err := s.repository.ListByDistance(ctx, req)
	if err != nil {
		return []MerchantsResponse{}, err
	}

	return CreateMerchantsResponse(merchants), nil
}
