package merchantitems

import "time"

var (
	MinName = 2
	MaxName = 30
)

type ItemCategory string

var (
	Beverage   ItemCategory = "Beverage"
	Food       ItemCategory = "Food"
	Snack      ItemCategory = "Snack"
	Condiments ItemCategory = "Condiments"
	Additions  ItemCategory = "Additions"
)

var ProductCategories = []interface{}{Beverage, Food, Snack, Condiments, Additions}

type MerchantItems struct {
	ID          uint64
	UID         string
	MerchantID  uint64
	MerchantUID string
	Name        string
	Category    ItemCategory
	Price       int
	ImageURL    string
	CreatedAt   time.Time
}
