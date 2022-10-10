package address

type Address struct {
	Id         int64  `db:"ID"`
	Building   int64  `db:"BUILDING"`
	Street     string `db:"STREET"`
	City       string `db:"CITY"`
	District   string `db:"DISTRICT"`
	Region     string `db:"REGION"`
	PostalCode string `db:"POSTAL_CODE"`
}
