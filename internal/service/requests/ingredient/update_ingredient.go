package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateIngredientRequest struct {
	IngredientID int64 `url:"-" json:"-"`
	Data         resources.Ingredient
}

func NewUpdateIngredientRequest(r *http.Request) (UpdateIngredientRequest, error) {
	request := UpdateIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.IngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateIngredientRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/name": validation.Validate(&r.Data.Attributes.Name, validation.Required,
			validation.By(helpers.IsInteger)),
	}).Filter()
}
