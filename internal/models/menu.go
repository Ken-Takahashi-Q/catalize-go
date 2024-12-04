package models

type Menu struct {
	MenuItemOld []MenuItemOld `bson:"MenuOld" json:"menu_old"`
	MenuString  string        `bson:"Menu" json:"menu_string"`
	Menu        []MenuItem    `json:"menu"`
}

type MenuItemOld struct {
	MenuID     int     `bson:"MenuID" json:"menu_id"`
	MenuNameTH string  `bson:"MenuNameTH" json:"menu_name_th"`
	MenuNameEN string  `bson:"MenuNameEN" json:"menu_name_en"`
	Price      float64 `bson:"Price" json:"price"`
}

type MenuItem struct {
	MenuID     int     `bson:"MenuID" json:"menu_id"`
	MenuNameTH string  `bson:"MenuNameTH" json:"menu_name_th"`
	MenuNameEN string  `bson:"MenuNameEN" json:"menu_name_en"`
	Price      float64 `bson:"Price" json:"price"`
	Category   int     `bson:"Category" json:"category"`
	Status     int     `bson:"Status" json:"status"`
	Image      string  `bson:"Image" json:"image"`
}

type MenuCategory struct {
	Category   int    `bson:"Category" json:"category"`
	CategoryTH string `bson:"CategoryTH" json:"category_th"`
	CategoryEN string `bson:"CategoryEN" json:"category_en"`
}
