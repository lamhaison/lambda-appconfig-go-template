package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var env string
	env = getEnvOrDefault("ENV", "demo")
	log.Println("Env " + env)

	var pathParams []string
	if request.PathParameters != nil {
		pathParams = strings.Split(request.PathParameters["proxy"], "/")

	} else {
		pathParams = []string{"demo", "default"}
	}

	var appConfigName string
	var appConfigProfile string
	var appConfigEnv = env

	appConfigName = pathParams[0]
	appConfigProfile = pathParams[1]
	appConfigEnv = env

	response, err := http.Get("http://localhost:2772/applications/" + appConfigName + "/environments/" + appConfigEnv + "/configurations/" + appConfigProfile)
	HandleError(err)

	bodyBytes, err := io.ReadAll(response.Body)
	HandleError(err)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(bodyBytes),
	}, nil

}

func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
