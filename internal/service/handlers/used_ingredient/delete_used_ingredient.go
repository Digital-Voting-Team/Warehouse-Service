package handlers

import (
	"net/http"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/used_ingredient"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func DeleteUsedIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewDeleteUsedIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	usedIngredient, err := helpers.UsedIngredientsQuery(r).FilterById(request.UsedIngredientID).Get()
	if usedIngredient == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	err = helpers.UsedIngredientsQuery(r).Delete(request.UsedIngredientID)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to used ingredient")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	w.WriteHeader(http.StatusOK)
}
