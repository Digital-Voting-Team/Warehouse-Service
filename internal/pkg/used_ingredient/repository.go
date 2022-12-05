package ingredient

import (
	"database/sql"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/common"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"
	"time"
)

const usedIngredientsTableName = "warehouse.used_ingredient"

var (
	queryInsertUsedIngredient = `INSERT INTO WAREHOUSE.USED_INGREDIENT (QUANTITY, ORIGIN, PRICE, DELETION_DATE, REASON, NAME)
					  	  VALUES (:1, :2, :3, :4, :5, :6) RETURNING ID INTO :7`
)

type query struct {
	db        *sqlx.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func NewQuery(db *sqlx.DB) Query {
	return &query{
		db:        db,
		sql:       sq.Select("*").From(usedIngredientsTableName).PlaceholderFormat(sq.Colon),
		sqlUpdate: sq.Update(usedIngredientsTableName).PlaceholderFormat(sq.Colon),
	}
}

func (q *query) New() Query {
	return NewQuery(q.db)
}

func (q *query) Get() (*UsedIngredient, error) {
	var result UsedIngredient

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]UsedIngredient, error) {
	var result []UsedIngredient
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(usedIngredient UsedIngredient) (UsedIngredient, error) {

	result, err := q.db.Queryx(queryInsertUsedIngredient, usedIngredient.Quantity,
		usedIngredient.Origin, usedIngredient.Price, usedIngredient.DeletionDate,
		usedIngredient.Reason, usedIngredient.DeletionDate, &usedIngredient.Id)

	if err != nil {
		return UsedIngredient{}, err
	}

	defer result.Close()

	return usedIngredient, nil

}

func (q *query) Update(usedIngredient UsedIngredient) (UsedIngredient, error) {
	var result *UsedIngredient
	clauses := structs.Map(usedIngredient)
	clauses["name"] = usedIngredient.Name
	clauses["quantity"] = usedIngredient.Quantity
	clauses["origin"] = usedIngredient.Origin
	clauses["price"] = usedIngredient.Price
	clauses["deletion_date"] = usedIngredient.DeletionDate
	clauses["reason"] = usedIngredient.Reason

	sqlString, args, _ := q.sqlUpdate.SetMap(clauses).ToSql()

	_, err := q.db.Exec(sqlString, args...)

	result, err = q.Get()

	return *result, err
}

func (q *query) Delete(id int64) error {
	stmt := sq.Delete(usedIngredientsTableName).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Colon)
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

func (q *query) FilterByName(names ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"name": names})
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

func (q *query) FilterByDeletionDate(dates ...time.Time) Query {
	q.sql = q.sql.Where(sq.Eq{"deletion_date": dates})
	return q
}

func (q *query) FilterByReason(reasons ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"reason": reasons})
	return q
}
