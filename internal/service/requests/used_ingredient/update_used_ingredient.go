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

type UpdateUsedIngredientRequest struct {
	UsedIngredientID int64 `url:"-" json:"-"`
	Data             resources.UsedIngredient
}

func NewUpdateUsedIngredientRequest(r *http.Request) (UpdateUsedIngredientRequest, error) {
	request := UpdateUsedIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.UsedIngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateUsedIngredientRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/name":          validation.Validate(&r.Data.Attributes.Name, validation.Required),
		"/data/attributes/quantity":      validation.Validate(&r.Data.Attributes.Quantity, validation.Required),
		"/data/attributes/origin":        validation.Validate(&r.Data.Attributes.Origin, validation.Required),
		"/data/attributes/price":         validation.Validate(&r.Data.Attributes.Price, validation.Required),
		"/data/attributes/deletion_date": validation.Validate(&r.Data.Attributes.DeletionDate, validation.Required),
		"/data/attributes/reason":        validation.Validate(&r.Data.Attributes.Reason, validation.Required),
	}).Filter()
}
