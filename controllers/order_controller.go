package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func AddOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var order models.Order

		userId := int(math.Round(context.Get(r, "userId").(float64)))

		// Parse the incoming request as a multipart form with a maximum file size of 32MB
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		files := r.MultipartForm.File["files"] // Get the list of uploaded files from the request
		form := r.Form

		order.Message = form.Get("message")
		order.PaymentExtId = form.Get("paymentId")
		order.PaymentTypeId = form.Get("paymentTypeId")
		order.UserId = userId
		order.StatusId = 1

		// create order
		configs.DB.Create(&order)

		// upload files
		go uploadFiles(files, order.Id)

		utils.SendResponse(rw, http.StatusCreated, "success", map[string]interface{}{"data": order})
	}
}

func GetOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		orderId := mux.Vars(r)["id"]

		var order models.Order

		// look up order and preload "has-many" relation to order files.
		err := configs.DB.Model(&models.Order{}).Preload("orderFiles").First(&order, orderId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": order})
	}
}

func DeleteOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		orderId := mux.Vars(r)["id"]

		var order models.Order

		// delete order
		err := configs.DB.Delete(&order, orderId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "order deleted."})
	}
}

func UpdateOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		orderId := mux.Vars(r)["id"]

		var order models.Order

		// get user and check if user exists
		err := configs.DB.Preload("orderFiles").First(&order, orderId).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// Decode the request body into a User struct
		var updateorder models.Order
		err = json.NewDecoder(r.Body).Decode(&updateorder)
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// update user in DB
		dberr := configs.DB.Model(&order).Updates(updateorder).Error

		// show internal error if needed
		if dberr != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": dberr.Error()})
			return
		}

		// return updates user and successfull status code
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": order})
	}
}

func AddFileToOrder() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		orderId := mux.Vars(r)["id"]

		var order models.Order

		// get order and check if it exists
		err := configs.DB.Preload("orderFiles").First(&order, orderId).Error
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

		uploadErr := uploadFiles(files, order.Id)

		if uploadErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": uploadErr.Error()})
			return
		}

		configs.DB.Preload("orderFiles").First(&order, orderId)

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": order})
	}
}

func DeleteOrderFile() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var orderFile models.OrderFile

		// get query params
		orderId := mux.Vars(r)["id"]
		orderFileId := mux.Vars(r)["orderFileId"]

		// convert to int
		intorderFileId, err := strconv.Atoi(orderFileId)

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// assign converted id to new instance of order file
		orderFile.Id = intorderFileId

		// get order file dataset
		queryErr := configs.DB.Where("order_id = ?", orderId).First(&orderFile, orderFileId).Error

		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": queryErr.Error()})
			return
		}

		// delete file
		go deleteFile(orderFile.FilePath)

		// delete dataset on database
		dbErr := configs.DB.Where("order_id = ?", orderId).Delete(&orderFile).Error
		if dbErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": dbErr})
			return
		}

		// return success
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "file deleted."})
	}
}

func uploadFiles(files []*multipart.FileHeader, orderId int) error {

	for _, fileHeader := range files {
		var orderFile models.OrderFile

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		// Save the uploaded file to disk
		filepath := "order-files/" + fileHeader.Filename
		f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			return err
		}

		orderFile.Name = fileHeader.Filename
		orderFile.FilePath = filepath
		orderFile.OrderId = orderId

		configs.DB.Save(&orderFile)
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
