package whatsapp

type Message struct {
	From      string `json:"from"`
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Text      *Text  `json:"text"`
	Media     *Image `json:"image"`
	Type      string `json:"type"`
}
