package models

type Menu struct {
	MenuItem []MenuItem `bson:"menu" json:"menu"`
}

type MenuItem struct {
	MenuID     int     `bson:"MenuID" json:"menu_id"`
	MenuNameTH string  `bson:"MenuNameTH" json:"menu_name_th"`
	MenuNameEN string  `bson:"MenuNameEN" json:"menu_name_en"`
	Price      float64 `bson:"Price" json:"price"`
}
