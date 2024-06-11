package merchants

import (
	"database/sql"
	"time"
)

type MerchantUIDResponse struct {
	UID string `json:"merchantId"`
}

type LocationResponse struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"long"`
}

type MerchantsResponse struct {
	UID       string           `json:"merchantId"`
	Name      string           `json:"name"`
	Category  string           `json:"merchantCategory"`
	ImageURL  string           `json:"imageUrl"`
	Location  LocationResponse `json:"location"`
	CreatedAt int              `json:"createdAt"`
}

func CreateMerchantsResponse(merchants []Merchants) []MerchantsResponse {
	merchantsResponse := make([]MerchantsResponse, 0)
	for _, m := range merchants {
		merchantsResponse = append(merchantsResponse, MerchantsResponse{
			UID:       m.UID,
			Name:      m.Name,
			Category:  string(m.Category),
			ImageURL:  m.ImageURL,
			Location:  LocationResponse{Lat: m.Lat, Lng: m.Lng},
			CreatedAt: m.CreatedAt.Nanosecond(),
		})
	}

	return merchantsResponse
}

type MerchantItemResponse struct {
	UID             string `json:"itemId"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageURL        string `json:"imageUrl"`
	CreatedAt       int    `json:"createdAt"`
}

type MerchantWithItemsResponse struct {
	Merchant MerchantsResponse      `json:"merchant"`
	Items    []MerchantItemResponse `json:"items"`
}

func CreateMerchantsWithItemsResponse(merchants []MerchantsWithItem) []MerchantWithItemsResponse {
	merchantMap := make(map[uint64]*MerchantWithItemsResponse)

	for _, merchant := range merchants {
		if _, exists := merchantMap[merchant.ID]; !exists {
			merchantMap[merchant.ID] = &MerchantWithItemsResponse{
				Merchant: MerchantsResponse{
					UID:      merchant.UID,
					Name:     merchant.Name,
					Category: string(merchant.Category),
					ImageURL: merchant.ImageURL,
					Location: LocationResponse{
						Lat: merchant.Lat,
						Lng: merchant.Lng,
					},
					CreatedAt: merchant.CreatedAt.Nanosecond(),
				},
				Items: []MerchantItemResponse{},
			}
		}

		if merchant.Item.UID == "" {
			continue
		}

		merchantMap[merchant.ID].Items = append(merchantMap[merchant.ID].Items, MerchantItemResponse{
			UID:             merchant.Item.UID,
			Name:            merchant.Item.Name,
			ProductCategory: getString(merchant.Item.Category),
			Price:           merchant.Item.Price,
			ImageURL:        merchant.Item.ImageURL,
			CreatedAt:       getTime(merchant.Item.CreatedAt).Nanosecond(),
		})
	}

	var result []MerchantWithItemsResponse
	for _, merchant := range merchantMap {
		result = append(result, *merchant)
	}

	return result
}

func getString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func getTime(ns sql.NullTime) time.Time {
	if ns.Valid {
		return ns.Time
	}
	return time.Time{}
}
