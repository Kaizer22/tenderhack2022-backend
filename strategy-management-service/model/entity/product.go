package entity

// Product example
type Product struct {
	Id          int64     `pg:"id,pk" json:"id" example:"1" format:"int64"`
	Name        string    `pg:"product_name" json:"name" example:"Pepsi"`
	Description string    `pg:"description" json:"description" example:"Carbonated sweet drink"`
	Category    *Category `pg:"rel:has-one" json:"category,omitempty"`
	CategoryId  int64     `pg:"category_id" json:"category_id" example:"3"`
}

// ProductData example
type ProductData struct {
	Name        string `json:"name" example:"Pepsi"`
	Description string `json:"description" example:"Carbonated sweet drink"`
	CategoryId  int64  `json:"category_id" example:"3"`
}
