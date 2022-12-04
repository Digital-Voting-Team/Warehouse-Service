package requests

import (
	"github.com/go-chi/chi"
	"github.com/spf13/cast"

	"net/http"

	"gitlab.com/distributed_lab/urlval"
)

type GetUsedIngredientRequest struct {
	UsedIngredientID int64 `url:"-"`
}

func NewGetUsedIngredientRequest(r *http.Request) (GetUsedIngredientRequest, error) {
	request := GetUsedIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UsedIngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
