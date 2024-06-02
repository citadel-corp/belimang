package merchants

import "time"

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
	CreatedAt time.Time        `json:"createdAt"`
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
			CreatedAt: m.CreatedAt,
		})
	}

	return merchantsResponse
}
