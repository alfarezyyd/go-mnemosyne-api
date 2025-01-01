package whatsapp

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      Text   `json:"text" validate:"required"`
	Type      string `json:"type"`
}
