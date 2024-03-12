package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/joho/godotenv"
	"github.com/seal/scansearch/pkg/types"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
func CountryLoopup(r *http.Request) string {
	ip := strings.Split(ReadUserIP(r), ":")
	resp, err := http.Get("https://api.iplocation.net?ip=" + ip[0])
	if err != nil {
		Error(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Error(err)
	}
	var IPResponse types.IpLoopup
	json.Unmarshal(body, &IPResponse)
	// Problem here, Google uses wrong type of response codes.
	// Switch statement to country standard?
	// Currently accepting GB as uk, if it breaks in future correct it
	return IPResponse.CountryCode2
}
func Sort(RecievedParams map[string]string, ResponseResults types.ShoppingResults) types.ShoppingResults {
	// take recieved params, and a shopping result as an input
	// Switch the sort param and return results either high to low or low to high, or default
	// default = order of relavency

	v, ok := RecievedParams["sort"]
	if !ok {
		// Gotten no sort param
		return ResponseResults
	}
	switch v {
	case "hl":
		sort.SliceStable(ResponseResults, func(i, j int) bool {
			return ResponseResults[i].ExtractedPrice > ResponseResults[j].ExtractedPrice
		})
		return ResponseResults
	case "lh":
		sort.SliceStable(ResponseResults, func(i, j int) bool {
			return ResponseResults[i].ExtractedPrice < ResponseResults[j].ExtractedPrice
		})
		return ResponseResults
	default:
		return ResponseResults
	}

}
func SerpError(err error) types.TbsResponse {
	var Response types.TbsResponse
	Response.Success = false
	Response.Message = err.Error()
	return Response
}
func ReturnJsonTbs(Response types.TbsResponse, w http.ResponseWriter, status int) error {
	ResponseJson, err := json.Marshal(Response)
	if err != nil {
		Error(err)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(ResponseJson)
	return nil

}
func EnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
func HttpError(err error, StatusCode int, w http.ResponseWriter) {
	ErrorStruct := &ErrorResponse{
		Success: false,
		Message: err.Error(),
	}
	response, err := json.Marshal(ErrorStruct)
	if err != nil {
		Error(err) // Yeah not very good practice
	}

	w.WriteHeader(StatusCode)
	w.Write(response)
}

type ErrorResponse struct { // Add to models package later
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Error(err error) {
	log.Println(err)
	SendTelegramMessage(err.Error())
}

func SendTelegramMessage(msg string) {

	url := "https://api.telegram.org/bot" + EnvVariable("TELEGRAM_BOT_TOKEN") + "/sendMessage?chat_id=" + EnvVariable("TELEGRAM_CHAT_ID") + "&text=" + msg
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
}
