package domain

import "github.com/doug-martin/goqu/v9"

var (
	TabNameItem = goqu.T("items")
	ColId       = "id"
	ColName     = "name"
	ColPrice    = "price"
	ColType     = "type"
	ColImg      = "img"
)

type Item struct {
	Id    int32  `json:"id"`
	Name  string `json:"name"`
	Price int32  `json:"price"`
	Type  int32  `json:"type"`
	Img   string `json:"img"`
}
