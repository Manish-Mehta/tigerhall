package mockEmail

import (
	"github.com/Manish-Mehta/tigerhall/pkg/email"
)

type MockAWSEmailClient struct {
}

func (m *MockAWSEmailClient) CreateClient(params email.ClientParam) {
}

var CreateEmailInputsMockAWSEmailClient = func(emailComposition email.EmailComposition) (interface{}, error) {
	return nil, nil
}

func (m *MockAWSEmailClient) CreateEmailInputs(emailComposition email.EmailComposition) (interface{}, error) {
	return CreateEmailInputsMockAWSEmailClient(emailComposition)
}

var SendEmailMockAWSEmailClient = func(emailInput interface{}) bool {
	return false
}

func (m *MockAWSEmailClient) SendEmail(emailInput interface{}) bool {
	return SendEmailMockAWSEmailClient(emailInput)
}
