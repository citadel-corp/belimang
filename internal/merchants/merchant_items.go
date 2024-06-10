package merchants

import (
	"database/sql"
)

type MerchantItems struct {
	ID          uint64
	UID         string
	MerchantID  uint64
	MerchantUID string
	Name        string
	Category    sql.NullString
	Price       int
	ImageURL    string
	CreatedAt   sql.NullTime
}
