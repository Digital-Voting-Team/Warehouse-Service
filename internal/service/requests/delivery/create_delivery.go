package requests

import (
	"encoding/json"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type CreateDeliveryRequest struct {
	Data resources.Delivery
}

func NewCreateDeliveryRequest(r *http.Request) (CreateDeliveryRequest, error) {
	var request CreateDeliveryRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *CreateDeliveryRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/relations/source": validation.Validate(&r.Data.Relationships.Source, validation.Required,
			validation.By(helpers.IsInteger)),
		"/data/relations/destination": validation.Validate(&r.Data.Relationships.Destination, validation.Required,
			validation.By(helpers.IsInteger)),
		"/data/attributes/price": validation.Validate(&r.Data.Attributes.Price, validation.Required),
		"/data/attributes/date":  validation.Validate(&r.Data.Attributes.Date, validation.Required),
	}).Filter()
}
