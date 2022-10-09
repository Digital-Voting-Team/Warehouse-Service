package warehouse

type Warehouse struct {
	Id        int64 `id:"ID"`
	CafeId    int64 `id:"CAFE_ID"`
	AddressId int64 `id:"ADDRESS_ID"`
	Capacity  int64 `id:"CAPACITY"`
}
