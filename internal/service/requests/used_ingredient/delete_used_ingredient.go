package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteUsedIngredientRequest struct {
	UsedIngredientID int64 `url:"-"`
}

func NewDeleteUsedIngredientRequest(r *http.Request) (DeleteUsedIngredientRequest, error) {
	request := DeleteUsedIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UsedIngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
