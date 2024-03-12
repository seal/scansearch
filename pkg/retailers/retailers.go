package retailers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/seal/scansearch/pkg/serp"
	"github.com/seal/scansearch/pkg/types"
)

type sizeTbsJson struct {
	Size string `json:"size"`
	URL  string `json:"url,omitempty"`
	Tbs  string `json:"tbs"`
}

func ParseSizeTbs() {
	file, err := os.Open("tbs.json")
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var tbs []sizeTbsJson
	err = json.Unmarshal(body, &tbs)
	if err != nil {
		panic(err)
	}
	var tbsBLANK []sizeTbsJson
	for _, v := range tbs {
		var tbsSINGLE sizeTbsJson
		tbsSINGLE.Size = v.Size
		u, err := url.Parse(v.URL)
		if err != nil {
			panic(err)
		}
		tbsBeforeSplit := u.Query().Get("tbs")
		splitTBS := strings.Replace(tbsBeforeSplit, "mr:1,", "",-1)
		tbsSINGLE.Tbs  = splitTBS
		tbsBLANK = append(tbsBLANK, tbsSINGLE)
	}
	for _, v := range tbsBLANK {
		log.Println(v.Size, v.Tbs)
	}
	ToWrite, err := json.Marshal(tbsBLANK)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("tbsfiltered.json", ToWrite, 0644)
	if err != nil {
		panic(err)
	}

}
func ParseRetailers(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	var Retailers []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 5 {
			continue // new lines
		}
		// Format provided is www.domain.whatever
		split := strings.Split(line, ".")

		Retailers = append(Retailers, split[1])
	}
	log.Println(Retailers)
	var failedretailers []string
	var tbstostore []string
	for _, v := range Retailers {
		var Response types.TbsResponse
		query := make(map[string]string)
		query["q"] = v + " shirt"

		body, err := serp.ApiRequest(query)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, &Response)
		var options []string
		var tbs []string
		for _, y := range Response.Filters {
			if y.Type == "Seller" {
				log.Println("Printing options:")
				//log.Print(y.Options)
				for zx, zy := range y.Options {
					log.Println("Key:", zx, "value", zy)
					options = append(options, zy.Text)
					tbs = append(tbs, zy.Tbs)
				}
				log.Println("Please enter a tbs option for retailer:", v)
				var i int
				_, err := fmt.Scanf("%d", &i)
				if err != nil {
					log.Println(err)
				}
				if i == 50 {
					failedretailers = append(failedretailers, v)
					continue
				}
				tbstostore = append(tbstostore, tbs[i])
				log.Println("Storing ", tbs[i])
			}
		}
	}
	log.Println("Failed retailers:", failedretailers)
	f, err := os.Create("tbs.txt")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	for _, v := range tbstostore {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(500 * time.Second)
}
