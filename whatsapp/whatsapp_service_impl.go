package whatsapp

import (
	"bytes"
	"cloud.google.com/go/vertexai/genai"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
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
	"os"
	"strings"
	"time"
)

type ServiceImpl struct {
	whatsAppRepository Repository
	gormConnection     *gorm.DB
	viperConfig        *viper.Viper
	engTranslator      ut.Translator
	vertexClient       *config.VertexClient
	noteService        note.Service
	googleCloudStorage *config.GoogleCloudStorage
}

type Content struct {
	Parts []string `json:"Parts"`
	Role  string   `json:"Role"`
}
type Candidates struct {
	Content *Content `json:"Content"`
}
type ContentResponse struct {
	Candidates *[]Candidates `json:"Candidates"`
}

func NewService(whatsAppRepository Repository,
	gormConnection *gorm.DB,
	viperConfig *viper.Viper,
	engTranslator ut.Translator,
	vertexClient *config.VertexClient,
	noteService note.Service,
	storage *config.GoogleCloudStorage) *ServiceImpl {
	return &ServiceImpl{
		whatsAppRepository: whatsAppRepository,
		gormConnection:     gormConnection,
		viperConfig:        viperConfig,
		engTranslator:      engTranslator,
		vertexClient:       vertexClient,
		noteService:        noteService,
		googleCloudStorage: storage,
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
	fmt.Println(payloadMessageDto)
	if len(payloadMessageDto.Entry[0].Changes[0].Value.Messages) == 0 {
		return
	}

	err := whatsAppService.gormConnection.Transaction(func(gormTransaction *gorm.DB) error {
		allWhatsAppMessage := mapper.MapPayloadIntoWhatsAppMessageModel(payloadMessageDto)
		for _, message := range allWhatsAppMessage {
			if message.SenderPhoneNumber != "" {
				whatsAppService.SendMessage(message.SenderPhoneNumber, "Permintaan anda sedang diproses")
				var openingPrompt string
				var content *genai.GenerateContentResponse
				var err error
				switch message.Type {
				case "text":
					fmt.Println(*(message.Text))
					if strings.HasPrefix(*(message.Text), "/") {
						switch strings.ToLower(strings.Replace(*(message.Text), "/", "", 1)) {
						case "getall":
							var allNote []model.Note
							gormTransaction.Joins("JOIN users ON users.id = notes.user_id").Where("users.phone_number = ?", message.SenderPhoneNumber).Find(&allNote)
							whatsAppService.SendMessage(message.SenderPhoneNumber, mapper.MapAllNoteIntoString(allNote))
							break
						}
						return nil
					}
					openingPrompt = fmt.Sprintf("Saya memiliki teks seperti ini %s", *(message.Text))
					content, err =
						whatsAppService.vertexClient.GenerateContent(
							fmt.Sprintf(
								`
%s
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
`, openingPrompt, (time.Now()).Format("2006-01-02 15.04")))
					break
				case "image":
					openingPrompt = fmt.Sprintf("Lakukan OCR pada gambar yang dilampirkan dan")
					mediaURL, err := whatsAppService.retrieveMediaLocation(*(message.MediaId))
					if err != nil {
						fmt.Println("Error getting media URL:", err)
					}

					publicUrl := whatsAppService.downloadMedia(*(message.MediaId), mediaURL, strings.TrimPrefix(*(message.MimeType), "image/"))
					if err != nil {
						fmt.Println("Error downloading media:", err)
					}
					content, err =
						whatsAppService.vertexClient.GenerateContentWithImage(
							publicUrl,
							fmt.Sprintf(
								`
%s
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
`, openingPrompt, (time.Now()).Format("2006-01-02 15.04")))
					break
				}

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
				fmt.Println(generateResponse)
				fmt.Println(allParsedNote[0])
				var userModel model.User
				err = gormTransaction.Where("phone_number = ?", message.SenderPhoneNumber).First(&userModel).Error
				helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
				userJwtClaim := userDto.JwtClaimDto{
					Email:       &userModel.Email,
					PhoneNumber: &userModel.PhoneNumber.String,
				}
				ginContext.Set("claims", &userJwtClaim)
				whatsAppService.noteService.HandleCreate(ginContext, &allParsedNote[0])
				whatsAppService.SendMessage(message.SenderPhoneNumber, "Catatan berhasil ditambahkan")
			}
		}
		return nil
	})
	helper.LogError(err)
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
		fmt.Printf("failed to send message: %s\n", string(body))
	}
}

func (whatsAppService *ServiceImpl) retrieveMediaLocation(mediaId string) (string, error) {
	fmt.Println("Retrieving media location")
	url := fmt.Sprintf("https://graph.facebook.com/v21.0/%s", mediaId)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", whatsAppService.viperConfig.GetString("META_GRAPH_API_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error response: %s", resp.Status)
	}

	var result struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	fmt.Println("Success retrieve media location")

	return result.URL, nil
}

func (whatsAppService *ServiceImpl) downloadMedia(mediaId string, mediaUrl string, mimeType string) string {
	projectRoot, _ := os.Getwd() // Mendapatkan root path proyek
	templateFile := fmt.Sprintf("%s/public/static/image_temp", projectRoot)

	req, err := http.NewRequest("GET", mediaUrl, nil)
	if err != nil {
		fmt.Println("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", whatsAppService.viperConfig.GetString("META_GRAPH_API_TOKEN")))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("error response: %s", resp.Status)
	}

	backgroundCtx := context.Background()
	file, err := os.Create(fmt.Sprintf("%s/%s.%s", templateFile, mediaId, mimeType))
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("failed to save file: %w", err)
	}

	bucketObject := whatsAppService.googleCloudStorage.StorageClient.Bucket("mnemosyne-bucket")
	helper.LogError(err)
	fileName := fmt.Sprintf("%s.%s", mediaId, mimeType)
	imageObject := bucketObject.Object(fileName)
	fmt.Println(imageObject)
	writer := imageObject.NewWriter(backgroundCtx)
	imageFile, err := os.Open(fmt.Sprintf("%s/%s.%s", templateFile, mediaId, mimeType))
	defer imageFile.Close()
	helper.LogError(err)
	_, err = io.Copy(writer, imageFile)
	if err != nil {
		fmt.Printf("failed to copy data to GCS writer: %v", err)
	}
	err = writer.Close() // Pastikan writer ditutup untuk menyelesaikan upload
	if err != nil {
		fmt.Printf("failed to close GCS writer: %v", err)
	}
	helper.LogError(err)
	fmt.Println("Success download media location")
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", whatsAppService.viperConfig.GetString("BUCKET_NAME"), fileName)
}
