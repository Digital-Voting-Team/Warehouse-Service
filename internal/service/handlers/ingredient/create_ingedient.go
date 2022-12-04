package handlers

import (
	"net/http"
	"warehouse-service/internal/pkg/ingredient"
	"warehouse-service/internal/service/helpers"
	requests "warehouse-service/internal/service/requests/ingredient"
	"warehouse-service/resources"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateIngredient(w http.ResponseWriter, r *http.Request) {
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//if *accessLevel < resources.Manager {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	request, err := requests.NewCreateIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	var resultIngredient ingredient.Ingredient

	currentIngredient := ingredient.Ingredient{
		Name: request.Data.Attributes.Name,
	}

	resultIngredient, err = helpers.IngredientsQuery(r).Insert(currentIngredient)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to create ingredient")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	result := resources.IngredientResponse{
		Data: resources.Ingredient{
			Key: resources.NewKeyInt64(resultIngredient.Id, resources.INGREDIENT),
			Attributes: resources.IngredientAttributes{
				Name: resultIngredient.Name,
			},
		},
	}
	ape.Render(w, result)
}
