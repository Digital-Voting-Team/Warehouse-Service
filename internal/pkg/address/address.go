package address

type Address struct {
	Id         int64    `db:"ID"`
	Street     [50]byte `db:"STREET"`
	City       [50]byte `db:"CITY"`
	District   [50]byte `db:"DISTRICT"`
	Region     [50]byte `db:"REGION"`
	PostalCode [10]byte `db:"POSTAL_CODE"`
}
