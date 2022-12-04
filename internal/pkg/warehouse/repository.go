package warehouse

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"
	"warehouse-service/internal/pkg/common"
)

const addressesTableName = "warehouse.warehouse"

var (
	queryInsertWarehouse = `INSERT INTO WAREHOUSE.WAREHOUSE (CAFE_ID, ADDRESS_ID, CAPACITY)
					  	 	VALUES (:1, :2, :3) RETURNING ID INTO :4`
)

type query struct {
	db        *sqlx.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func NewQuery(db *sqlx.DB) Query {
	return &query{
		db:        db,
		sql:       sq.Select("*").From(addressesTableName).PlaceholderFormat(sq.Colon),
		sqlUpdate: sq.Update(addressesTableName).PlaceholderFormat(sq.Colon),
	}
}

func (q *query) New() Query {
	return NewQuery(q.db)
}

func (q *query) Get() (*Warehouse, error) {
	var result Warehouse

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]Warehouse, error) {
	var result []Warehouse
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(warehouse Warehouse) (Warehouse, error) {
	result, err := q.db.Queryx(queryInsertWarehouse, warehouse.CafeId,
		warehouse.AddressId, warehouse.Capacity, &warehouse.Id)

	if err != nil {
		return Warehouse{}, err
	}

	defer result.Close()

	return warehouse, nil
}

func (q *query) Update(address Warehouse) (Warehouse, error) {
	var result *Warehouse
	clauses := structs.Map(address)
	clauses["cafe_id"] = address.CafeId
	clauses["address_id"] = address.AddressId
	clauses["capacity"] = address.Capacity

	sqlString, args, _ := q.sqlUpdate.SetMap(clauses).ToSql()

	_, err := q.db.Exec(sqlString, args...)

	result, err = q.Get()

	return *result, err
}

func (q *query) Delete(id int64) error {
	stmt := sq.Delete(addressesTableName).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Colon)
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

func (q *query) FilterByCafeId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"cafe_id": ids})
	return q
}

func (q *query) FilterByAddressId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"address_id": ids})
	return q
}

func (q *query) FilterByCapacity(capacities ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"capacity": capacities})
	return q
}
