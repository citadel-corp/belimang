package order

import (
	"context"
	"fmt"

	"github.com/citadel-corp/belimang/internal/common/haversine"
	"github.com/citadel-corp/belimang/internal/common/id"
	merchantitems "github.com/citadel-corp/belimang/internal/merchant_items"
	"github.com/citadel-corp/belimang/internal/merchants"
)

type Service interface {
	CalculateEstimate(ctx context.Context, req CalculateOrderEstimateRequest, userID string) (*CalculateOrderEstimateResponse, error)
	CreateOrder(ctx context.Context, req CreateOrderRequest, userID string) (*CreateOrderResponse, error)
}

type orderService struct {
	repository              Repository
	merchantRepository      merchants.Repository
	merchantItemsRepository merchantitems.Repository
}

func NewService(repository Repository, merchantRepository merchants.Repository, merchantItemsRepository merchantitems.Repository) Service {
	return &orderService{
		repository:              repository,
		merchantRepository:      merchantRepository,
		merchantItemsRepository: merchantItemsRepository,
	}
}

// CalculateEstimate implements Service.
func (s *orderService) CalculateEstimate(ctx context.Context, req CalculateOrderEstimateRequest, userID string) (*CalculateOrderEstimateResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	// validate only 1 starting point
	startingPointCount := 0
	merchantIDs := make([]string, len(req.Orders))
	startingMerchantID := ""
	allItems := make([]OrderItemRequest, 0)
	calculateEstimateItems := make(Items, 0)
	for i, order := range req.Orders {
		if order.IsStartingPoint {
			startingPointCount += 1
			startingMerchantID = order.MerchantID
		}
		merchantIDs[i] = order.MerchantID
		for _, item := range order.Items {
			allItems = append(allItems, item)
			calculateEstimateItems = append(calculateEstimateItems, Item{
				ItemID:     item.ItemID,
				MerchantID: order.MerchantID,
				Quantity:   item.Quantity,
			})
		}
	}
	if startingPointCount != 1 {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, ErrStartingPointInvalid)
	}
	merchantList, err := s.merchantRepository.ListByUIDs(ctx, merchantIDs)
	if err != nil {
		return nil, err
	}
	if len(req.Orders) != len(merchantList) {
		return nil, ErrSomeMerchantNotFound
	}
	itemList, err := s.merchantItemsRepository.ListByUIDs(ctx, merchantIDs)
	if err != nil {
		return nil, err
	}
	if len(allItems) != len(itemList) {
		return nil, ErrSomeItemNotFound
	}
	totalPrice := 0
	for _, item := range itemList {
		totalPrice += item.Price
	}
	// calculate delivery time
	deliveryTime, err := haversine.CalculateDeliveryTime(req.UserLocation.Lat, req.UserLocation.Long, startingMerchantID, merchantList)
	if err != nil {
		return nil, err
	}
	calculatedEstimate := &CalculatedEstimate{
		ID:                    id.GenerateStringID(16),
		UserID:                userID,
		TotalPrice:            totalPrice,
		Lat:                   req.UserLocation.Lat,
		Long:                  req.UserLocation.Long,
		Merchants:             CalculatedEstimateMerchants(merchantIDs),
		Items:                 calculateEstimateItems,
		EstimatedDeliveryTime: deliveryTime,
		Ordered:               false,
	}
	err = s.repository.InsertCalculatedEstimate(ctx, calculatedEstimate)
	if err != nil {
		return nil, err
	}

	return &CalculateOrderEstimateResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: deliveryTime,
		CalculatedEstimateID:           calculatedEstimate.ID,
	}, nil
}

// CreateOrder implements Service.
func (s *orderService) CreateOrder(ctx context.Context, req CreateOrderRequest, userID string) (*CreateOrderResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrValidationFailed, err)
	}
	calculatedEstimate, err := s.repository.GetCalculatedEstimate(ctx, req.CalculatedEstimateID)
	if err != nil {
		return nil, err
	}
	order := &Order{
		ID:                   id.GenerateStringID(16),
		CalculatedEstimateID: calculatedEstimate.ID,
		UserID:               userID,
	}
	err = s.repository.InsertOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	merchantItemMap := make(map[string][]Item) // key: merchant id
	for _, item := range calculatedEstimate.Items {
		merchantItemMap[item.MerchantID] = append(merchantItemMap[item.MerchantID], item)
	}
	for k, v := range merchantItemMap {
		err = s.repository.InsertOrderItem(ctx, &OrderItem{
			ID:         id.GenerateStringID(16),
			OrderID:    order.ID,
			MerchantID: k,
			Items:      v,
		})
		if err != nil {
			return nil, err
		}
	}

	return &CreateOrderResponse{
		OrderID: order.ID,
	}, nil
}
