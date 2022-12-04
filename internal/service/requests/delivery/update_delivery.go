package requests

import (
	"encoding/json"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"net/http"
	"warehouse-service/internal/service/helpers"
	"warehouse-service/resources"

	"github.com/go-chi/chi"
	"github.com/spf13/cast"
	"gitlab.com/distributed_lab/urlval"
)

type UpdateDeliveryRequest struct {
	DeliveryID int64 `url:"-" json:"-"`
	Data       resources.Delivery
}

func NewUpdateDeliveryRequest(r *http.Request) (UpdateDeliveryRequest, error) {
	request := UpdateDeliveryRequest{}

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	request.DeliveryID = cast.ToInt64(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, errors.Wrap(err, "failed to unmarshal")
	}

	return request, request.validate()
}

func (r *UpdateDeliveryRequest) validate() error {
	return helpers.MergeErrors(validation.Errors{
		"/data/relations/source": validation.Validate(&r.Data.Relationships.Source, validation.Required,
			validation.By(helpers.IsInteger)),
		"/data/relations/destination": validation.Validate(&r.Data.Relationships.Destination, validation.Required,
			validation.By(helpers.IsInteger)),
		"/data/attributes/price": validation.Validate(&r.Data.Attributes.Price, validation.Required),
		"/data/attributes/date":  validation.Validate(&r.Data.Attributes.Date, validation.Required),
	}).Filter()
}
