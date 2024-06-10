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
	merchantListToVisit := make([]*merchants.Merchants, 0)
	for _, merchant := range merchantList {
		if merchant.UID == startingMerchantID {
			startingPoint = haversine.NewCoordinates(merchant.Lat, merchant.Lng)
		} else {
			merchantPoints = append(merchantPoints, haversine.NewCoordinates(merchant.Lat, merchant.Lng))
			merchantListToVisit = append(merchantListToVisit, merchant)
			visited[merchant.UID] = false
		}
	}
	if IsMoreThan3KM2(endPoint, merchantList) {
		return 0, ErrDistanceTooFar
	}

	numPoints := len(merchantListToVisit)
	i := 0
	currDist := 0.0
	point := startingPoint
	for i < numPoints {
		points := GetPointsToCalculate(merchantListToVisit, visited)
		merchant, dist := NearestNeighbor(point, points)
		visited[merchant.UID] = true
		currDist += dist
		point = haversine.NewCoordinates(merchant.Lat, merchant.Lng)
		i += 1
	}
	lastDist := haversine.Distance(point, endPoint).Kilometers()
	currDist += lastDist
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
			point,
			haversine.NewCoordinates(merchant.Lat, merchant.Lng),
		).Kilometers()
		if d < dist {
			dist = d
			res = merchant
		}
	}
	return res, dist
}

func FarthestNeighbor(point haversine.Coordinates, merchantList []*merchants.Merchants) (*merchants.Merchants, float64) {
	var res *merchants.Merchants
	dist := -math.MaxFloat64
	for _, merchant := range merchantList {
		d := haversine.Distance(
			point,
			haversine.NewCoordinates(merchant.Lat, merchant.Lng),
		).Kilometers()
		if d > dist {
			dist = d
			res = merchant
		}
	}
	return res, dist
}

func IsMoreThan3KM2(point haversine.Coordinates, merchantList []*merchants.Merchants) bool {
	_, dist := FarthestNeighbor(point, merchantList)
	circleArea := math.Pi * dist * dist
	if circleArea > 3 {
		return true
	}
	return false
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
