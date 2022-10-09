package ingredient_warehouse

import "time"

type IngredientWarehouse struct {
	Id             int64     `db:"ID"`
	IngredientId   int64     `db:"INGREDIENT_ID"`
	WarehouseId    int64     `db:"WAREHOUSE_ID"`
	Quantity       int64     `db:"QUANTITY"`
	Origin         [50]byte  `db:"ORIGIN"`
	Price          float64   `db:"PRICE"`
	ExpirationDate time.Time `db:"EXPIRATION_DATE"`
	DeliveryDate   time.Time `db:"DELIVERY_DATE"`
}
