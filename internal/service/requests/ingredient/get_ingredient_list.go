package requests

import (
	"net/http"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetIngredientListRequest struct {
	pgdb.OffsetPageParams
	FilterName []string `filter:"name"`
}

func NewGetIngredientListRequest(r *http.Request) (GetIngredientListRequest, error) {
	var request GetIngredientListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
