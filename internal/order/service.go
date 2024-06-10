package order

import (
	"context"
	"fmt"
	"math"

	"github.com/LucaTheHacker/go-haversine"
	"github.com/citadel-corp/belimang/internal/common/id"
	merchantitems "github.com/citadel-corp/belimang/internal/merchant_items"
	"github.com/citadel-corp/belimang/internal/merchants"
)

type Service interface {
	CalculateEstimate(ctx context.Context, req CalculateOrderEstimateRequest, userID string) (*CalculateOrderEstimateResponse, error)
	CreateOrder(ctx context.Context, req CreateOrderRequest, userID string) (*CreateOrderResponse, error)
	SearchOrder(ctx context.Context, req SearchOrderPayload, userID string) ([]*SearchOrderResponse, error)
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
	var startingPoint haversine.Coordinates
	endPoint := haversine.NewCoordinates(req.UserLocation.Lat, req.UserLocation.Long)
	merchantPoints := make([]haversine.Coordinates, 0)
	visited := make(map[string]bool) // string: merchant id, bool: has visited
	for _, merchant := range merchantList {
		if merchant.UID == startingMerchantID {
			startingPoint = haversine.NewCoordinates(merchant.Lat, merchant.Lng)
		} else {
			merchantPoints = append(merchantPoints, haversine.NewCoordinates(merchant.Lat, merchant.Lng))
			visited[merchant.UID] = false
		}
	}
	deliveryTime, err := calculateDeliveryTime(startingPoint, endPoint, merchantList, visited)
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

// SearchOrder implements Service.
func (s *orderService) SearchOrder(ctx context.Context, req SearchOrderPayload, userID string) ([]*SearchOrderResponse, error) {
	orderItemMerchants, err := s.repository.SearchOrderItemMerchants(ctx, req, userID)
	if err != nil {
		return nil, err
	}
	actualMerchantIDs := make([]string, 0)
	for _, orderItemMerchant := range orderItemMerchants {
		actualMerchantIDs = append(actualMerchantIDs, orderItemMerchant.MerchantID)
	}
	// searchResults, err := s.repository.SearchOrder(ctx, req, userID)\
	orderItemDetailsMap := make(map[string][]Item) // key: orderID-merchantID
	for _, orderItemMerchant := range orderItemMerchants {
		key := fmt.Sprintf("%s-%s", orderItemMerchant.OrderID, orderItemMerchant.MerchantID)
		orderItemDetailsMap[key] = append(orderItemDetailsMap[key], orderItemMerchant.OrderItems)
	}
	items, err := s.merchantItemsRepository.ListByMerchantUIDAndName(ctx, actualMerchantIDs, req.Name)
	if err != nil {
		return nil, err
	}
	itemsMap := make(map[string]*merchantitems.MerchantItems) // key: itemID
	for _, item := range items {
		itemsMap[item.UID] = item
	}

	orderItemMerchantsMap := make(map[string][]*SearchOrderResponse) // key: orderID
	res := make([]*SearchOrderResponse, 0)
	for _, orderItemMerchant := range orderItemMerchants {
		searchOrderDetailItemResponse := make([]SearchOrderDetailItemResponse, 0)
		key := fmt.Sprintf("%s-%s", orderItemMerchant.OrderID, orderItemMerchant.MerchantID)
		searchOrderDetailItems := orderItemDetailsMap[key]
		for _, searchOrderDetailItem := range searchOrderDetailItems {
			item := itemsMap[searchOrderDetailItem.ItemID]
			searchOrderDetailItemResponse = append(searchOrderDetailItemResponse, SearchOrderDetailItemResponse{
				MerchantItemResponse: merchantitems.MerchantItemResponse{
					UID:             item.UID,
					Name:            item.Name,
					ProductCategory: item.Category,
					Price:           item.Price,
					ImageURL:        item.ImageURL,
					CreatedAt:       item.CreatedAt,
				},
				Quantity: searchOrderDetailItem.Quantity,
			})
		}
		orderItemMerchantsMap[orderItemMerchant.OrderID] = append(orderItemMerchantsMap[orderItemMerchant.OrderID], &SearchOrderResponse{
			OrderID: orderItemMerchant.OrderID,
			Orders: SearchOrderDetailResponse{
				Merchant: merchants.MerchantsResponse{
					UID:      orderItemMerchant.MerchantID,
					Name:     orderItemMerchant.MerchantName,
					Category: orderItemMerchant.MerchantCategory,
					ImageURL: orderItemMerchant.MerchantImageURL,
					Location: merchants.LocationResponse{
						Lat: orderItemMerchant.MerchantLat,
						Lng: orderItemMerchant.MerchantLong,
					},
					CreatedAt: orderItemMerchant.MerchantCreatedAt,
				},
				Items: searchOrderDetailItemResponse,
			},
		})
		// if orderID, ok := orderItemMerchantsMap[orderItemMerchant.OrderID]; ok {

		// } else {

		// }
	}
	// for k, v := range orderItemMerchantsMap {
	// 	res = append(res, v)
	// }
	return res, nil
}

func calculateDeliveryTime(startingPoint, endPoint haversine.Coordinates, merchantList []*merchants.Merchants, visited map[string]bool) (int, error) {
	numPoints := len(merchantList)
	i := 0
	currDist := 0.0
	point := startingPoint
	for i < numPoints {
		points := getPointsToCalculate(merchantList, visited)
		merchant, dist := nearestNeighbor(point, points)
		visited[merchant.UID] = true
		currDist += dist
		if currDist > 3.0 {
			return 0, ErrDistanceTooFar
		}
		point = haversine.NewCoordinates(merchant.Lat, merchant.Lng)
		i += 1
	}
	currDist += haversine.Distance(
		haversine.NewCoordinates(point.Latitude, point.Longitude),
		haversine.NewCoordinates(endPoint.Latitude, endPoint.Longitude),
	).Kilometers()
	speedInMS := 11.11 // m/s
	currDist *= 1000   // convert to meter
	timeSecond := currDist / speedInMS
	return int(timeSecond / 60), nil
}

func nearestNeighbor(point haversine.Coordinates, merchantList []*merchants.Merchants) (*merchants.Merchants, float64) {
	var res *merchants.Merchants
	dist := math.MaxFloat64
	for _, merchant := range merchantList {
		d := haversine.Distance(
			haversine.NewCoordinates(point.Latitude, point.Longitude),
			haversine.NewCoordinates(merchant.Lat, merchant.Lng),
		).Kilometers()
		if d < dist {
			dist = d
			res = merchant
		}
	}
	return res, dist
}

func getPointsToCalculate(merchantList []*merchants.Merchants, visited map[string]bool) []*merchants.Merchants {
	res := make([]*merchants.Merchants, 0)
	for _, merchant := range merchantList {
		if !visited[merchant.UID] {
			res = append(res, merchant)
		}
	}
	return res
}
