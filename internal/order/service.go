package order

import (
	"context"
	"fmt"
	"math"

	"github.com/LucaTheHacker/go-haversine"
	"github.com/citadel-corp/belimang/internal/common/id"
	"github.com/citadel-corp/belimang/internal/merchants"
)

type Service interface {
	CalculateEstimate(ctx context.Context, req CalculateOrderEstimateRequest, userID string) (*CalculateOrderEstimateResponse, error)
}

type orderService struct {
	repository         Repository
	merchantRepository merchants.Repository
}

func NewService(repository Repository, merchantRepository merchants.Repository) Service {
	return &orderService{
		repository:         repository,
		merchantRepository: merchantRepository,
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
	for i, order := range req.Orders {
		if order.IsStartingPoint {
			startingPointCount += 1
			startingMerchantID = order.MerchantID
		}
		merchantIDs[i] = order.MerchantID
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
	// TODO: calculate price
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
		TotalPrice:            0,
		Lat:                   req.UserLocation.Lat,
		Long:                  req.UserLocation.Long,
		EstimatedDeliveryTime: deliveryTime,
		Ordered:               false,
	}
	err = s.repository.InsertCalculatedEstimate(ctx, calculatedEstimate)
	if err != nil {
		return nil, err
	}

	return &CalculateOrderEstimateResponse{
		TotalPrice:                     0,
		EstimatedDeliveryTimeInMinutes: deliveryTime,
		CalculatedEstimateID:           calculatedEstimate.ID,
	}, nil
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
