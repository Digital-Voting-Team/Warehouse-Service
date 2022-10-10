package delivery

import (
	"github.com/jmoiron/sqlx"
)

var (
	queryInsertDelivery = `INSERT INTO WAREHOUSE.DELIVERY (SOURCE_ID, DESTINATION_ID, DELIVERY_PRICE, DELIVERY_DATE)
					  	  VALUES (:1, :2, :3, :4) RETURNING ID INTO :5`
	querySelectDelivery = `SELECT * FROM WAREHOUSE.DELIVERY WHERE ID=:1`
	queryDeleteDelivery = `DELETE FROM WAREHOUSE.DELIVERY WHERE ID = :1`
	queryGetDeliveryId  = `SELECT ID FROM WAREHOUSE.DELIVERY WHERE SOURCE_ID = :1 AND
																   DESTINATION_ID = :2 AND
																   DELIVERY_PRICE = :3 AND
																   DELIVERY_DATE = :4`
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Insert(delivery *Delivery) (int64, error) {
	id, _ := r.GetId(delivery)

	if id == -1 {
		result, err := r.db.Queryx(queryInsertDelivery, delivery.SourceId, delivery.DestinationId,
			delivery.DeliveryPrice, delivery.DeliveryDate, &delivery.Id)

		if err != nil {
			return -1, err
		}

		defer result.Close()

		return delivery.Id, nil
	}

	return id, nil
}

func (r *Repository) Select(id int64) (*Delivery, error) {
	result, err := r.db.Queryx(querySelectDelivery, id)

	if err != nil {
		return nil, err
	}

	defer result.Close()

	delivery := new(Delivery)

	result.Next()
	err = result.StructScan(delivery)

	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (r *Repository) GetId(delivery *Delivery) (int64, error) {
	result, err := r.db.Queryx(queryGetDeliveryId, delivery.SourceId, delivery.DestinationId,
		delivery.DeliveryPrice, delivery.DeliveryDate)

	if err != nil {
		return -1, err
	}

	defer result.Close()

	result.Next()
	err = result.Scan(&delivery.Id)

	if err != nil {
		return -1, err
	}

	return delivery.Id, nil
}

func (r *Repository) Update(delivery *Delivery) (int64, error) {
	id, _ := r.GetId(delivery)

	if id == -1 {
		return r.Insert(delivery)
	}

	return id, nil
}

func (r *Repository) Delete(id int64) error {
	_, err := r.db.Exec(queryDeleteDelivery, id)

	if err != nil {
		return err
	}

	return nil
}
