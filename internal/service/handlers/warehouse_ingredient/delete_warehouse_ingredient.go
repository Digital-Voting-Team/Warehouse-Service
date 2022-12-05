package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/warehouse_ingredient"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteWarehouseIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewDeleteWarehouseIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	warehouseIngredient, err := helpers.WarehouseIngredientsQuery(r).FilterById(request.WarehouseIngredientID).Get()
	if warehouseIngredient == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	err = helpers.WarehouseIngredientsQuery(r).Delete(request.WarehouseIngredientID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to delete warehouseIngredient")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
