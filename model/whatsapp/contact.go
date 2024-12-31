package whatsapp

type Contact struct {
	Profile    Profile `json:"profile"`
	WhatsAppId string  `json:"wa_id"`
}
