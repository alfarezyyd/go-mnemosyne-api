package whatsapp

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}
