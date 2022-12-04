package requests

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetDeliveryListRequest struct {
	pgdb.OffsetPageParams
	FilterSourceId      []int64     `filter:"source_id"`
	FilterDestinationId []int64     `filter:"destination_id"`
	FilterPrice         []float64   `filter:"price"`
	FilterDate          []time.Time `filter:"date"`
}

func NewGetDeliveryListRequest(r *http.Request) (GetDeliveryListRequest, error) {
	var request GetDeliveryListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
