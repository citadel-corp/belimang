package merchantitems

import "time"

type MerchantItemUIDResponse struct {
	UID string `json:"itemId"`
}

type MerchantItemResponse struct {
	UID             string       `json:"itemId"`
	Name            string       `json:"name"`
	ProductCategory ItemCategory `json:"productCategory"`
	Price           int          `json:"price"`
	ImageURL        string       `json:"imageUrl"`
	CreatedAt       time.Time    `json:"createdAt"`
}

func CreateMerchantItemListResponse(items []MerchantItems) []MerchantItemResponse {
	itemsResponse := make([]MerchantItemResponse, 0)
	for _, item := range items {
		itemsResponse = append(itemsResponse, MerchantItemResponse{
			UID:             item.UID,
			Name:            item.Name,
			ProductCategory: item.Category,
			Price:           item.Price,
			ImageURL:        item.ImageURL,
			CreatedAt:       item.CreatedAt,
		})
	}
	return itemsResponse
}
