package requests

import (
	"encoding/json"
	"net/http"
	"warehouse-service/internal/service/helpers"
	"warehouse-service/resources"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateWarehouseIngredientRequest struct {
	Data resources.WarehouseIngredient
}

func NewCreateWarehouseIngredientRequest(r *http.Request) (CreateWarehouseIngredientRequest, error) {
	var request CreateWarehouseIngredientRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateWarehouseIngredientRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/attributes/quantity":        validation.Validate(&r.Data.Attributes.Quantity, validation.Required),
		"/data/attributes/origin":          validation.Validate(&r.Data.Attributes.Origin, validation.Required),
		"/data/attributes/price":           validation.Validate(&r.Data.Attributes.Price, validation.Required),
		"/data/attributes/expiration_date": validation.Validate(&r.Data.Attributes.ExpirationDate, validation.Required),
		"/data/attributes/ingredient":      validation.Validate(&r.Data.Relationships.Ingredient, validation.Required),
		"/data/attributes/warehouse":       validation.Validate(&r.Data.Relationships.Warehouse, validation.Required),
		"/data/attributes/delivery":        validation.Validate(&r.Data.Relationships.Delivery, validation.Required),
	}).Filter()
}
