package common

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"gitlab.com/distributed_lab/kit/pgdb"
)

func ApplyPageParams(p *pgdb.OffsetPageParams, sql sq.SelectBuilder, cols ...string) sq.SelectBuilder {
	if p.Limit == 0 {
		p.Limit = 15
	}

	if p.Order == "" {
		p.Order = pgdb.OrderTypeDesc
	}

	offset := p.Limit * p.PageNumber

	var orderString string

	switch p.Order {
	case pgdb.OrderTypeAsc:
		orderString = "asc"
	case pgdb.OrderTypeDesc:
		orderString = "desc"
	default:
		panic(fmt.Errorf("unexpected order type: %v", p.Order))
	}

	for _, col := range cols {
		sql = sql.OrderBy(fmt.Sprintf("%s %s", col, orderString))
	}

	return sql.Suffix(fmt.Sprintf("OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", offset, p.Limit))
}
