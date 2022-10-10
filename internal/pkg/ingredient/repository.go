package ingredient

import (
	"github.com/jmoiron/sqlx"
)

var (
	queryInsertIngredient = `INSERT INTO WAREHOUSE.INGREDIENT (NAME)
					  	  VALUES (:1) RETURNING ID INTO :2`
	querySelectIngredient = `SELECT * FROM WAREHOUSE.INGREDIENT WHERE ID=:1`
	queryDeleteIngredient = `DELETE FROM WAREHOUSE.INGREDIENT WHERE ID = :1`
	queryGetIngredientId  = `SELECT ID FROM WAREHOUSE.INGREDIENT WHERE NAME = :1`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ingredient *Ingredient) (int64, error) {
	id, _ := r.GetId(ingredient)

	if id == -1 {
		result, err := r.db.Queryx(queryInsertIngredient, ingredient.Name, &ingredient.Id)

		if err != nil {
			return -1, err
		}

		defer result.Close()

		return ingredient.Id, nil
	}

	return id, nil
}

func (r *Repository) Select(id int64) (*Ingredient, error) {
	result, err := r.db.Queryx(querySelectIngredient, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	ingredient := new(Ingredient)

	result.Next()
	err = result.StructScan(ingredient)

	if err != nil {
		return nil, err
	}

	return ingredient, nil
}

func (r *Repository) GetId(ingredient *Ingredient) (int64, error) {
	result, err := r.db.Queryx(queryGetIngredientId, ingredient.Name)

	if err != nil {
		return -1, err
	}

	defer result.Close()

	result.Next()
	err = result.Scan(&ingredient.Id)

	if err != nil {
		return -1, err
	}

	return ingredient.Id, nil
}

func (r *Repository) Update(ingredient *Ingredient) (int64, error) {
	id, _ := r.GetId(ingredient)

	if id == -1 {
		return r.Insert(ingredient)
	}

	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(queryDeleteIngredient, id)

	if err != nil {
		return err
	}

	return nil
}
