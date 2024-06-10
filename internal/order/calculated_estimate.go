package order

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type CalculatedEstimate struct {
	ID                    string
	UserID                string
	TotalPrice            int
	Lat                   float64
	Long                  float64
	Merchants             CalculatedEstimateMerchants
	Items                 Items
	EstimatedDeliveryTime int
	Ordered               bool
	CreatedAt             time.Time
}

type CalculatedEstimateMerchants []string

type Item struct {
	ItemID     string `json:"itemId"`
	MerchantID string `json:"merchantId"`
	Quantity   int    `json:"quantity"`
}

type Items []Item

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a CalculatedEstimateMerchants) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *CalculatedEstimateMerchants) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

// Make the Attrs struct implement the driver.Valuer interface. This method
// simply returns the JSON-encoded representation of the struct.
func (a Items) Value() (driver.Value, error) {
	return json.Marshal(a)
}

// Make the Attrs struct implement the sql.Scanner interface. This method
// simply decodes a JSON-encoded value into the struct fields.
func (a *Items) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}
