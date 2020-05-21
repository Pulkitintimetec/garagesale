package product

// ProductStructure is an item we sell.
type ProductStructure struct {
	Name     string `json:"Name"`
	Cost     string `json:"Cost"`
	Quantity string `json:"Quantity"`
	// DateCreated time.Time `db:"date_created" json:"date_created"`
	// DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}
