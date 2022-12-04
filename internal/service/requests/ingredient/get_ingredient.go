package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetIngredientRequest struct {
	IngredientID int64 `url:"-"`
}

func NewGetIngredientRequest(r *http.Request) (GetIngredientRequest, error) {
	request := GetIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.IngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
