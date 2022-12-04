/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type WarehouseIngredient struct {
	Key
	Attributes    WarehouseIngredientAttributes    `json:"attributes"`
	Relationships WarehouseIngredientRelationships `json:"relationships"`
}
type WarehouseIngredientResponse struct {
	Data     WarehouseIngredient `json:"data"`
	Included Included            `json:"included"`
}

type WarehouseIngredientListResponse struct {
	Data     []WarehouseIngredient `json:"data"`
	Included Included              `json:"included"`
	Links    *Links                `json:"links"`
}

// MustWarehouseIngredient - returns WarehouseIngredient from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustWarehouseIngredient(key Key) *WarehouseIngredient {
	var warehouseIngredient WarehouseIngredient
	if c.tryFindEntry(key, &warehouseIngredient) {
		return &warehouseIngredient
	}
	return nil
}
