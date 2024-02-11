package email

var SupportedServices = map[string]string{
	"Brevo": "brevo",
	"AWS":   "aws",
	// Add more service in future here
}

type EmailClientInterface interface {
	CreateClient(ClientParam)
	CreateEmailInputs(EmailComposition) (interface{}, error)
	SendEmail(interface{}) bool
}

type ClientParam struct {
	Region string
}

type Body struct {
	Html string
	Text string
}

type Message struct {
	Body    Body
	Subject string
}

type EmailComposition struct {
	ToAddresses  []*string
	BccAddresses []*string
	CcAddresses  []*string
	Message      Message
	FromAddress  string
}

var serviceMap = map[string]EmailClientInterface{}

var InitService = func(service string) {
	if service == SupportedServices["Brevo"] {
		serviceMap[service] = &brevoEmailClient{}
	}
}

var GetServiceClient = func(service string) EmailClientInterface {
	return serviceMap[service]
}
