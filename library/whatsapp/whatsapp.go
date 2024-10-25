package whatsapp

import (
	"log"
	"net/url"

	"github.com/fritz-immanuel/eral-promo-library-go/configs"
	"github.com/fritz-immanuel/eral-promo-library-go/library/client"
	"github.com/fritz-immanuel/eral-promo-library-go/library/types"
	"github.com/fritz-immanuel/eral-promo-library-go/models"
	"github.com/gin-gonic/gin"
)

/*
SendMessage()  Sends WA message.

Second parameter is the recipient number which must be given in country code format (<country_code><number>).

Third parameter is the message you would like to send.
*/
func SendMessage(ctx *gin.Context, recipientNumber string, message string) *types.Error {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	WAretailsoftClient := client.NewHTTPClient(client.HTTPClient{
		APIURL:            config.SendWhatsappAPI,
		MaxNetworkRetries: 3,
		ClientName:        "servicesales",
	})

	data := url.Values{}
	data.Set("token", config.SendWhatsappToken)
	data.Set("number", recipientNumber)
	data.Set("message", message)

	WAResponse := models.WhatsappResponse{}
	if err := WAretailsoftClient.CallClientFormEncode(ctx, "send_message", "POST", data, &WAResponse, false); err != nil && err.Error != nil {
		return &types.Error{
			Path:       ".WhatsappHelper->SendMessage()",
			StatusCode: err.StatusCode,
			Message:    err.Message,
			Error:      err.Error,
			Type:       "golang-error",
		}
	}

	log.Println(">>> WA Response: ", WAResponse)

	return nil
}

/*
SendImage()  Sends WA message.

Second parameter is the recipient number which must be given in country code format (<country_code><number>).

Third parameter is the image link you would like to send.

Fourth parameter is the caption to the image.
*/
func SendImage(ctx *gin.Context, recipientNumber string, imgURL string, caption string) *types.Error {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	WAretailsoftClient := client.NewHTTPClient(client.HTTPClient{
		APIURL:            config.SendWhatsappAPI,
		MaxNetworkRetries: 3,
		ClientName:        "servicesales",
	})

	data := url.Values{}
	data.Set("token", config.SendWhatsappToken)
	data.Set("number", recipientNumber)
	data.Set("file", imgURL)
	data.Set("caption", caption)

	WAResponse := models.WhatsappResponse{}
	if err := WAretailsoftClient.CallClientFormEncode(ctx, "send_image", "POST", data, &WAResponse, false); err != nil && err.Error != nil {
		return &types.Error{
			Path:       ".WhatsappHelper->SendImage()",
			StatusCode: err.StatusCode,
			Message:    err.Message,
			Error:      err.Error,
			Type:       "golang-error",
		}
	}

	return nil
}
