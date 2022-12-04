package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetWarehouseIngredientRequest struct {
	WarehouseIngredientID int64 `url:"-"`
}

func NewGetWarehouseIngredientRequest(r *http.Request) (GetWarehouseIngredientRequest, error) {
	request := GetWarehouseIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.WarehouseIngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
