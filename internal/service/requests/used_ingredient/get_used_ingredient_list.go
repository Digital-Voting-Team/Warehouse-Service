package requests

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetUsedIngredientListRequest struct {
	pgdb.OffsetPageParams
	FilterName         []string    `filter:"name"`
	FilterQuantity     []int64     `filter:"quantity"`
	FilterOrigin       []string    `filter:"origin"`
	FilterPrice        []float64   `filter:"price"`
	FilterDeletionDate []time.Time `filter:"deletion_date"`
	FilterReason       []string    `filter:"reason"`
}

func NewGetUsedIngredientListRequest(r *http.Request) (GetUsedIngredientListRequest, error) {
	var request GetUsedIngredientListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
