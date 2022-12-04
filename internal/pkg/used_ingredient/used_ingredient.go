package ingredient

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

type Query interface {
	New() Query

	Get() (*UsedIngredient, error)
	Select() ([]UsedIngredient, error)

	Insert(UsedIngredient) (UsedIngredient, error)
	Update(UsedIngredient) (UsedIngredient, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterByName(names ...string) Query
	FilterByQuantity(quantities ...int64) Query
	FilterByOrigin(origins ...string) Query
	FilterByPrice(prices ...float64) Query
	FilterByDeletionDate(dates ...time.Time) Query
	FilterByReason(reasons ...string) Query
}

type UsedIngredient struct {
	Id           int64     `db:"ID" structs:"-"`
	Name         string    `db:"NAME" structs:"name"`
	Quantity     int64     `db:"QUANTITY" structs:"quantity"`
	Origin       string    `db:"ORIGIN" structs:"origin"`
	Price        float64   `db:"PRICE" structs:"price"`
	DeletionDate time.Time `db:"DELETION_DATE" structs:"deletion_date"`
	Reason       string    `db:"REASON" structs:"reason"`
}
