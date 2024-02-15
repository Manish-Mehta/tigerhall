package sight

import (
	"log"
	"net/http"

	"github.com/Manish-Mehta/tigerhall/api/dto"
	"github.com/Manish-Mehta/tigerhall/internal/config"
	"github.com/Manish-Mehta/tigerhall/internal/util"
	"github.com/Manish-Mehta/tigerhall/model/datastore"
	"github.com/Manish-Mehta/tigerhall/model/entities"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
	imageHandler "github.com/Manish-Mehta/tigerhall/pkg/image-handler"
)

type SightService interface {
	create(*dto.CreateSightingRequest, string) *errorHandler.Error
	Create(*dto.CreateSightingRequest) *errorHandler.Error
	List(page int, limit int) (*[]*dto.ListSightingResponse, *errorHandler.Error)
}
type sightService struct {
	dataStore  datastore.SightStore
	tigerStore datastore.TigerStore
}

func (service *sightService) Create(request *dto.CreateSightingRequest) *errorHandler.Error {
	validator := util.NewValidator()

	err := validator.ValCoord(dto.Coordinate{Lat: request.Lat, Lon: request.Lon})
	if err != nil {
		return &errorHandler.Error{
			Err:        "Invalid coordinates",
			ErrMsg:     "Provide valid coordinates",
			StatusCode: http.StatusInternalServerError,
		}
	}

	imgType, err := validator.ValImage(request.Image)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Invalid Image",
			ErrMsg:     err.Error(),
			StatusCode: http.StatusBadRequest,
		}
	}

	sightEntity := &entities.Sight{}
	err = service.dataStore.GetLatest(sightEntity, &entities.Sight{TigerID: request.TigerID}, []string{"id", "lat", "lon", "seen_at"})
	if err != nil {
		return &errorHandler.Error{
			Err:        "Sight fetch failed",
			ErrMsg:     "Error while getting Sight data",
			StatusCode: http.StatusInternalServerError,
		}
	}

	if sightEntity.ID == 0 {
		log.Println("No sighting exist, creating first")
		return service.create(request, imgType)
	}

	if sightEntity.SeenAt.Compare(request.SeenAt) != -1 {
		return &errorHandler.Error{
			Err:        "Invalid seen time",
			ErrMsg:     "Seen at time is older than previous sighting",
			StatusCode: http.StatusBadRequest,
		}
	}

	distance, err := service.dataStore.GetDistance(
		dto.Coordinate{Lat: sightEntity.Lat, Lon: sightEntity.Lon},
		dto.Coordinate{Lat: request.Lat, Lon: request.Lon},
	)
	if err != nil {
		log.Println(err)
		return &errorHandler.Error{
			Err:        "Sight check failed",
			ErrMsg:     "Error while creating sight in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}
	if distance <= 5 {
		return &errorHandler.Error{
			Err:        "Previous sight exists within 5km",
			ErrMsg:     "Sighting should be further than 5 km from the previous one",
			StatusCode: http.StatusConflict,
		}
	}

	return service.create(request, imgType)
}

func (service *sightService) create(request *dto.CreateSightingRequest, imgType string) *errorHandler.Error {

	filePath, imgErr := imageHandler.ProcessImage(request.Image, request.TigerID, imgType)
	if imgErr != nil {
		return imgErr
	}

	sightEntity := &entities.Sight{
		TigerID:  request.TigerID,
		UserID:   request.UserID,
		Lat:      request.Lat,
		Lon:      request.Lon,
		SeenAt:   request.SeenAt,
		ImageURL: filePath,
	}
	err := service.dataStore.Create(sightEntity)
	if err != nil {
		return &errorHandler.Error{
			Err:        "Sight creation failed",
			ErrMsg:     "Error while creating user in the system",
			StatusCode: http.StatusInternalServerError,
		}
	}

	service.tigerStore.Update(
		&entities.Tiger{
			ID:       sightEntity.TigerID,
			LastSeen: sightEntity.SeenAt,
			Lat:      sightEntity.Lat,
			Lon:      sightEntity.Lon,
		},
	)

	// Start the event of email sending Asynch(Chan) with Queue
	go func(id uint) { config.TIGER_SIGHTING_CHAN <- id }(sightEntity.TigerID)
	return nil
}

func (service *sightService) List(page int, limit int) (*[]*dto.ListSightingResponse, *errorHandler.Error) {

	var sightings []*entities.Sight
	var sightingRes []*dto.ListSightingResponse

	err := service.dataStore.List(&sightings, page, limit, []string{"id", "lat", "lon", "seen_at", "image_url", "tiger_id", "user_id"})
	if err != nil {
		log.Println(err)
		return nil, &errorHandler.Error{
			Err:        "Sighting fetch failed",
			ErrMsg:     "Error while getting Sights",
			StatusCode: http.StatusInternalServerError,
		}
	}
	for _, sighting := range sightings {
		sightingResponse := &dto.ListSightingResponse{
			ID:       sighting.ID,
			Lat:      sighting.Lat,
			Lon:      sighting.Lon,
			SeenAt:   sighting.SeenAt,
			ImageURL: sighting.ImageURL,
			TigerID:  sighting.TigerID,
			UserID:   sighting.UserID,
		}
		sightingRes = append(sightingRes, sightingResponse)
	}
	return &sightingRes, nil
}

func NewSightService(ds datastore.SightStore, ts datastore.TigerStore) SightService {
	return &sightService{
		dataStore:  ds,
		tigerStore: ts,
	}
}
