package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/seal/scansearch/pkg/models"
	"github.com/seal/scansearch/pkg/utils"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(DB *gorm.DB) UserController {
	return UserController{DB}
}

func (uc *UserController) DeleteWardrobe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := ctx.Value("currentUser").(models.User)
	// assume we get same type as tbs params, just json array
	// Body = ["one", "two"]
	var WardrobeIDs []int
	err := json.NewDecoder(r.Body).Decode(&WardrobeIDs)
	if err != nil {
		utils.Error(err)
		utils.HttpError(err, 400, w)
		return
	}
	for _, v := range WardrobeIDs {
		Wardrobe := models.Wardrobe{
			ID:     v,
			UserID: currentUser.ID,
		}
		uc.DB.Where("user_id=?", currentUser.ID).Delete(&Wardrobe)
	}
	fmt.Fprint(w, `{
    "success":true,
    "message":"Success"
    }`)

	w.Header().Set("Content-Type", "application/json")
}
func (uc *UserController) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := ctx.Value("currentUser").(models.User)
	userResponse := &models.UserResponse{
		ID:        currentUser.ID,
		FirstName: currentUser.FirstName,
		LastName:  currentUser.LastName,
		Email:     currentUser.Email,
		Plan:      currentUser.Plan,
		Username:  currentUser.Username,
	}
	response, err := json.Marshal(userResponse)
	if err != nil {
		err = fmt.Errorf("%w : Error marshalling user response into json", err)
		utils.Error(err)
		utils.HttpError(err, 500, w)
		return
	}
	fmt.Fprint(w, string(response))

	w.Header().Set("Content-Type", "application/json")
}
func (uc *UserController) PutWardrobe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := ctx.Value("currentUser").(models.User)

	var payload *models.WardrobeInput
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.Error(err)
		utils.HttpError(err, 400, w)
		return
	}
	newWardrobe := models.Wardrobe{
		Like:              payload.Like,
		DesiredPrice:      payload.DesiredPrice,
		SerpapiProductApi: payload.SerpapiProductApi,
		ImageURL:          payload.ImageURL,
		Position:          payload.Position,
		Title:             payload.Title,
		Link:              payload.Link,
		ProductLink:       payload.ProductLink,
		Source:            payload.Source,
		Price:             payload.Price,
		ExtractedPrice:    payload.ExtractedPrice,
		Rating:            payload.Rating,
		//Extensions:        payload.ProductDetails.Extensions,
		Thumbnail: payload.Thumbnail,
		Delivery:  payload.Delivery,
	}
	if uc.DB.Model(&newWardrobe).Where("id=?", currentUser.ID).Updates(&newWardrobe).RowsAffected == 0 {
		uc.DB.Create(&newWardrobe)
	}
	response, err := json.Marshal(newWardrobe)
	if err != nil {
		err = fmt.Errorf("%w : Error marshalling response for new wardrobe, wardrobe was created successfully", err)
		utils.Error(err)
		utils.HttpError(err, 500, w)
		return
	}
	fmt.Fprint(w, string(response))

	w.Header().Set("Content-Type", "application/json")
}

func (uc *UserController) GetWardrobe(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	currentUser := ctx.Value("currentUser").(models.User)
	userID := currentUser.ID // for getting wardrobe
	var Wardrobe = []models.Wardrobe{}
	uc.DB.Where("user_id = ?", userID).Find(&Wardrobe)
	response, err := json.Marshal(Wardrobe)
	if err != nil {
		utils.Error(errors.New("Error marshalling response wardrobe" + err.Error()))
		utils.HttpError(err, 401, w)
		return
	}
	fmt.Fprint(w, string(response))

	w.Header().Set("Content-Type", "application/json")

}
func (uc *UserController) AddWardrobe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	currentUser := ctx.Value("currentUser").(models.User)
	var payload *models.WardrobeInput
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		utils.Error(err)
		utils.HttpError(err, 400, w)
		return
	}

	newWardrobe := models.Wardrobe{
		UserID:            currentUser.ID,
		Like:              payload.Like,
		DesiredPrice:      payload.DesiredPrice,
		SerpapiProductApi: payload.SerpapiProductApi,
		ImageURL:          payload.ImageURL,

		Position:       payload.Position,
		Title:          payload.Title,
		Link:           payload.Link,
		ProductLink:    payload.ProductLink,
		Source:         payload.Source,
		Price:          payload.Price,
		ExtractedPrice: payload.ExtractedPrice,
		Rating:         payload.Rating,
		//Extensions:     payload.ProductDetails.Extensions,
		Thumbnail: payload.Thumbnail,
		Delivery:  payload.Delivery,
	}
	result := uc.DB.Create(&newWardrobe)
	if result.Error != nil {
		utils.Error(errors.New("Error adding to new wardrobe" + err.Error()))
		utils.HttpError(err, 500, w)
		return
	}
	WardrobeResponse := models.WardrobeResponse{
		ID:                newWardrobe.ID,
		UserID:            currentUser.ID,
		Like:              payload.Like,
		DesiredPrice:      payload.DesiredPrice,
		SerpapiProductApi: payload.SerpapiProductApi,
		ImageURL:          payload.ImageURL,
		Position:          payload.Position,
		Title:             payload.Title,
		Link:              payload.Link,
		ProductLink:       payload.ProductLink,
		Source:            payload.Source,
		Price:             payload.Price,
		ExtractedPrice:    payload.ExtractedPrice,
		Rating:            payload.Rating,
		//Extensions:        payload.ProductDetails.Extensions,
		Thumbnail: payload.Thumbnail,
		Delivery:  payload.Delivery,
	}

	response, err := json.Marshal(WardrobeResponse)
	if err != nil {
		err = errors.New("Error marshalling wardrobe response" + err.Error())
		utils.Error(err)
		utils.HttpError(err, 500, w)
		return
	}
	fmt.Fprint(w, string(response))
	w.Header().Set("Content-Type", "application/json")
}
