package haversine

import (
	"math"

	"github.com/LucaTheHacker/go-haversine"
	"github.com/citadel-corp/belimang/internal/merchants"
)

func CalculateDeliveryTime(lat, lng float64, startingMerchantID string, merchantList []*merchants.Merchants) (int, error) {
	var startingPoint haversine.Coordinates
	endPoint := haversine.NewCoordinates(lat, lng)
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

	numPoints := len(merchantList)
	i := 0
	currDist := 0.0
	point := startingPoint
	for i < numPoints {
		points := GetPointsToCalculate(merchantList, visited)
		merchant, dist := NearestNeighbor(point, points)
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

func NearestNeighbor(point haversine.Coordinates, merchantList []*merchants.Merchants) (*merchants.Merchants, float64) {
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

func GetPointsToCalculate(merchantList []*merchants.Merchants, visited map[string]bool) []*merchants.Merchants {
	res := make([]*merchants.Merchants, 0)
	for _, merchant := range merchantList {
		if !visited[merchant.UID] {
			res = append(res, merchant)
		}
	}
	return res
}
