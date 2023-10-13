package domain

type Masterlist struct {
	PsgcCode string `csv:"10-digit PSGC"       json:"psgc_code"`
	Name     string `csv:"Name"                json:"name"`
	Code     string `csv:"Correspondence Code" json:"-"`
	Level    string `csv:"Geographic Level"    json:"-"`
} //@name Masterlist
//? comment above is for renaming stuct
