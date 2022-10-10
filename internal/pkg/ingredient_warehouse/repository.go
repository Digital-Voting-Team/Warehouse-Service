package ingredient_warehouse

import (
	"github.com/jmoiron/sqlx"
)

var (
	queryInsertIngredientWarehouse = `INSERT INTO WAREHOUSE.INGREDIENT_WAREHOUSE (INGREDIENT_ID, WAREHOUSE_ID, QUANTITY, ORIGIN, PRICE, EXPIRATION_DATE, DELIVERY_DATE)
					  	 			  VALUES (:1, :2, :3, :4, :5, :6, :7) RETURNING ID INTO :8`
	querySelectIngredientWarehouse = `SELECT * FROM WAREHOUSE.INGREDIENT_WAREHOUSE WHERE ID=:1`
	queryDeleteIngredientWarehouse = `DELETE FROM WAREHOUSE.INGREDIENT_WAREHOUSE WHERE ID = :1`
	queryGetIngredientWarehouseId  = `SELECT ID FROM WAREHOUSE.INGREDIENT_WAREHOUSE WHERE INGREDIENT_ID = :1 AND
																						  WAREHOUSE_ID = :2 AND
																						  QUANTITY = :3 AND
																						  ORIGIN = :4 AND
																						  PRICE = :5 AND
																						  EXPIRATION_DATE = :6 AND
																						  DELIVERY_DATE = :7`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(ingredientWarehouse *IngredientWarehouse) (int64, error) {
	id, _ := r.GetId(ingredientWarehouse)

	if id == -1 {
		result, err := r.db.Queryx(queryInsertIngredientWarehouse, ingredientWarehouse.IngredientId,
			ingredientWarehouse.WarehouseId, ingredientWarehouse.Quantity, ingredientWarehouse.Origin,
			ingredientWarehouse.Price, ingredientWarehouse.ExpirationDate, ingredientWarehouse.DeliveryDate,
			&ingredientWarehouse.Id)

		if err != nil {
			return -1, err
		}

		defer result.Close()

		return ingredientWarehouse.Id, nil
	}

	return id, nil
}

func (r *Repository) Select(id int64) (*IngredientWarehouse, error) {
	result, err := r.db.Queryx(querySelectIngredientWarehouse, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	warehouse := new(IngredientWarehouse)

	result.Next()
	err = result.StructScan(warehouse)

	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (r *Repository) GetId(ingredientWarehouse *IngredientWarehouse) (int64, error) {
	result, err := r.db.Queryx(queryGetIngredientWarehouseId, ingredientWarehouse.IngredientId,
		ingredientWarehouse.WarehouseId, ingredientWarehouse.Quantity, ingredientWarehouse.Origin,
		ingredientWarehouse.Price, ingredientWarehouse.ExpirationDate, ingredientWarehouse.DeliveryDate)

	if err != nil {
		return -1, err
	}

	defer result.Close()

	result.Next()
	err = result.Scan(&ingredientWarehouse.Id)

	if err != nil {
		return -1, err
	}

	return ingredientWarehouse.Id, nil
}

func (r *Repository) Update(warehouse *IngredientWarehouse) (int64, error) {
	id, _ := r.GetId(warehouse)

	if id == -1 {
		return r.Insert(warehouse)
	}

	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(queryDeleteIngredientWarehouse, id)

	if err != nil {
		return err
	}

	return nil
}
