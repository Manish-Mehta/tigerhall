package queue

import (
	"fmt"
	"log"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
	emailServices "github.com/Manish-Mehta/tigerhall/pkg/email"
)

type SightingEventService interface {
	GetUsersForTiger(tigerID uint) ([]datastore.SighitngUserDetails, error)
}

type sightingEventService struct {
	dataStore datastore.SightStore
}

func (ss *sightingEventService) GetUsersForTiger(tigerID uint) ([]datastore.SighitngUserDetails, error) {
	return ss.dataStore.GetUsersForTigerSight(tigerID)
}

func NewSightingEventService(ds datastore.SightStore) SightingEventService {
	return &sightingEventService{dataStore: ds}
}

func NewSightingEventManager() SightingEventService {
	sightStore := datastore.NewSightStore(db.GetDBClient().GetClient())
	return NewSightingEventService(sightStore)
}

func SightingListener() {
	log.Println("***************************** SightingListener started *****************************")

	sightService := NewSightingEventManager()

	for tigerId := range config.TIGER_SIGHTING_CHAN {
		log.Println("tigerSighting event recieved")
		userDetails, err := sightService.GetUsersForTiger(tigerId)

		if err != nil || len(userDetails) == 0 {
			log.Println(err)
			log.Println("User Details may empty")
			continue
		}

		emailEvent := config.EmailEvent{}
		for i, _ := range userDetails {
			emailEvent.UserEmails = append(emailEvent.UserEmails, &userDetails[i].Email)
		}
		emailEvent.TigerName = userDetails[0].Tname
		config.EMAIL_EVENT_CHAN <- emailEvent
	}
}

func EmailListener() {
	for event := range config.EMAIL_EVENT_CHAN {
		log.Println("***************************** Processing Email Event *****************************")

		tigerName := event.TigerName
		userEmails := event.UserEmails

		emailBody := fmt.Sprintf(`
			<html>
			<body>
				<p>Hello Wildlife Enthusiast,</p>
				<p>We're excited to share that Tiger %s was recently sighted in the wild! Head over to our app to view all the thrilling details.</p>
				<p>Happy exploring!</p>
			</body>
			</html>`,
			tigerName,
		)

		emailComposition := emailServices.EmailComposition{
			ToAddresses: userEmails,
			FromAddress: config.EMAIL_FROM_ADDRESS,
			Message: emailServices.Message{
				Subject: fmt.Sprintf("Tiger %s reappeared", tigerName),
				Body: emailServices.Body{
					Html: emailBody,
					Text: "",
				},
			},
		}

		emailService := emailServices.GetServiceClient(config.EMAIL_SERVICE)
		emailInputs, _ := emailService.CreateEmailInputs(emailComposition)

		success := emailService.SendEmail(emailInputs)
		if !success {
			fmt.Printf("Email sending failed for Tiger %s \n", tigerName)
		} else {
			fmt.Printf("Email sent for Tiger %s \n", tigerName)
		}
	}
}
