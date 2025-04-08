package v20230530

// Residual 残差图
type Residual struct {
	Vars          []*ResidualVar `json:"Vars"`
	AvailableXvar []string       `json:"AvailableXvar"`
}

// ResidualVar 残差图的值
type ResidualVar struct {
	Values []float64 `json:"Values"`
	Name   string    `json:"Name"`
}
