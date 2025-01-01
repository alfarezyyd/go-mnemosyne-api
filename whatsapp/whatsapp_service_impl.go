package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/note"
	noteDto "go-mnemosyne-api/note/dto"
	userDto "go-mnemosyne-api/user/dto"
	"go-mnemosyne-api/whatsapp/dto"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

type ServiceImpl struct {
	whatsAppRepository Repository
	gormConnection     *gorm.DB
	viperConfig        *viper.Viper
	engTranslator      ut.Translator
	vertexClient       *config.VertexClient
	noteService        note.Service
	validationInstance *validator.Validate
}

type Content struct {
	Parts []string `json:Parts`
	Role  string   `json:Role`
}
type Candidates struct {
	Content *Content `json:Content`
}
type ContentResponse struct {
	Candidates *[]Candidates `json:Candidates`
}

func NewService(whatsAppRepository Repository,
	gormConnection *gorm.DB,
	viperConfig *viper.Viper,
	engTranslator ut.Translator,
	vertexClient *config.VertexClient,
	noteService note.Service) *ServiceImpl {
	return &ServiceImpl{
		whatsAppRepository: whatsAppRepository,
		gormConnection:     gormConnection,
		viperConfig:        viperConfig,
		engTranslator:      engTranslator,
		vertexClient:       vertexClient,
		noteService:        noteService,
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
	err := whatsAppService.validationInstance.Struct(&payloadMessageDto)
	helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	err = whatsAppService.gormConnection.Transaction(func(gormTransaction *gorm.DB) error {
		allWhatsAppMessage := mapper.MapPayloadIntoWhatsAppMessageModel(payloadMessageDto)

		for _, message := range allWhatsAppMessage {
			if message.SenderPhoneNumber != "" && message.Text != "" {
				whatsAppService.SendMessage(message.SenderPhoneNumber, "Permintaan anda sedang diproses")
				content, err :=
					whatsAppService.vertexClient.GenerateContent(
						fmt.Sprintf(
							`
Saya memiliki teks seperti ini %s
Saya ingin Anda mengurai teks ke dalam format JSON dengan skema berikut:
{
"title": "string (diperlukan, 3-100 karakter)",
"content": "string (opsional, maks 255 karakter)",
"priority": "string (diperlukan, salah satu dari: High, Low, Medium, default: Low)",
"due_date": "string (format: YYYY-MM-DD HH:mm)",
"is_pinned": "boolean (default: false)",
"is_archived": "boolean (default: false)"
}
Harap perhatikan aturan berikut
1. Jika teks berisi informasi yang jelas dan ringkas, gunakan sebagai judul
2. Jika teks tidak cukup panjang untuk menjadi judul atau tidak ada cukup informasi untuk judul, perlakukan judul sebagai "Tanpa Judul" atau kosongkan jika diinginkan
3. Jika ada deskripsi terperinci setelah judul, gunakan sebagai konten
4. Jika tidak ada konten yang eksplisit, biarkan kosong atau tetapkan nilai default seperti "Tidak ada konten yang disediakan"
5. Jika teks berisi kata kunci yang menunjukkan urgensi, tetapkan prioritas ke "Tinggi".
6. Jika tidak ada kata yang menunjukkan urgensi, tetapkan prioritas ke "Sedang".
7. Jika ada kata seperti "mendesak" atau "segera", tetapkan prioritas ke "Tinggi".
8. Jika ada tanggal atau waktu yang disebutkan dalam teks, ekstrak tanggal tersebut dan tentukan due_date, tetapi pastikan untuk mengawalinya dengan kata kunci yang relevan seperti deadline, collected, asked.
9. Jika disebutkan waktu (misalnya, jam 7'), bandingkan dengan waktu saat ini. Jika waktu yang disebutkan sudah lewat, tambahkan 1 hari ke waktu saat ini
11. Jika ada kata yang menunjukkan bahwa item tersebut penting, tetapkan is_pinned menjadi true
12. Jika tidak ada indikasi pentingnya catatan tersebut, tetapkan is_pinned menjadi false
13. Jika teks berisi kata yang menunjukkan bahwa catatan tersebut sudah selesai atau tidak perlu diprioritaskan, tetapkan is_archived menjadi true
14. Jika tidak ada indikasi pengarsipan, tetapkan is_archived menjadi false
15. Ekstrak waktu terlebih dahulu. Kemudian, tentukan apakah tanggal perlu digeser berdasarkan apakah waktu telah berlalu atau belum
16. Periksa apakah waktu yang ditentukan telah lewat. Jika ya, tambahkan 1 hari.
17. Waktu saat ini adalah %s GMT+7 dalam format 24 jam
18. Jika tidak ditemukan informasi tanggal dan waktu, due_date dapat kosong
HANYA KEMBALIKAN FORMAT JSON, JANGAN KEMBALIKAN YANG LAIN
`, message.Text, (time.Now()).Format("2006-01-02 15.04")))
				fmt.Println(err)
				marshalResponse, _ := json.MarshalIndent(content, "", "  ")
				var generateResponse ContentResponse
				if err := json.Unmarshal(marshalResponse, &generateResponse); err != nil {
					log.Fatal(err)
				}
				var allParsedNote []noteDto.CreateNoteDto
				for _, cad := range *generateResponse.Candidates {
					if cad.Content != nil {
						var parsedNote noteDto.CreateNoteDto
						for _, part := range cad.Content.Parts {
							err := json.Unmarshal([]byte(part), &parsedNote)
							helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
						}
						parsedNote.CategoryId = 1
						allParsedNote = append(allParsedNote, parsedNote)
					}
				}
				var userModel model.User
				fmt.Println("CP11")
				err = gormTransaction.Where("phone_number = ?", message.SenderPhoneNumber).First(&userModel).Error
				helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
				userJwtClaim := userDto.JwtClaimDto{
					Email:       &userModel.Email,
					PhoneNumber: &userModel.PhoneNumber.String,
				}
				ginContext.Set("claims", &userJwtClaim)
				whatsAppService.noteService.HandleCreate(ginContext, &allParsedNote[0])
				helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
			}
		}
		whatsAppService.SendMessage(allWhatsAppMessage[0].SenderPhoneNumber, "Catatan berhasil ditambahkan")
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
