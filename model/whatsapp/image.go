package whatsapp

type Image struct {
	Caption  string `json:"caption"`
	MimeType string `json:"mime_type"`
	SHA256   string `json:"sha256"`
	ID       string `json:"id"`
}
