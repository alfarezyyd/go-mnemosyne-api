package whatsapp

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/spf13/viper"
	"go-mnemosyne-api/whatsapp/dto"
	"gorm.io/gorm"
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

func (whatsAppService *ServiceImpl) HandleProcessWebhook(ginContext *gin.Context, payloadMessageDto *dto.PayloadMessageDto) {

}

func (whatsAppService *ServiceImpl) HandleCreate(ginContext *gin.Context) {}
