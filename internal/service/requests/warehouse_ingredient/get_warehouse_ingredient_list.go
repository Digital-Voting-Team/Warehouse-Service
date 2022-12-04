package requests

import (
	"net/http"
	"time"

	"gitlab.com/distributed_lab/kit/pgdb"
	"gitlab.com/distributed_lab/urlval"
)

type GetWarehouseIngredientListRequest struct {
	pgdb.OffsetPageParams
	FilterIngredientId   []int64     `filter:"ingredient_id"`
	FilterWarehouseId    []int64     `filter:"warehouse_id"`
	FilterQuantity       []int64     `filter:"quantity"`
	FilterOrigin         []string    `filter:"origin"`
	FilterPrice          []float64   `filter:"price"`
	FilterExpirationDate []time.Time `filter:"expiration_date"`
	FilterDeliveryId     []int64     `filter:"delivery_id"`
}

func NewGetWarehouseIngredientListRequest(r *http.Request) (GetWarehouseIngredientListRequest, error) {
	var request GetWarehouseIngredientListRequest

	err := urlval.Decode(r.URL.Query(), &request)
	if err != nil {
		return request, err
	}

	return request, nil
}
