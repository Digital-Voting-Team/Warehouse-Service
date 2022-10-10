package warehouse

type Warehouse struct {
	Id        int64 `db:"ID"`
	CafeId    int64 `db:"CAFE_ID"`
	AddressId int64 `db:"ADDRESS_ID"`
	Capacity  int64 `db:"CAPACITY"`
}
