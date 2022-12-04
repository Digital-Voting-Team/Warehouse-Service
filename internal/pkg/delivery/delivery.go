package delivery

import (
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

type Query interface {
	New() Query

	Get() (*Delivery, error)
	Select() ([]Delivery, error)

	Insert(Delivery) (Delivery, error)
	Update(Delivery) (Delivery, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterBySourceId(ids ...int64) Query
	FilterByDestinationId(ids ...int64) Query
	FilterByPrice(prices ...float64) Query
	FilterByDate(dates ...time.Time) Query
}

type Delivery struct {
	Id            int64     `db:"ID" structs:"-"`
	SourceId      int64     `db:"SOURCE_ID" structs:"source_id"`
	DestinationId int64     `db:"DESTINATION_ID" structs:"destination_id"`
	Price         float64   `db:"PRICE" structs:"price"`
	Date          time.Time `db:"DATE" structs:"date"`
}
