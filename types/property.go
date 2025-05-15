package types

type Property struct {
	ID         string `json:"id" firestore:"omit"`
	Number     string `json:"number"`
	StreetName string `json:"streetName"`
	Town       string `json:"town"`
	Postcode   string `json:"postcode"`
}
