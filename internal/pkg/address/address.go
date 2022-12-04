package address

import "gitlab.com/distributed_lab/kit/pgdb"

type Query interface {
	New() Query

	Get() (*Address, error)
	Select() ([]Address, error)

	Insert(Address) (Address, error)
	Update(Address) (Address, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterByBuilding(numbers ...int64) Query
	FilterByStreet(streets ...string) Query
	FilterByCity(cities ...string) Query
	FilterByDistrict(districts ...string) Query
	FilterByRegion(regions ...string) Query
	FilterByPostalCode(codes ...string) Query
}

type Address struct {
	Id         int64  `db:"ID" structs:"-"`
	Building   int64  `db:"BUILDING" structs:"building"`
	Street     string `db:"STREET" structs:"street"`
	City       string `db:"CITY" structs:"city"`
	District   string `db:"DISTRICT" structs:"district"`
	Region     string `db:"REGION" structs:"region"`
	PostalCode string `db:"POSTAL_CODE" structs:"postal_code"`
}
