package merchants

import "time"

var (
	MinName = 2
	MaxName = 30
)

type MerchantCategory string

var (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

var MerchantCategories = []interface{}{
	SmallRestaurant,
	MediumRestaurant,
	LargeRestaurant,
	MerchandiseRestaurant,
	BoothKiosk,
	ConvenienceStore,
}

type Merchants struct {
	ID        uint64
	UID       string
	Name      string
	Category  MerchantCategory
	ImageURL  string
	Lat       float64
	Lng       float64
	CreatedAt time.Time
}

type MerchantsWithItem struct {
	ID        uint64
	UID       string
	Name      string
	Category  MerchantCategory
	ImageURL  string
	Lat       float64
	Lng       float64
	CreatedAt time.Time
	Item      MerchantItems
}
