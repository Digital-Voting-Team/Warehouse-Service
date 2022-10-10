package address

import (
	"github.com/jmoiron/sqlx"
)

var (
	queryInsertAddress = `INSERT INTO WAREHOUSE.ADDRESS (BUILDING, STREET, CITY, DISTRICT, REGION, POSTAL_CODE)
					  	  VALUES (:1, :2, :3, :4, :5, :6) RETURNING ID INTO :7`
	querySelectAddress = `SELECT * FROM WAREHOUSE.ADDRESS WHERE ID=:1`
	queryDeleteAddress = `DELETE FROM WAREHOUSE.ADDRESS WHERE ID = :1`
	queryGetAddressId  = `SELECT ID FROM WAREHOUSE.ADDRESS WHERE BUILDING = :1 AND
																 STREET = :2 AND
																 CITY = :3 AND
																 DISTRICT = :4 AND
																 REGION = :5 AND
																 POSTAL_CODE = :6`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(address *Address) (int64, error) {
	id, _ := r.GetId(address)

	if id == -1 {
		result, err := r.db.Queryx(queryInsertAddress, address.Building, address.Street, address.City,
			address.District, address.Region, address.PostalCode, &address.Id)

		if err != nil {
			return -1, err
		}

		defer result.Close()

		return address.Id, nil
	}

	return id, nil
}

func (r *Repository) Select(id int64) (*Address, error) {
	result, err := r.db.Queryx(querySelectAddress, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	address := new(Address)

	result.Next()
	err = result.StructScan(address)

	if err != nil {
		return nil, err
	}

	return address, nil
}

func (r *Repository) GetId(address *Address) (int64, error) {
	result, err := r.db.Queryx(queryGetAddressId, address.Building, address.Street, address.City,
		address.District, address.Region, address.PostalCode)

	if err != nil {
		return -1, err
	}

	defer result.Close()

	result.Next()
	err = result.Scan(&address.Id)

	if err != nil {
		return -1, err
	}

	return address.Id, nil
}

func (r *Repository) Update(address *Address) (int64, error) {
	id, _ := r.GetId(address)

	if id == -1 {
		return r.Insert(address)
	}

	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(queryDeleteAddress, id)

	if err != nil {
		return err
	}

	return nil
}
