package ingredient

type Ingredient struct {
	Id   int64    `db:"ID"`
	Name [50]byte `db:"NAME"`
}
