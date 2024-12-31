package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
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
}

func NewService(whatsAppRepository Repository,
	gormConnection *gorm.DB,
	viperConfig *viper.Viper,
	engTranslator ut.Translator) *ServiceImpl {
	return &ServiceImpl{
		whatsAppRepository: whatsAppRepository,
		gormConnection:     gormConnection,
		viperConfig:        viperConfig,
		engTranslator:      engTranslator,
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
	whatsAppService.gormConnection.Transaction(func(gormTransaction *gorm.DB) error {
		allWhatsAppMessage := mapper.MapPayloadIntoWhatsAppMessageModel(payloadMessageDto)
		for _, message := range allWhatsAppMessage {
			if message.SenderPhoneNumber != "" {
				whatsAppService.SendMessage(message.SenderPhoneNumber, "Permintaan anda sedang diproses")
			}

			gormTransaction.Create(&allWhatsAppMessage)
		}
		return nil
	})
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
