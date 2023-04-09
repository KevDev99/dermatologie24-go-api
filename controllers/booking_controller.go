package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/gorilla/mux"
)

func AddBooking() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var booking models.Booking

		// Parse the incoming request as a multipart form with a maximum file size of 32MB
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		files := r.MultipartForm.File["files"] // Get the list of uploaded files from the request
		form := r.Form

		booking.Subject = form.Get("subject")
		booking.Message = form.Get("message")

		configs.DB.Save(&booking)

		// upload files
		go uploadFiles(files, booking.Id)

		utils.SendResponse(rw, http.StatusCreated, "success", map[string]interface{}{"data": booking})
	}
}

func GetBooking() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookingId := mux.Vars(r)["id"]

		var booking models.Booking

		// look up booking and preload "has-many" relation to booking files.
		err := configs.DB.Model(&models.Booking{}).Preload("BookingFiles").First(&booking, bookingId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": booking})
	}
}

func DeleteBooking() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookingId := mux.Vars(r)["id"]

		var booking models.Booking

		// delete booking
		err := configs.DB.Delete(&booking, bookingId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "booking deleted."})
	}
}

func UpdateBooking() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookingId := mux.Vars(r)["id"]

		var booking models.Booking

		// get user and check if user exists
		err := configs.DB.Preload("BookingFiles").First(&booking, bookingId).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// Decode the request body into a User struct
		var updateBooking models.Booking
		err = json.NewDecoder(r.Body).Decode(&updateBooking)
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// update user in DB
		dberr := configs.DB.Model(&booking).Updates(updateBooking).Error

		// show internal error if needed
		if dberr != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": dberr.Error()})
			return
		}

		// return updates user and successfull status code
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": booking})
	}
}

func AddFileToBooking() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		bookingId := mux.Vars(r)["id"]

		var booking models.Booking

		// get booking and check if it exists
		err := configs.DB.Preload("BookingFiles").First(&booking, bookingId).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// parse form
		parseErr := r.ParseMultipartForm(32 << 20) // Parse the incoming request as a multipart form with a maximum file size of 32MB

		if parseErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": parseErr.Error()})
			return
		}

		files := r.MultipartForm.File["files"] // Get the list of uploaded files from the request

		uploadErr := uploadFiles(files, booking.Id)

		if uploadErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": uploadErr.Error()})
			return
		}

		configs.DB.Preload("BookingFiles").First(&booking, bookingId)

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": booking})
	}
}

func DeleteBookingFile() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var bookingFile models.BookingFile

		// get query params
		bookingId := mux.Vars(r)["id"]
		bookingFileId := mux.Vars(r)["bookingFileId"]

		// convert to int
		intBookingFileId, err := strconv.Atoi(bookingFileId)

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// assign converted id to new instance of booking file
		bookingFile.Id = intBookingFileId

		// get booking file dataset
		queryErr := configs.DB.Where("booking_id = ?", bookingId).First(&bookingFile, bookingFileId).Error

		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": queryErr.Error()})
			return
		}

		// delete file
		go deleteFile(bookingFile.FilePath)

		// delete dataset on database
		dbErr := configs.DB.Where("booking_id = ?", bookingId).Delete(&bookingFile).Error
		if dbErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": dbErr})
			return
		}

		// return success
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "file deleted."})
	}
}

func uploadFiles(files []*multipart.FileHeader, bookingId int) error {

	for _, fileHeader := range files {
		var bookingFile models.BookingFile

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// Save the uploaded file to disk
		filepath := "booking-files/" + fileHeader.Filename
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		bookingFile.Name = fileHeader.Filename
		bookingFile.FilePath = filepath
		bookingFile.BookingId = bookingId

		configs.DB.Save(&bookingFile)

	}

	return nil
}

func deleteFile(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
