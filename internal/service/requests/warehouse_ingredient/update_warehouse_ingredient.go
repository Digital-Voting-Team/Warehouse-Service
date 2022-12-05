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

type UpdateWarehouseIngredientRequest struct {
	WarehouseIngredientID int64 `url:"-" json:"-"`
	Data                  resources.WarehouseIngredient
}

func NewUpdateWarehouseIngredientRequest(r *http.Request) (UpdateWarehouseIngredientRequest, error) {
	request := UpdateWarehouseIngredientRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.WarehouseIngredientID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateWarehouseIngredientRequest) validate() error {
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
