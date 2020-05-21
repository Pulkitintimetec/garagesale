package product

import "time"

// ProductStructure is an item we sell.
type ProductStructure struct {
	Name        string    `json:"Name"`
	Cost        string    `json:"Cost"`
	Quantity    string    `json:"Quantity"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
	ID          string    `json:"id"`
	UserID      string    `json:"userid"`
}

// NewProduct is an item given by user
type NewProduct struct {
	Name     string `json:"Name"     validate:"required"`
	Cost     string `json:"Cost"     validate:"gte=0"`
	Quantity string `json:"Quantity" validate:"gte=1"`
	ID       string `json:"id" validate:"required"`
	UserID   string `json:"userid"`
}

// Sale represents one item of a transaction where some amount of a product was
// sold. Quantity is the number of units sold and Paid is the total price paid.
// Note that due to haggling the Paid value might not equal Quantity sold *
// Product cost.
type Sale struct {
	ID          string    `json:"ID"`
	ProductID   string    ` json:"ProductID"`
	Quantity    string    ` json:"Quantity"`
	Paid        string    `json:"Paid"`
	DateCreated time.Time `json:"date_created"`
}

// NewSale is what we require from clients for recording new transactions.
type NewSale struct {
	Quantity string `json:"Quantity" validate:"gte=0"`
	Paid     string `json:"Paid"  validate:"gte=0"`
}

// UpdateProductStructure used as a structure for update Data of Product
type UpdateProductStructure struct {
	Name     *string `json:"Name"`
	Cost     *string `json:"Cost" `
	Quantity *string `json:"Quantity" validate:"gte=1"`
}
