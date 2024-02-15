package imageHandler

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/nfnt/resize"

	"github.com/Manish-Mehta/tigerhall/internal/config"
	errorHandler "github.com/Manish-Mehta/tigerhall/pkg/error-handler"
)

// Future Scope: Update this function in moduler sub functions as a part of a struct with pointer recievers
// Ideal Scenario: Instead of storing on local storage it should be uploaded to cloud/NAS storage, including the resizing process
func ProcessImage(imageFileHeader *multipart.FileHeader, fileNamePrefix uint, imgType string) (string, *errorHandler.Error) {

	imgFile, err := imageFileHeader.Open()
	if err != nil {
		log.Println(err)
		return "", &errorHandler.Error{
			Err:        "Broken image",
			ErrMsg:     "Image file is not valid",
			StatusCode: http.StatusBadRequest,
		}
	}

	imgBytes := make([]byte, config.MAX_UPLOAD_IMAGE_SIZE)
	if n, err := imgFile.Read(imgBytes); n == 0 || err != nil {
		log.Println(n)
		log.Println(err)
		return "", &errorHandler.Error{
			Err:        "Error in reading image file",
			ErrMsg:     "Could not read the image file",
			StatusCode: http.StatusBadRequest,
		}
	}
	defer imgFile.Close()
	imgReader := bytes.NewReader(imgBytes)

	var imgDecoded image.Image
	if imgType == "jpg" {
		imgDecoded, err = jpeg.Decode(imgReader)
	} else {
		imgDecoded, err = png.Decode(imgReader)
	}
	if err != nil {
		log.Println(err)
		return "", &errorHandler.Error{
			Err:        "Corrupt Image",
			ErrMsg:     "Image file encoding not correct",
			StatusCode: http.StatusBadRequest,
		}
	}

	resizedImg := resize.Resize(250, 200, imgDecoded, resize.Lanczos3)
	fileName := fmt.Sprintf("%d_%d.%s", fileNamePrefix, time.Now().Unix(), imgType)
	// homeDir, err := os.UserHomeDir()

	// filePath := filepath.Join(homeDir, config.IMAGE_STORAGE_PATH, fileName)
	newFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		return "", &errorHandler.Error{
			Err:        "File processing issue",
			ErrMsg:     "Error while processing the image file",
			StatusCode: http.StatusInternalServerError,
		}
	}
	defer newFile.Close()

	if imgType == "jpg" {
		err = jpeg.Encode(newFile, resizedImg, nil)
	} else {
		err = png.Encode(newFile, resizedImg)
	}
	if err != nil {
		log.Println(err)
		return "", &errorHandler.Error{
			Err:        "File processing issue",
			ErrMsg:     "Error while saving the image file",
			StatusCode: http.StatusInternalServerError,
		}
	}
	return filePath, nil
}
