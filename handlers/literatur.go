package handlers

import (
	"encoding/json"
	literaturdto "literatur/dto/literatur"
	dto "literatur/dto/result"
	"literatur/models"
	"literatur/repositories"
	"net/http"
	"strconv"

	"os"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type handlerLiteratur struct {
	LiteraturRepository repositories.LiteraturRepository
}

// Create `path_file` Global variable here ...
var PathFile = os.Getenv("PATH_FILE")

func HandlerLiteratur(LiteraturRepository repositories.LiteraturRepository) *handlerLiteratur {
	return &handlerLiteratur{LiteraturRepository}
}

func (h *handlerLiteratur) FindLiteraturs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	literaturs, err := h.LiteraturRepository.FindLiteraturs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Embed Path File on Image property here ...
	for i, p := range literaturs {
		literaturs[i].Image = os.Getenv("PATH_FILE") + p.Image
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: literaturs}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerLiteratur) GetLiteratur(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var literatur models.Literatur
	literatur, err := h.LiteraturRepository.GetLiteratur(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create Embed Path File on Image property here ...
	literatur.Image = os.Getenv("PATH_FILE") + literatur.Image

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: convertResponseLiteratur(literatur)}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerLiteratur) CreateLiteratur(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// get data user token
	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userId := int(userInfo["id"].(float64))

	// Get dataFile from midleware and store to filename variable here ...
	dataContex := r.Context().Value("dataFile") // add this code
	filename := dataContex.(string)             // add this code

	// price, _ := strconv.Atoi(r.FormValue("price"))
	// qty, _ := strconv.Atoi(r.FormValue("qty"))
	// [4,8,2] => category id
	// var categoriesId []int
	// for _, r := range r.FormValue("category_id") {
	// 	if int(r-'0') >= 0 {
	// 		categoriesId = append(categoriesId, int(r-'0'))
	// 	}
	// }

	request := literaturdto.LiteraturRequest{
		Title:           r.FormValue("title"),
		PublicationDate: r.FormValue("publicationdate"),
		Pages:           r.FormValue("pages"),
		ISBN:            r.FormValue("isbn"),
		Author:          r.FormValue("author"),
		Image:           r.FormValue("image"),
	}

	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// category, _ := h.LiteraturRepository.FindCategoriesById(categoriesId)

	literatur := models.Literatur{
		Title:           request.Title,
		UserID:          userId,
		PublicationDate: request.PublicationDate,
		Pages:           request.Pages,
		ISBN:            request.ISBN,
		Author:          request.Author,
		Image:           filename,
	}

	literatur, err = h.LiteraturRepository.CreateLiteratur(literatur)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	literatur, _ = h.LiteraturRepository.GetLiteratur(literatur.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: literatur}
	json.NewEncoder(w).Encode(response)
}

func convertResponseLiteratur(u models.Literatur) models.LiteraturResponse {
	return models.LiteraturResponse{
		Title:           u.Title,
		User:            u.User,
		PublicationDate: u.PublicationDate,
		Pages:           u.PublicationDate,
		ISBN:            u.ISBN,
		Author:          u.Author,
		Image:           u.Image,
	}
}
