package helpers

import (
	"log"

	"github.com/fritz-immanuel/eral-promo-library-go/configs"
	"github.com/fritz-immanuel/eral-promo-library-go/library/client"
)

func HitExternalClient() *client.HTTPClient {
	config, errConfig := configs.GetConfiguration()
	if errConfig != nil {
		log.Fatalln("failed to get configuration: ", errConfig)
	}

	authToken := client.Bearer
	authToken.Token = config.ExternalToken

	accessToken := client.AccessToken
	accessToken.Token = config.ExternalAccessToken

	retailsoftClient := client.NewHTTPClient(client.HTTPClient{
		APIURL: config.ExternalURL,
		AuthorizationTypes: []client.AuthorizationType{
			authToken,
			accessToken,
		},
		MaxNetworkRetries: 3,
		ClientName:        "servicesales",
	})

	return retailsoftClient
}
