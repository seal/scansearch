package serp

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/seal/scansearch/pkg/database"
	"github.com/seal/scansearch/pkg/models"
	"github.com/seal/scansearch/pkg/types"
	"github.com/seal/scansearch/pkg/utils"
)

func ScanPrices(WaitTimeLike string, WaitTimeLove string) {
	go func() {
		ContinuousScan(WaitTimeLike, false)
	}()
	go func() {
		ContinuousScan(WaitTimeLove, true)
	}()
}
func SpecificRequest(url string, price float64) (error, types.SerpProductApiResponse) { // Return true if price is the correct one desired, float64 is the current price
	var s types.SerpProductApiResponse
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utils.Error(err)
		return err, s
	}
	q := req.URL.Query()
	q.Add("apikey", utils.EnvVariable("serpapi"))
	req.URL.RawQuery = q.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Error(err)
		return err, s
	}
	err = json.NewDecoder(resp.Body).Decode(&s)
	if err != nil {
		return err, s
	}
	return nil, s
}
func ContinuousScan(WaitTime string, like bool) {

	for {

		var Wardrobe = []models.Wardrobe{}
		database.Instance.Find(&Wardrobe)
		// Like is true, we want likes only
		var LikedWardrobe = []models.Wardrobe{}

		var LovedWardrobe = []models.Wardrobe{}
		for _, v := range Wardrobe {
			if v.Like {
				LikedWardrobe = append(LikedWardrobe, v)
			} else {

				LovedWardrobe = append(LovedWardrobe, v)

			}
		}
		switch like {
		case true:
			for _, v := range LikedWardrobe { // For each liked item, check the wardrobe and
				err, ResponseStruct := SpecificRequest(v.SerpapiProductApi, v.DesiredPrice)
				if err != nil {
					utils.Error(err)
					return
				}

				currencySymbols := []string{"$", "€", "£", "¥", "₹"}
				var Prices []string
				for _, y := range ResponseStruct.ProductResults.Prices {

					for _, v := range currencySymbols {
						y = strings.Replace(y, v, "", -1)
					}
					Prices = append(Prices, y)
				}
				for k := range Prices {
					price := Prices[k]
					//price = Regex.Replace(price, "[^0-9.]", "")
					m1 := regexp.MustCompile(`[^0-9.]`)

					price = m1.ReplaceAllString(price, "")

					ParsedPrice, err := strconv.ParseFloat(price, 64)
					if err != nil {
						utils.Error(err)
					}

					log.Println("Parsed price:", ParsedPrice, "desired", v.DesiredPrice)
					if ParsedPrice < v.DesiredPrice {
						// Send email
						var user models.User
						result := database.Instance.First(&user, "id= ?", v.UserID)

						log.Println("Sending email to user", user.Email, " Price on item has dropped, their desired price is:", v.DesiredPrice, "The current price is : ", ParsedPrice)
						if result.Error != nil {
							err := errors.New("Error grabbing user for wardrobe price drop" + result.Error.Error())
							utils.Error(err)
							return
						}
						emailData := utils.EmailDataPriceDropped{
							URL:          ResponseStruct.SearchMetadata.GoogleProductURL,
							FirstName:    user.FirstName,
							Subject:      "The price has dropped on an item you have liked",
							DesiredPrice: v.DesiredPrice,
							CurrentPrice: ParsedPrice,
							ImageURL:     ResponseStruct.ProductResults.Media[0].Link,
						}

						utils.SendEmailPrice(&user, &emailData)
						updatedWardrobe := models.Wardrobe{
							DesiredPrice: ParsedPrice, // Stop spamming emails
						}
						database.Instance.Model(&updatedWardrobe).Where("id=?", v.ID).Updates(&updatedWardrobe)
						// Should set the desired price to lower
					}
				}

			}
		case false:

			for _, v := range LovedWardrobe { // For each liked item, check the wardrobe and
				err, ResponseStruct := SpecificRequest(v.SerpapiProductApi, v.DesiredPrice)
				if err != nil {
					utils.Error(err)
					return
				}

				currencySymbols := []string{"$", "€", "£"}
				var Prices []string
				for _, y := range ResponseStruct.ProductResults.Prices {

					for _, v := range currencySymbols {
						y = strings.Replace(y, v, "", -1)
					}
					Prices = append(Prices, y)
				}
				for k := range Prices {

					price := Prices[k]
					//price = Regex.Replace(price, "[^0-9.]", "")
					m1 := regexp.MustCompile(`[^0-9.]`)

					price = m1.ReplaceAllString(price, "")

					ParsedPrice, err := strconv.ParseFloat(price, 64)
					if err != nil {
						utils.Error(err)
					}

					log.Println("Parsed price:", ParsedPrice, "desired", v.DesiredPrice)
					if ParsedPrice < v.DesiredPrice {
						// Send email
						var user models.User
						result := database.Instance.First(&user, "id= ?", v.UserID)

						log.Println("Sending email to user", user.Email, " Price on item has dropped, their desired price is:", v.DesiredPrice, "The current price is : ", ParsedPrice)
						if result.Error != nil {
							err := errors.New("Error grabbing user for wardrobe price drop" + result.Error.Error())
							utils.Error(err)
							return
						}

						emailData := utils.EmailDataPriceDropped{
							URL:          ResponseStruct.SearchMetadata.GoogleProductURL,
							FirstName:    user.FirstName,
							Subject:      "The price has dropped on an item you have liked",
							DesiredPrice: v.DesiredPrice,
							CurrentPrice: ParsedPrice,
							ImageURL:     ResponseStruct.ProductResults.Media[0].Link,
						}
						utils.SendEmailPrice(&user, &emailData)

						updatedWardrobe := models.Wardrobe{
							DesiredPrice: ParsedPrice, // Stop spamming emails
						}
						database.Instance.Model(&updatedWardrobe).Where("id=?", v.ID).Updates(&updatedWardrobe)
					}
				}
			}
		}

		WaitTimeParse, err := strconv.Atoi(WaitTime)
		if err != nil {
			utils.Error(errors.New("Wait time unable to be parsed" + err.Error()))
			return
		}
		time.Sleep(time.Second * time.Duration(WaitTimeParse))
	}
}
