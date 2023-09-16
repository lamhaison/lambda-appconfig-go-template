package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jarcoal/httpmock"
	"github.com/valyala/fastjson"
	"log"
	"testing"
)

const MOCK_RESPONSE = `{"enabled":true}`

func TestGetEnvOrDefault(t *testing.T) {
	var env string
	env = getEnvOrDefault("ENV", "dev")
	if env != "dev" {
		t.Error("Failed when getEnvOrDefault " + env)
	}
	log.Println("Test response: " + env)
}

func TestHandler(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET",
		"http://localhost:2772/applications/demo/environments/demo/configurations/default",
		httpmock.NewStringResponder(200, MOCK_RESPONSE))
	resp, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "GET", Body: "{}"})
	if err != nil {
		t.Error(err)
	}
	err = fastjson.Validate(resp.Body)
	if err != nil {
		t.Error(err)
	}
	log.Println("test response: " + resp.Body)

}

func TestHandlerWithPathParam(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET",
		"http://localhost:2772/applications/demo/environments/demo/configurations/default",
		httpmock.NewStringResponder(200, MOCK_RESPONSE))
	pathParamMap := map[string]string{"proxy": "demo/default"}
	resp, err := handler(events.APIGatewayProxyRequest{HTTPMethod: "GET", Body: "{}", PathParameters: pathParamMap})
	if err != nil {
		t.Error(err)
	}
	err = fastjson.Validate(resp.Body)
	if err != nil {
		t.Error(err)
	}
	log.Println("response: " + resp.Body)
}
