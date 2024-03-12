package serp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/seal/scansearch/pkg/types"
	"github.com/seal/scansearch/pkg/utils"
)

func Standard(RecievedParams map[string]string) types.TbsResponse {
	SerpParams := make(map[string]string)
	// Here we assume the only param we have is query
	SerpParams["q"] = RecievedParams["query"]
	SerpParams["gl"] = RecievedParams["country"]
	body, err := ApiRequest(SerpParams)
	if err != nil {
		// Do proper stuff
		return utils.SerpError(err)
	}
	var SearchOld types.SearchResponse
	json.Unmarshal(body, &SearchOld)
	Search := RemoveTbs(SearchOld)
	var Response types.TbsResponse
	Response.Success = true
	Response.Message = "Success"
	Response.ShoppingResults = utils.Sort(RecievedParams, Search.ShoppingResults)
	Response.Filters = Search.Filters
	return Response
	/*
		response := serp.ApiRequest(SerpParams)
		var Search types.SearchResponse
		json.Unmarshal(response, &Search)
		u, err := json.Marshal(Search.ShoppingResults)
		if err != nil {
			utils.Error(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(u)

	*/
}
func RemoveTbs(inc types.SearchResponse) types.SearchResponse { //Actually remove merchagg
	out := inc
	for x, v := range inc.Filters {
		if v.Type == "Seller" {
			continue
		}
		for y, k := range v.Options {
			if strings.Contains(k.Tbs, "merchagg") {
				Split := strings.Split(k.Tbs, ",merchagg:")
				Split2 := strings.Split(Split[1], ",")
				Split2[0] = Split[0]
				New := strings.Join(Split2, ",")
				out.Filters[x].Options[y].Tbs = New
			}
		}
	}
	return out
}
func PriceIncluded(RecievedParams map[string]string) types.TbsResponse {
	// Switch the min / max price and adjust properly then call ApiRequest
	SerpParams := make(map[string]string)
	var tbs string
	if RecievedParams["minprice"] != "" && RecievedParams["maxprice"] != "" {
		tbs = tbs + fmt.Sprintf("price:1,ppr_min:%s,ppr_max:%s", RecievedParams["minprice"], RecievedParams["maxprice"])
	} else if RecievedParams["minprice"] != "" && RecievedParams["maxprice"] == "" {
		// min price only
		tbs = tbs + fmt.Sprintf("price:1,ppr_min:%s", RecievedParams["minprice"])
	} else if RecievedParams["minprice"] == "" && RecievedParams["maxprice"] != "" {
		// max price only
		tbs = tbs + fmt.Sprintf("price:1,ppr_max:%s", RecievedParams["maxprice"])
	}
	// No mr:1 at this point
	if RecTbs, ok := RecievedParams["tbs"]; ok {
		tbs = tbs + RecTbs
	}
	SerpParams["q"] = RecievedParams["query"] // Need a query lol
	if !strings.Contains(strings.ToLower(tbs), "mr;1") {

		tbs = "mr;1," + tbs
	}
	SerpParams["tbs"] = tbs
	SerpParams["gl"] = RecievedParams["country"]
	body, err := ApiRequest(SerpParams)
	if err != nil {
		return utils.SerpError(err)

	}
	var SearchOld types.SearchResponse
	json.Unmarshal(body, &SearchOld)
	Search := RemoveTbs(SearchOld)
	var Response types.TbsResponse
	Response.Success = true // Very nice, all works
	Response.Message = "Success"
	Response.Filters = Search.Filters
	Response.ShoppingResults = utils.Sort(RecievedParams, Search.ShoppingResults)
	return Response
}

func TbsIncluded(RecievedParams map[string]string) types.TbsResponse {
	// Switch the tbs ensuring no repeat of mr; etc, call ApiRequest and return response
	// Min / max price will not be there
	SerpParams := make(map[string]string)
	SerpParams["gl"] = RecievedParams["country"]
	SerpParams["q"] = RecievedParams["query"]
	SerpParams["tbs"] = "mr;1," + RecievedParams["tbs"]
	body, err := ApiRequest(SerpParams)
	if err != nil {
		return utils.SerpError(err)
	}
	var SearchOld types.SearchResponse
	json.Unmarshal(body, &SearchOld)
	Search := RemoveTbs(SearchOld)
	var Response types.TbsResponse
	Response.Success = true
	Response.Message = "Success"
	Response.Filters = Search.Filters
	Response.ShoppingResults = utils.Sort(RecievedParams, Search.ShoppingResults)

	return Response
}
func ApiRequest(params map[string]string) ([]byte, error) { // Must have all correct params
	req, err := http.NewRequest("GET", "https://serpapi.com/search.json", nil)
	if err != nil {
		utils.Error(err)
		return nil, err
	}
	params["apikey"] = utils.EnvVariable("serpapi")
	//params["apikey"] = "774222bfb26a0fa940ccfdf6e5e86e4814c1e5f843409e0c0f7b012e77718c1c"
	params["tbm"] = "shop" // Prevent from adding in 5 diff funcs
	q := req.URL.Query()
	var tbsadded bool // merchant
	for k, v := range params {
		if k == "tbs" {
			tbsadded= true
			if strings.Contains(v, "merchagg") {
				q.Add(k, v)
				continue
			} else {
				q.Add(k, v+",merchagg:"+utils.EnvVariable("TbsSorted"))
				continue
			}
		}
		q.Add(k, v)
	}
	if !tbsadded{
		q.Add("tbs", "mr:1,merchagg:"+utils.EnvVariable("TbsSorted"))
	}
	/*
		tbs, ok := params["tbs"]
		if ok {
			//  if we have a tbs, and it's a specific merchant, we want to add only the proper merchant
			// else if it's say, size, we want to keep our tbs sorted thing
			if strings.Contains(tbs, "merchagg") {
				decoded, err := url.QueryUnescape(tbs)
				if err != nil {
					utils.Error(err)
				}
				q.Add("tbs", decoded)

				// If we have the tbs, and we have the merchant
				// It's most likely URL encoded in the json payload.

			} else {

				q.Add("tbs", tbs+",merchagg:"+utils.EnvVariable("TbsSorted"))
			}

		} else {
			q.Add("tbs", "mr:1,merchagg:"+utils.EnvVariable("TbsSorted"))
		}*/
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.Error(err)
		return nil, err

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Error(err)
		return nil, err
	}

	return body, nil
}
