package warehouse

import "gitlab.com/distributed_lab/kit/pgdb"

type Query interface {
	New() Query

	Get() (*Warehouse, error)
	Select() ([]Warehouse, error)

	Insert(Warehouse) (Warehouse, error)
	Update(Warehouse) (Warehouse, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterByCafeId(ids ...int64) Query
	FilterByAddressId(ids ...int64) Query
	FilterByCapacity(capacities ...int64) Query
}

type Warehouse struct {
	Id        int64 `db:"ID" structs:"-"`
	CafeId    int64 `db:"CAFE_ID" structs:"cafe_id"`
	AddressId int64 `db:"ADDRESS_ID" structs:"address_id"`
	Capacity  int64 `db:"CAPACITY" structs:"capacity"`
}
