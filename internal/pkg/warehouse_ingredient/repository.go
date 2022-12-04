package warehouse_ingredient

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
	"warehouse-service/internal/pkg/common"
)

const warehouseIngredientTableName = "warehouse.warehouse_ingredient"

var (
	queryInsertWarehouseIngredient = `INSERT INTO WAREHOUSE.WAREHOUSE_INGREDIENT (INGREDIENT_ID, WAREHOUSE_ID, QUANTITY, ORIGIN, PRICE, EXPIRATION_DATE, DELIVERY_ID)
					  	 			  VALUES (:1, :2, :3, :4, :5, :6, :7) RETURNING ID INTO :8`
)

type query struct {
	db        *sqlx.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func NewQuery(db *sqlx.DB) Query {
	return &query{
		db:        db,
		sql:       sq.Select("*").From(warehouseIngredientTableName).PlaceholderFormat(sq.Colon),
		sqlUpdate: sq.Update(warehouseIngredientTableName).PlaceholderFormat(sq.Colon),
	}
}

func (q *query) New() Query {
	return NewQuery(q.db)
}

func (q *query) Get() (*WarehouseIngredient, error) {
	var result WarehouseIngredient

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]WarehouseIngredient, error) {
	var result []WarehouseIngredient
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(warehouseIngredient WarehouseIngredient) (WarehouseIngredient, error) {

	result, err := q.db.Queryx(queryInsertWarehouseIngredient, warehouseIngredient.IngredientId,
		warehouseIngredient.WarehouseId, warehouseIngredient.Quantity, warehouseIngredient.Origin,
		warehouseIngredient.Price, warehouseIngredient.ExpirationDate, warehouseIngredient.DeliveryId,
		&warehouseIngredient.Id)

	if err != nil {
		return WarehouseIngredient{}, err
	}

	defer result.Close()

	return warehouseIngredient, nil

}

func (q *query) Update(warehouseIngredient WarehouseIngredient) (WarehouseIngredient, error) {
	var result *WarehouseIngredient
	clauses := structs.Map(warehouseIngredient)
	clauses["ingredient_id"] = warehouseIngredient.IngredientId
	clauses["warehouse_id"] = warehouseIngredient.WarehouseId
	clauses["quantity"] = warehouseIngredient.Quantity
	clauses["origin"] = warehouseIngredient.Origin
	clauses["price"] = warehouseIngredient.Price
	clauses["expiration_date"] = warehouseIngredient.ExpirationDate
	clauses["delivery_id"] = warehouseIngredient.DeliveryId

	sqlString, args, _ := q.sqlUpdate.SetMap(clauses).ToSql()

	_, err := q.db.Exec(sqlString, args...)

	result, err = q.Get()

	return *result, err
}

func (q *query) Delete(id int64) error {
	stmt := sq.Delete(warehouseIngredientTableName).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Colon)
	sqlString, args, _ := stmt.ToSql()
	_, err := q.db.Exec(sqlString, args...)
	return err
}

func (q *query) Page(pageParams pgdb.OffsetPageParams) Query {
	q.sql = common.ApplyPageParams(&pageParams, q.sql, "id")

	return q
}

func (q *query) FilterById(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"id": ids})
	q.sqlUpdate = q.sqlUpdate.Where(sq.Eq{"id": ids})
	return q
}

func (q *query) FilterByIngredientId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"ingredient_id": ids})
	return q
}

func (q *query) FilterByWarehouseId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"warehouse_id": ids})
	return q
}

func (q *query) FilterByQuantity(quantities ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"quantity": quantities})
	return q
}

func (q *query) FilterByOrigin(origins ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"origin": origins})
	return q
}

func (q *query) FilterByPrice(prices ...float64) Query {
	q.sql = q.sql.Where(sq.Eq{"price": prices})
	return q
}

func (q *query) FilterByExpirationDate(expirationDates ...time.Time) Query {
	q.sql = q.sql.Where(sq.Eq{"expiration_date": expirationDates})
	return q
}

func (q *query) FilterByDeliveryId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"delivery_id": ids})
	return q
}
