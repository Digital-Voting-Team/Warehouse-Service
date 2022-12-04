/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Ingredient struct {
	Key
	Attributes IngredientAttributes `json:"attributes"`
}
type IngredientResponse struct {
	Data     Ingredient `json:"data"`
	Included Included   `json:"included"`
}

type IngredientListResponse struct {
	Data     []Ingredient `json:"data"`
	Included Included     `json:"included"`
	Links    *Links       `json:"links"`
}

// MustIngredient - returns Ingredient from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustIngredient(key Key) *Ingredient {
	var ingredient Ingredient
	if c.tryFindEntry(key, &ingredient) {
		return &ingredient
	}
	return nil
}
