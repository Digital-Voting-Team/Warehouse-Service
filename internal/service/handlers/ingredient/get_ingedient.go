package handlers

import (
	"net/http"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func GetIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewGetIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	ingredient, err := helpers.IngredientsQuery(r).FilterById(request.IngredientID).Get()
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to get ingredient from DB")
		ape.Render(w, problems.InternalError())
		return
	}
	if ingredient == nil {
		ape.Render(w, problems.NotFound())
		return
	}

	result := resources.IngredientResponse{
		Data: resources.Ingredient{
			Key: resources.NewKeyInt64(ingredient.Id, resources.INGREDIENT),
			Attributes: resources.IngredientAttributes{
				Name: ingredient.Name,
			},
		},
	}

	ape.Render(w, result)
}
