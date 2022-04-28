package entity

//Category example
type Category struct {
	Id          int64  `pg:"id, pk" json:"id" example:"42"`
	Name        string `pg:"category_name" json:"name" example:"Furniture"`
	Description string `pg:"description" json:"description" example:"Sofa, table, shelf e.g."`
}

type CategoryData struct {
	Name        string `json:"name" example:"Furniture"`
	Description string `json:"description" example:"Sofa, table, shelf e.g."`
}
