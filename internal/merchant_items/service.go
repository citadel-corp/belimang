package merchantitems

import (
	"context"

	"github.com/citadel-corp/belimang/internal/common/id"
	"github.com/citadel-corp/belimang/internal/common/response"
	"github.com/citadel-corp/belimang/internal/merchants"
)

type Service interface {
	Create(ctx context.Context, payload CreateMerchantItemPayload) (resp *MerchantItemUIDResponse, err error)
	List(ctx context.Context, payload ListMerchantItemsPayload) (resp []MerchantItemResponse, pagination *response.Pagination, err error)
}

type merchantItemService struct {
	repository         Repository
	merchantRepository merchants.Repository
}

func NewService(repository Repository, merchantRepository merchants.Repository) Service {
	return &merchantItemService{repository: repository, merchantRepository: merchantRepository}
}

func (s *merchantItemService) Create(ctx context.Context, payload CreateMerchantItemPayload) (resp *MerchantItemUIDResponse, err error) {
	// get merchant
	merchant, err := s.merchantRepository.GetByUID(ctx, payload.MerchantID)
	if err != nil {
		return
	}

	item := &MerchantItems{
		UID:        id.GenerateStringID(16),
		MerchantID: merchant.ID,
		Name:       payload.Name,
		Category:   ItemCategory(payload.ProductCategory),
		Price:      payload.Price,
		ImageURL:   payload.ImageURL,
	}
	err = s.repository.Create(ctx, item)
	if err != nil {
		return
	}

	resp = &MerchantItemUIDResponse{
		UID: item.UID,
	}

	return
}

func (s *merchantItemService) List(ctx context.Context, payload ListMerchantItemsPayload) (resp []MerchantItemResponse, pagination *response.Pagination, err error) {
	// get merchant
	merchant, err := s.merchantRepository.GetByUID(ctx, payload.MerchantUID)
	if err != nil {
		return
	}
	payload.MerchantID = merchant.ID

	if payload.Limit == 0 {
		payload.Limit = 5
	}

	merchantItems, pagination, err := s.repository.List(ctx, payload)
	if err != nil {
		return
	}

	return CreateMerchantItemListResponse(merchantItems), pagination, nil
}
