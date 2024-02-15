package email

import (
	"log"
	"net/http"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/pkg/resty"
)

type brevoEmailClient struct{}

// Does nothing as brevo doesn't have a client and direct API will be used
func (brevoEmailClient *brevoEmailClient) CreateClient(params ClientParam) {
	log.Println("Email service started")
}

func (brevoEmailClient *brevoEmailClient) CreateEmailInputs(emailComposition EmailComposition) (interface{}, error) {

	inputs := map[string]interface{}{
		"sender": map[string]interface{}{
			"name":  "Manish Test",
			"email": emailComposition.FromAddress,
		},
		"subject":     emailComposition.Message.Subject,
		"htmlContent": emailComposition.Message.Body.Html,
	}

	type Email struct {
		Name  string
		Email string
	}

	if len(emailComposition.ToAddresses) > 0 {
		toEmails := []Email{}
		for _, id := range emailComposition.ToAddresses {
			toEmails = append(toEmails, Email{Email: *id, Name: "User"})
		}
		inputs["to"] = toEmails
	}
	if len(emailComposition.BccAddresses) > 0 {
		bccEmails := []Email{}
		for _, id := range emailComposition.BccAddresses {
			bccEmails = append(bccEmails, Email{Email: *id, Name: "User"})
		}
		inputs["bcc"] = bccEmails
	}
	return &inputs, nil

}

// Returns true if email was successfully sent otherwise log and returns false
func (brevoEmailClient *brevoEmailClient) SendEmail(emailInput interface{}) bool {

	rc := resty.GetRestyClient()
	restyRes, err := rc.Send(
		rc.GetClient().R().
			SetHeader("Content-Type", "application/json").
			SetHeader("api-key", config.EMAIL_API_KEY).
			SetBody(emailInput),
		config.EMAIL_API_ENDPOINT,
		resty.POST,
	)

	_, err = rc.CheckResponse(restyRes, err, http.StatusCreated, "Brevo")
	if err != nil {
		log.Println(err)
		log.Println("Email sending failed")
		return false
	}
	log.Println("Email sent")
	return true
}
