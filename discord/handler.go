package discord

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"go-mnemosyne-api/config"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"go-mnemosyne-api/mapper"
	"go-mnemosyne-api/model"
	"go-mnemosyne-api/model/vertex"
	"go-mnemosyne-api/note"
	noteDto "go-mnemosyne-api/note/dto"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	discordService     Service
	dbConnection       *gorm.DB
	vertexClient       *config.VertexClient
	noteService        note.Service
	googleCloudStorage *config.GoogleCloudStorage
	viperConfig        *viper.Viper
}

func NewHandler(discordService Service,
	vertexClient *config.VertexClient,
	noteService note.Service,
	dbConnection *gorm.DB,
	cloudStorage *config.GoogleCloudStorage,
	viperConfig *viper.Viper) *Handler {
	return &Handler{
		discordService:     discordService,
		vertexClient:       vertexClient,
		noteService:        noteService,
		dbConnection:       dbConnection,
		googleCloudStorage: cloudStorage,
		viperConfig:        viperConfig,
	}
}

func (discordHandler *Handler) OnMessageCreate(discSession *discordgo.Session, messagePayload *discordgo.MessageCreate) {
	if messagePayload.Author.ID == discSession.State.User.ID {
		return
	}
	discSession.ChannelMessageSend(messagePayload.ChannelID, "Permintaan sedang diproses!")
	// Periksa apakah pesan memiliki lampiran (media)
	var openingPrompt string
	var content *genai.GenerateContentResponse
	var err error
	discordHandler.dbConnection.Transaction(func(gormTransaction *gorm.DB) error {
		if len(messagePayload.Attachments) > 0 {
			for _, attachment := range messagePayload.Attachments {
				fmt.Printf("File received: %s\n", attachment.Filename)
				fmt.Printf("Download URL: %s\n", attachment.URL)

				// Jika Anda ingin mengunduh file
				discordHandler.downloadMedia(attachment.URL, attachment.Filename)
				filePath, err := discordHandler.uploadIntoCloudStorage(attachment.Filename)
				helper.LogError(err)
				fmt.Printf("File uploaded: %s\n", *(filePath))
				content, err = discordHandler.vertexClient.GenerateContentWithImage(
					*(filePath),
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
			}
		} else {
			if messagePayload.Content != "" {
				if strings.HasPrefix(messagePayload.Content, "/") {
					switch strings.ToLower(strings.Replace(messagePayload.Content, "/", "", 1)) {
					case "getall":
						var allNote []model.Note
						gormTransaction.Joins("JOIN users ON users.id = notes.user_id").Where("users.phone_number = ?", "6289637577001").Find(&allNote)
						discSession.ChannelMessageSend(messagePayload.ChannelID, mapper.MapAllNoteIntoString(allNote))
						break
					}
					return err
				}

				openingPrompt = fmt.Sprintf("Saya memiliki teks seperti ini %s", messagePayload.Content)
				content, err = discordHandler.vertexClient.GenerateContent(
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
			}
		}
		if err != nil {
			return err
		}
		marshalResponse, _ := json.MarshalIndent(content, "", "  ")
		var generateResponse vertex.ContentResponse
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
		var noteModel model.Note
		var isCategoryExists bool
		err = gormTransaction.Where("phone_number = ?", "6289637577001").First(&userModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		err = gormTransaction.Model(&model.Category{}).Select("COUNT(*) > 0").Where("id = ?", 1).Where("user_id = ?", userModel.ID).Find(&isCategoryExists).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		mapper.MapNoteDtoIntoNoteModel(&allParsedNote[0], &noteModel)
		noteModel.UserID = userModel.ID
		err = gormTransaction.Create(&noteModel).Error
		helper.CheckErrorOperation(err, exception.ParseGormError(err))
		discSession.ChannelMessageSend(messagePayload.ChannelID, "Catatan berhasil ditambahkan")
		return nil
	})

	fmt.Println("ONLINe!")
}

// Fungsi untuk mengunduh media dari URL
func (discordHandler *Handler) downloadMedia(url, filename string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to download file: %s\n", err)
		return
	}
	defer resp.Body.Close()
	projectRoot, _ := os.Getwd() // Mendapatkan root path proyek
	templateFile := fmt.Sprintf("%s/public/static/image_temp", projectRoot)

	file, err := os.Create(fmt.Sprintf("%s/%s", templateFile, filename))
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Printf("Failed to save file: %s\n", err)
		return
	}

	fmt.Printf("File %s downloaded successfully.\n", filename)
}

func (discordHandler *Handler) uploadIntoCloudStorage(filename string) (*string, error) {
	projectRoot, _ := os.Getwd() // Mendapatkan root path proyek
	templateFile := fmt.Sprintf("%s/public/static/image_temp", projectRoot)
	open, err := os.Open(fmt.Sprintf("%s/%s", templateFile, filename))
	if err != nil {
		return nil, err
	}
	bucketObject := discordHandler.googleCloudStorage.StorageClient.Bucket("mnemosyne-bucket")
	imageObject := bucketObject.Object(filename)
	backgroundCtx := context.Background()

	writer := imageObject.NewWriter(backgroundCtx)
	_, err = io.Copy(writer, open)
	if err != nil {
		fmt.Printf("failed to copy data to GCS writer: %v", err)
	}
	err = writer.Close() // Pastikan writer ditutup untuk menyelesaikan upload
	if err != nil {
		fmt.Printf("failed to close GCS writer: %v", err)
	}
	helper.LogError(err)
	fmt.Println("Success download media location")
	filePath := fmt.Sprintf("https://storage.googleapis.com/%s/%s", discordHandler.viperConfig.GetString("BUCKET_NAME"), filename)
	return &filePath, nil
}
