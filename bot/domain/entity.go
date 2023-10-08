package domain

type Stock struct {
	Code  string
	Value string
}

func NewStock(code, value string) Stock {
	return Stock{
		Code:  code,
		Value: value,
	}
}
