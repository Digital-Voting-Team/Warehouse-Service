package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateUsedIngredientRequest struct {
	Data resources.UsedIngredient
}

func NewCreateUsedIngredientRequest(r *http.Request) (CreateUsedIngredientRequest, error) {
	var request CreateUsedIngredientRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateUsedIngredientRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/name":          validation.Validate(&r.Data.Attributes.Name, validation.Required),
		"/data/attributes/quantity":      validation.Validate(&r.Data.Attributes.Quantity, validation.Required),
		"/data/attributes/origin":        validation.Validate(&r.Data.Attributes.Origin, validation.Required),
		"/data/attributes/price":         validation.Validate(&r.Data.Attributes.Price, validation.Required),
		"/data/attributes/deletion_date": validation.Validate(&r.Data.Attributes.DeletionDate, validation.Required),
		"/data/attributes/reason":        validation.Validate(&r.Data.Attributes.Reason, validation.Required),
	}).Filter()
}
