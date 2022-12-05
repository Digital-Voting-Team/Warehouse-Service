package handlers

import (
	"github.com/Digital-Voting-Team/warehouse-service/internal/pkg/ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/internal/service/helpers"
	requests "github.com/Digital-Voting-Team/warehouse-service/internal/service/requests/ingredient"
	"github.com/Digital-Voting-Team/warehouse-service/resources"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func UpdateIngredient(w http.ResponseWriter, r *http.Request) {
	request, err := requests.NewUpdateIngredientRequest(r)
	if err != nil {
		helpers.Log(r).WithError(err).Info("wrong request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	currentIngredient, err := helpers.IngredientsQuery(r).FilterById(request.IngredientID).Get()
	if currentIngredient == nil {
		helpers.Log(r).WithError(err).Info("did not found ingredient to update")
		ape.Render(w, problems.NotFound())
		return
	}

	//userId := r.Context().Value("userId").(int64)
	//accessLevel := r.Context().Value("accessLevel").(*resources.AccessLevel)
	//_, _, ingredientId, err := helpers.GetIdsForGivenUser(r, userId)
	//if err != nil {
	//	helpers.Log(r).WithError(err).Info("wrong relations")
	//	ape.RenderErr(w, problems.InternalError())
	//	return
	//}
	//if *accessLevel != resources.Admin && ingredientId != ingredient.ID {
	//	helpers.Log(r).Info("insufficient user permissions")
	//	ape.RenderErr(w, problems.Forbidden())
	//	return
	//}

	newIngredient := ingredient.Ingredient{
		Name: request.Data.Attributes.Name,
	}

	var resultIngredient ingredient.Ingredient
	resultIngredient, err = helpers.IngredientsQuery(r).FilterById(currentIngredient.Id).Update(newIngredient)
	if err != nil {
		helpers.Log(r).WithError(err).Error("failed to update ingredient")
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
