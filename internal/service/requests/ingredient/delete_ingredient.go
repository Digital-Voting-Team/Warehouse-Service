package requests

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type DeleteIngredientRequest struct {
	IngredientID int64 `url:"-"`
}

func NewDeleteIngredientRequest(r *http.Request) (DeleteIngredientRequest, error) {
	request := DeleteIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.IngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	return request, nil
}
