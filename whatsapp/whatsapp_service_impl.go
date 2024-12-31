package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/whatsapp/dto"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type ServiceImpl struct {
	whatsAppRepository Repository
	gormConnection     *gorm.DB
	viperConfig        *viper.Viper
	engTranslator      ut.Translator
	vertexClient       *config.VertexClient
}

func NewService(whatsAppRepository Repository,
	gormConnection *gorm.DB,
	viperConfig *viper.Viper,
	engTranslator ut.Translator,
	vertexClient *config.VertexClient) *ServiceImpl {
	return &ServiceImpl{
		whatsAppRepository: whatsAppRepository,
		gormConnection:     gormConnection,
		viperConfig:        viperConfig,
		engTranslator:      engTranslator,
		vertexClient:       vertexClient,
	}
}

func (whatsAppService *ServiceImpl) HandleVerifyTokenWebhook(ginContext *gin.Context) {
	mode := ginContext.Query("hub.mode")
	token := ginContext.Query("hub.verify_token")
	challenge := ginContext.Query("hub.challenge")

	if mode == "subscribe" && token == whatsAppService.viperConfig.GetString("META_WEBHOOK_VERIFY_TOKEN") {
		ginContext.String(http.StatusOK, challenge) // Kirim kembali challenge jika valid
		return
	}
	ginContext.String(http.StatusForbidden, "Forbidden")
}

func (whatsAppService *ServiceImpl) HandleMessageWebhook(ginContext *gin.Context, payloadMessageDto *dto.PayloadMessageDto) {
	err := whatsAppService.gormConnection.Transaction(func(gormTransaction *gorm.DB) error {
		allWhatsAppMessage := mapper.MapPayloadIntoWhatsAppMessageModel(payloadMessageDto)

		for _, message := range allWhatsAppMessage {
			if message.SenderPhoneNumber != "" {
				whatsAppService.SendMessage(message.SenderPhoneNumber, "Permintaan anda sedang diproses")
				content, err := whatsAppService.vertexClient.GenerateContent(
					fmt.Sprintf(
						`
I have text like this %s
I want you to parse the text into JSON format with the following schema:
{
"title": "string (required, 3-100 characters)",
"content": "string (optional, max 255 characters)",
"priority": "string (required, one of: Low, Medium, High, default: false)",
"due_date": "string (format: YYYY-MM-DD HH:mm)",
"is_pinned": "boolean (default: false)",
"is_archived": "boolean (default: false)"
}
Please pay close attention to the following rules
1. If the text contains clear and concise information, use it as the title
2. If the text is not long enough to be a title or there is not enough information for a title, treat the title as "No Title" or leave it blank if desired
3. If there is a detailed description after the title, use it as the content
4. If there is no explicit content, leave it blank or set the value default like "No content provided"
5. If the text contains keywords that indicate urgency, set priority to "High".
6. If there are no words that indicate urgency, set priority to "Medium".
7. If there are words like "urgent" or "soon", set priority to "High".
8. If there is a date or time mentioned in the text, extract the date and determine the due_date, but make sure to preface it with relevant keywords like deadline, collected, requested.
9. If a time is mentioned (for example, '7 o'clock'), compare it to the current time (from the time the request was sent). If the time mentioned is past, add 1 day to the date when the request was sent to determine the due_date
10. If the text contains words like "tomorrow", "next week", or relative dates, specify the appropriate date.
11. If there is a word indicating that the item is important, set is_pinned to true
12. If there is no indication of the importance of the note, set is_pinned to false
13. If the text contains a word indicating that the note is done or does not need to be prioritized, set is_archived to true
14. If there is no indication of archiving, set is_archived to false
15. Extract the time first. Then, determine whether the date needs to be shifted based on whether the time has passed or not
16. If the current time is 3:00 PM, and the text says '7 o'clock': the due_date should be tomorrow (tomorrow's date) in the format (YYYY-MM-DD HH:mm).
17. If the current time is 5:00 AM, and the text says '7 o'clock': the due_date is the date the request was sent in the format (YYYY-MM-DD HH:mm).
18. Check whether the specified time has passed. If yes, add 1 day.
ONLY RETURN JSON FORMAT, DO NOT RETURN ANYTHING ELSE
`, message.Text))

				rb, err := json.MarshalIndent(content, "", "  ")
				fmt.Println(content.Candidates[0].Content.Parts[0])
				if err != nil {
					fmt.Println("json.MarshalIndent: %w", err)
				}
				fmt.Println(string(rb))
			}

			err := gormTransaction.Create(&allWhatsAppMessage).Error
			helper.CheckErrorOperation(err, exception.ParseGormError(err))
		}
		return nil
	})
	helper.CheckErrorOperation(err, exception.ParseGormError(err))
}

func (whatsAppService *ServiceImpl) HandleCreate(ginContext *gin.Context) {}

func (whatsAppService *ServiceImpl) SendMessage(targetNumber string, payloadMessage string) {
	// URL endpoint WhatsApp API
	endpointUrl := fmt.Sprintf(whatsAppService.viperConfig.GetString("META_ENDPOINT_SEND_MESSAGE"), "519218867943122")
	token := whatsAppService.viperConfig.GetString("META_GRAPH_API_TOKEN") // Ganti dengan token akses Anda

	// Payload untuk API
	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                targetNumber,
		"type":              "text",
		"text": map[string]string{
			"body": payloadMessage,
		},
	}

	// Konversi payload ke JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	// Kirim permintaan HTTP POST
	req, err := http.NewRequest("POST", endpointUrl, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	// Cek status respons
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("failed to send message: %s", string(body))
	}

}
