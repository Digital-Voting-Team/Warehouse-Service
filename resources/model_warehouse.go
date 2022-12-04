/*
 * GENERATED. Do not modify. Your changes might be overwritten!
 */

package resources

type Warehouse struct {
	Key
	Attributes    WarehouseAttributes    `json:"attributes"`
	Relationships WarehouseRelationships `json:"relationships"`
}
type WarehouseResponse struct {
	Data     Warehouse `json:"data"`
	Included Included  `json:"included"`
}

type WarehouseListResponse struct {
	Data     []Warehouse `json:"data"`
	Included Included    `json:"included"`
	Links    *Links      `json:"links"`
}

// MustWarehouse - returns Warehouse from include collection.
// if entry with specified key does not exist - returns nil
// if entry with specified key exists but type or ID mismatches - panics
func (c *Included) MustWarehouse(key Key) *Warehouse {
	var warehouse Warehouse
	if c.tryFindEntry(key, &warehouse) {
		return &warehouse
	}
	return nil
}
