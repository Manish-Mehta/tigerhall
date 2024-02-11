package email

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type awsEmailClient struct {
	client *ses.SES
}

func (awsEmailClient *awsEmailClient) CreateClient(params ClientParam) {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(params.Region),
	})
	if err != nil {
		// FATAL ERROR
		log.Fatal("Error in Creating AWS Session")
	}

	awsEmailClient.client = ses.New(sess)
}

func (awsEmailClient *awsEmailClient) CreateEmailInputs(emailComposition EmailComposition) (interface{}, error) {

	CharSet := "UTF-8"

	return &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses:  emailComposition.CcAddresses,
			ToAddresses:  emailComposition.ToAddresses,
			BccAddresses: emailComposition.BccAddresses,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailComposition.Message.Body.Html),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailComposition.Message.Body.Text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(emailComposition.Message.Subject),
			},
		},
		Source: aws.String(emailComposition.FromAddress),
		// Uncomment to use a configuration set
		// ConfigurationSetName: aws.String(ConfigurationSet),
	}, nil

}

// Returns true if email was successfully sent otherwise log and returns false
func (awsEmailClient *awsEmailClient) SendEmail(emailInput interface{}) bool {
	result, err := awsEmailClient.client.SendEmail(emailInput.(*ses.SendEmailInput))

	// Display error messages if they occur.
	if err != nil {
		var awsErr string
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				awsErr = ses.ErrCodeMessageRejected + " | " + aerr.Error()
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				awsErr = ses.ErrCodeMailFromDomainNotVerifiedException + " | " + aerr.Error()
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				awsErr = ses.ErrCodeConfigurationSetDoesNotExistException + " | " + aerr.Error()
			default:
				awsErr = aerr.Error()
			}
		} else {
			awsErr = aerr.Error()
		}
		log.Println("AWS SES error while sending email: " + awsErr)
		return false
	}

	log.Println("Email sent, id: " + *result.MessageId)
	return true
}
