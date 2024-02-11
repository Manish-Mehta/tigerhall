package resty

import (
	"github.com/go-resty/resty/v2"
)

type MockRestyClient struct {
}

var GetClientRestyClientMock = func() *resty.Client {
	return nil
}

var SendRestyClientMock = func(req *resty.Request, url string) (*resty.Response, error) {
	return nil, nil
}

var CheckResponseRestyClientMock = func(restyRes *resty.Response, err error, expectedStatusCode int) ([]byte, error) {
	return nil, nil
}

func (rc *MockRestyClient) GetClient() *resty.Client {
	return GetClientRestyClientMock()
}

func (rc *MockRestyClient) Send(req *resty.Request, url string, reqType requestType) (*resty.Response, error) {
	return SendRestyClientMock(req, url)
}

func (rc *MockRestyClient) CheckResponse(restyRes *resty.Response, err error, expectedStatusCode int, serviceName string) ([]byte, error) {
	return CheckResponseRestyClientMock(restyRes, err, expectedStatusCode)
}
