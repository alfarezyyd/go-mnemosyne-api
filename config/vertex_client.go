package config

import (
	"cloud.google.com/go/vertexai/genai"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go-mnemosyne-api/exception"
	"go-mnemosyne-api/helper"
	"net/http"
)

type VertexClient struct {
	viperConfig     *viper.Viper
	generativeModel *genai.GenerativeModel
}

func NewVertexClient(viperConfig *viper.Viper) *VertexClient {
	return &VertexClient{
		viperConfig: viperConfig,
	}
}

func (vertexClient *VertexClient) InitializeVertexClient() error {
	if vertexClient.generativeModel == nil {
		projectId := vertexClient.viperConfig.GetString("GOOGLE_CLOUD_PROJECT_ID")
		locationInstance := vertexClient.viperConfig.GetString("GOOGLE_CLOUD_LOCATION")
		backgroundContext := context.Background()
		client, err := genai.NewClient(backgroundContext, projectId, locationInstance)
		if err != nil {
			return fmt.Errorf("error creating client: %w", err)
		}
		vertexClient.generativeModel = client.GenerativeModel("gemini-1.5-pro-002")
		vertexClient.generativeModel.ResponseMIMEType = "application/json"
	}
	return nil
}

func (vertexClient *VertexClient) GenerateContent(promptPayload string) (*genai.GenerateContentResponse, error) {
	backgroundContext := context.Background()
	if vertexClient.generativeModel != nil {
		textPayload := genai.Text(promptPayload)
		resp, err := vertexClient.generativeModel.GenerateContent(backgroundContext, textPayload)
		if err != nil {
			helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
		}

		return resp, nil
	}
	return nil, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest)
}

func (vertexClient *VertexClient) GenerateContentWithImage(imagePath string, promptPayload string) (*genai.GenerateContentResponse, error) {
	ctx := context.Background()
	if vertexClient.generativeModel == nil {
		return nil, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest)
	}

	img := genai.FileData{
		MIMEType: "image/jpeg",
		FileURI:  imagePath,
	}
	textPayload := genai.Text(promptPayload)

	// Generate content with multipart prompt
	resp, err := vertexClient.generativeModel.GenerateContent(ctx, img, textPayload)
	if err != nil {
		helper.CheckErrorOperation(err, exception.NewClientError(http.StatusBadRequest, exception.ErrBadRequest))
	}

	return resp, nil
}
