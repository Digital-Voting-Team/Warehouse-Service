package delivery

import (
	"database/sql"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/common"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

const deliveriesTableName = "warehouse.delivery"

var (
	queryInsertDelivery = `INSERT INTO WAREHOUSE.DELIVERY (SOURCE_ID, DESTINATION_ID, PRICE, DATE)
					  	  VALUES (:1, :2, :3, :4) RETURNING ID INTO :5`
)

type query struct {
	db        *sqlx.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func NewQuery(db *sqlx.DB) Query {
	return &query{
		db:        db,
		sql:       sq.Select("*").From(deliveriesTableName).PlaceholderFormat(sq.Colon),
		sqlUpdate: sq.Update(deliveriesTableName).PlaceholderFormat(sq.Colon),
	}
}

func (q *query) New() Query {
	return NewQuery(q.db)
}

func (q *query) Get() (*Delivery, error) {
	var result Delivery

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]Delivery, error) {
	var result []Delivery
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(delivery Delivery) (Delivery, error) {
	result, err := q.db.Queryx(queryInsertDelivery, delivery.SourceId, delivery.DestinationId,
		delivery.Price, delivery.Date, &delivery.Id)

	if err != nil {
		return Delivery{}, err
	}

	defer result.Close()

	return delivery, nil
}

func (q *query) Update(delivery Delivery) (Delivery, error) {
	var result *Delivery
	clauses := structs.Map(delivery)
	clauses["source_id"] = delivery.SourceId
	clauses["destination_id"] = delivery.DestinationId
	clauses["price"] = delivery.Price
	clauses["date"] = delivery.Date

	sqlString, args, _ := q.sqlUpdate.SetMap(clauses).ToSql()

	_, err := q.db.Exec(sqlString, args...)

	result, err = q.Get()

	return *result, err
}

func (q *query) Delete(id int64) error {
	stmt := sq.Delete(deliveriesTableName).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Colon)
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

func (q *query) FilterBySourceId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"source_id": ids})
	return q
}

func (q *query) FilterByDestinationId(ids ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"destination_id": ids})
	return q
}

func (q *query) FilterByPrice(prices ...float64) Query {
	q.sql = q.sql.Where(sq.Eq{"price": prices})
	return q
}

func (q *query) FilterByDate(dates ...time.Time) Query {
	q.sql = q.sql.Where(sq.Eq{"date": dates})
	return q
}
