package address

import (
	"database/sql"
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/common"
	"github.com/jmoiron/sqlx"
	"gitlab.com/distributed_lab/kit/pgdb"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
)

const addressesTableName = "warehouse.address"

var (
	queryInsertAddress = `INSERT INTO WAREHOUSE.ADDRESS (BUILDING, STREET, CITY, DISTRICT, REGION, POSTAL_CODE)
					  	  VALUES (:1, :2, :3, :4, :5, :6) RETURNING ID INTO :7` // TODO : check "RETURNING * INTO :7"
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

func (q *query) Get() (*Address, error) {
	var result Address

	sqlString, args, _ := q.sql.ToSql()
	err := q.db.Get(&result, sqlString, args...)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return &result, err
}

func (q *query) Select() ([]Address, error) {
	var result []Address
	sqlString, args, _ := q.sql.ToSql()

	err := q.db.Select(&result, sqlString, args...)
	return result, err
}

func (q *query) Insert(address Address) (Address, error) {
	result, err := q.db.Queryx(queryInsertAddress, address.Building, address.Street, address.City,
		address.District, address.Region, address.PostalCode, &address.Id)

	if err != nil {
		return Address{}, err
	}

	defer result.Close()

	return address, nil
}

func (q *query) Update(address Address) (Address, error) {
	var result *Address
	clauses := structs.Map(address)
	clauses["building"] = address.Building
	clauses["street"] = address.Street
	clauses["city"] = address.City
	clauses["district"] = address.District
	clauses["region"] = address.Region
	clauses["postal_code"] = address.PostalCode

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

func (q *query) FilterByBuilding(numbers ...int64) Query {
	q.sql = q.sql.Where(sq.Eq{"building": numbers})
	return q
}

func (q *query) FilterByStreet(streets ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"street": streets})
	return q
}

func (q *query) FilterByCity(cities ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"city": cities})
	return q
}

func (q *query) FilterByDistrict(districts ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"district": districts})
	return q
}

func (q *query) FilterByRegion(regions ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"region": regions})
	return q
}

func (q *query) FilterByPostalCode(codes ...string) Query {
	q.sql = q.sql.Where(sq.Eq{"postal_code": codes})
	return q
}
