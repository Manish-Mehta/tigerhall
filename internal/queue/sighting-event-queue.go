package queue

import (
	"log"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/pkg/db"
)

// type SightingEvent struct {
// 	Tiger   Tiger
// 	Users   []User
// 	Comment string
// }

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

func SightingListener(tigerSighting <-chan uint) {
	log.Println("started")
	sightService := NewSightingEventManager()

	for tigerId := range tigerSighting {
		log.Println("tigerSighting event recieved")
		userDetails, err := sightService.GetUsersForTiger(tigerId)

		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("userDetails")
		log.Println(userDetails)

		// Call email queue with tigerId and []userId
	}
}

func EmailListener(emailEvent <-chan config.EmailEvent) {
	// for event := range emailEvent {
	// 	// tigerId := event.TigerID
	// 	// userEmails := event.UserEmails

	// 	// Send Email to All the users
	// }
}
