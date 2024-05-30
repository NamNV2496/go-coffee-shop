package domain

import "github.com/doug-martin/goqu/v9"

var (
	TabNameItem = goqu.T("items")
)

// type Item struct {
// 	Id    int32  `json:"id"`
// 	Name  string `json:"name"`
// 	Price int32  `json:"price"`
// 	Type  int32  `json:"type"`
// 	Image string `json:"image"`
// }
