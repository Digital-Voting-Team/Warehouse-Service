package warehouse_ingredient

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

type Query interface {
	New() Query

	Get() (*WarehouseIngredient, error)
	Select() ([]WarehouseIngredient, error)

	Insert(WarehouseIngredient) (WarehouseIngredient, error)
	Update(WarehouseIngredient) (WarehouseIngredient, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterByIngredientId(ids ...int64) Query
	FilterByWarehouseId(ids ...int64) Query
	FilterByQuantity(quantities ...int64) Query
	FilterByOrigin(origins ...string) Query
	FilterByPrice(prices ...float64) Query
	FilterByExpirationDate(expirationDates ...time.Time) Query
	FilterByDeliveryId(ids ...int64) Query
}
type WarehouseIngredient struct {
	Id             int64     `db:"ID" structs:"-"`
	IngredientId   int64     `db:"INGREDIENT_ID" structs:"ingredient_id"`
	WarehouseId    int64     `db:"WAREHOUSE_ID" structs:"warehouse_id"`
	Quantity       int64     `db:"QUANTITY" structs:"quantity"`
	Origin         string    `db:"ORIGIN" structs:"origin"`
	Price          float64   `db:"PRICE" structs:"price"`
	ExpirationDate time.Time `db:"EXPIRATION_DATE" structs:"expiration_date"`
	DeliveryId     int64     `db:"DELIVERY_ID" structs:"delivery_id"`
}
