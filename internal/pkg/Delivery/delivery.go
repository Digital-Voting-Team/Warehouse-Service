package Delivery

import "time"

type Delivery struct {
	Id            int64     `id:"ID"`
	SourceId      int64     `id:"SOURCE_ID"`
	DestinationId int64     `id:"DESTINATION_ID"`
	DeliveryPrice float64   `id:"DELIVERY_PRICE"`
	DeliveryDate  time.Time `id:"DELIVERY_DATE"`
}
