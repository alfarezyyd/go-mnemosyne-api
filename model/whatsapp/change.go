package whatsapp

type Change struct {
	Value MessageValue `json:"value"`
	Field string       `json:"field"`
}
