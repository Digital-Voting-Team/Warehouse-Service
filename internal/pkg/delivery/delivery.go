package delivery

import "time"

type Delivery struct {
	Id            int64     `db:"ID"`
	SourceId      int64     `db:"SOURCE_ID"`
	DestinationId int64     `db:"DESTINATION_ID"`
	DeliveryPrice float64   `db:"DELIVERY_PRICE"`
	DeliveryDate  time.Time `db:"DELIVERY_DATE"`
}
