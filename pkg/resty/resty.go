package resty

import (
	"encoding/json"
	"fmt" // used only for string formatting
	"log"

	"github.com/go-resty/resty/v2"
)

var client RestyClient = nil

type RestyClient interface {
	GetClient() *resty.Client
	Send(req *resty.Request, url string, reqType requestType) (*resty.Response, error)
	CheckResponse(restyRes *resty.Response, err error, expectedStatusCode int, serviceName string) ([]byte, error)
}

type restyClient struct {
	client *resty.Client
}

var CreateRestyClient = func() {
	if client == nil {
		client = &restyClient{client: resty.New()}
	}
}

var GetRestyClient = func() RestyClient {
	return client
}

func (rc *restyClient) GetClient() *resty.Client {
	return rc.client
}

func (rc *restyClient) Send(req *resty.Request, url string, reqType requestType) (*resty.Response, error) {
	var res *resty.Response
	var err error
	switch reqType {
	case GET:
		res, err = req.Get(url)
	case PUT:
		res, err = req.Put(url)
	case POST:
		res, err = req.Post(url)
	case PATCH:
		res, err = req.Patch(url)
	case DELETE:
		res, err = req.Delete(url)
	default:
		log.Println(fmt.Sprintf("request type: %s is not supported", reqType))
	}

	return res, err
}

func (rc *restyClient) CheckResponse(restyRes *resty.Response, err error, expectedStatusCode int, serviceName string) ([]byte, error) {
	if err != nil {
		log.Println("error getting details from " + serviceName)
		return nil, err
	}

	if restyRes.StatusCode() != expectedStatusCode {
		log.Println(fmt.Sprintf("expected status code %d, got %d, in service: %s", expectedStatusCode, restyRes.StatusCode(), serviceName))

		var error interface{}
		if err := json.Unmarshal(restyRes.Body(), &error); err != nil {
			log.Println(fmt.Sprintf("error parsing %s response: %s", serviceName, string(restyRes.Body())))
			return nil, err
		}

		return nil, fmt.Errorf("%s error: %+v", serviceName, error)
	}

	return restyRes.Body(), nil
}
