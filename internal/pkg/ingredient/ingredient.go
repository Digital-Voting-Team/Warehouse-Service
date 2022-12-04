package ingredient

import "gitlab.com/distributed_lab/kit/pgdb"

type Query interface {
	New() Query

	Get() (*Ingredient, error)
	Select() ([]Ingredient, error)

	Insert(Ingredient) (Ingredient, error)
	Update(Ingredient) (Ingredient, error)
	Delete(id int64) error

	Page(pageParams pgdb.OffsetPageParams) Query

	FilterById(ids ...int64) Query
	FilterByName(names ...string) Query
}

type Ingredient struct {
	Id   int64  `db:"ID" structs:"-"`
	Name string `db:"NAME" structs:"name"`
}
