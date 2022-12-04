/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type UsedIngredient struct {
	Key
	Attributes UsedIngredientAttributes `json:"attributes"`
}
type UsedIngredientResponse struct {
	Data     UsedIngredient `json:"data"`
	Included Included       `json:"included"`
}

type UsedIngredientListResponse struct {
	Data     []UsedIngredient `json:"data"`
	Included Included         `json:"included"`
	Links    *Links           `json:"links"`
}

// MustUsedIngredient - returns UsedIngredient from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustUsedIngredient(key Key) *UsedIngredient {
	var usedIngredient UsedIngredient
	if c.tryFindEntry(key, &usedIngredient) {
		return &usedIngredient
	}
	return nil
}
