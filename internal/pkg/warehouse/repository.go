package warehouse

import (
	"github.com/jmoiron/sqlx"
)

var (
	queryInsertWarehouse = `INSERT INTO WAREHOUSE.WAREHOUSE (CAFE_ID, ADDRESS_ID, CAPACITY)
					  	 	VALUES (:1, :2, :3) RETURNING ID INTO :4`
	querySelectWarehouse = `SELECT * FROM WAREHOUSE.WAREHOUSE WHERE ID=:1`
	queryDeleteWarehouse = `DELETE FROM WAREHOUSE.WAREHOUSE WHERE ID = :1`
	queryGetWarehouseId  = `SELECT ID FROM WAREHOUSE.WAREHOUSE WHERE CAFE_ID = :1 AND
																	 ADDRESS_ID = :2 AND
																	 CAPACITY = :3`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(warehouse *Warehouse) (int64, error) {
	id, _ := r.GetId(warehouse)

	if id == -1 {
		result, err := r.db.Queryx(queryInsertWarehouse, warehouse.CafeId,
			warehouse.AddressId, warehouse.Capacity, &warehouse.Id)

		if err != nil {
			return -1, err
		}

		defer result.Close()

		return warehouse.Id, nil
	}

	return id, nil
}

func (r *Repository) Select(id int64) (*Warehouse, error) {
	result, err := r.db.Queryx(querySelectWarehouse, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	warehouse := new(Warehouse)

	result.Next()
	err = result.StructScan(warehouse)

	if err != nil {
		return nil, err
	}

	return warehouse, nil
}

func (r *Repository) GetId(warehouse *Warehouse) (int64, error) {
	result, err := r.db.Queryx(queryGetWarehouseId, warehouse.CafeId,
		warehouse.AddressId, warehouse.Capacity)

	if err != nil {
		return -1, err
	}

	defer result.Close()

	result.Next()
	err = result.Scan(&warehouse.Id)

	if err != nil {
		return -1, err
	}

	return warehouse.Id, nil
}

func (r *Repository) Update(warehouse *Warehouse) (int64, error) {
	id, _ := r.GetId(warehouse)

	if id == -1 {
		return r.Insert(warehouse)
	}

	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(queryDeleteWarehouse, id)

	if err != nil {
		return err
	}

	return nil
}
