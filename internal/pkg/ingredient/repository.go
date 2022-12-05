package ingredient

import (
	"database/sql"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/common"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const ingredientsTableName = "warehouse.ingredient"

var (
	queryInsertIngredient = `INSERT INTO WAREHOUSE.INGREDIENT (NAME)
					  	  VALUES (:1) RETURNING ID INTO :2`
)

type query struct {
	db        *sqlx.DB
	sql       sq.SelectBuilder
	sqlUpdate sq.UpdateBuilder
}

func NewQuery(db *sqlx.DB) Query {
	return &query{
		db:        db,
		sql:       sq.Select("*").From(ingredientsTableName).PlaceholderFormat(sq.Colon),
		sqlUpdate: sq.Update(ingredientsTableName).PlaceholderFormat(sq.Colon),
	}
}

func (q *query) New() Query {
	return NewQuery(q.db)
}

func (q *query) Get() (*Ingredient, error) {
	var result Ingredient

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]Ingredient, error) {
	var result []Ingredient
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(ingredient Ingredient) (Ingredient, error) {
	result, err := q.db.Queryx(queryInsertIngredient, ingredient.Name, &ingredient.Id)

	if err != nil {
		return Ingredient{}, err
	}

	defer result.Close()

	return ingredient, nil
}

func (q *query) Update(ingredient Ingredient) (Ingredient, error) {
	var result *Ingredient
	clauses := structs.Map(ingredient)
	clauses["name"] = ingredient.Name

	sqlString, args, _ := q.sqlUpdate.SetMap(clauses).ToSql()

	_, err := q.db.Exec(sqlString, args...)

	result, err = q.Get()

	return *result, err
}

func (q *query) Delete(id int64) error {
	stmt := sq.Delete(ingredientsTableName).Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Colon)
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
