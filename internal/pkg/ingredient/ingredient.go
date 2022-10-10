package ingredient

type Ingredient struct {
	Id   int64  `db:"ID"`
	Name string `db:"NAME"`
}
