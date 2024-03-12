package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/seal/scansearch/pkg/api"
	"github.com/seal/scansearch/pkg/database"
	"github.com/seal/scansearch/pkg/serp"
	"github.com/seal/scansearch/pkg/utils"
)

var localhostFlag = flag.Bool("l", false, "true = run on port 5000")

func main() {
	//retailers.ParseRetailers("outputs/retailers.txt")

	//retailers.ParseSizeTbs()
	database.Connect("scan:" + utils.EnvVariable("MYSQL_SCAN_PASS") + "@tcp(161.35.174.192:3306)/scansearch?parseTime=true")
	database.Migrate()
	r := api.GetRouter()
	/*chi.Walk(r, func(method, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		log.Printf("%-10s | %s\n", method, route)
		return nil
	})*/
	serp.ScanPrices(utils.EnvVariable("LIKE_WAIT_TIME"), utils.EnvVariable("LOVE_WAIT_TIME"))

	//time.Sleep(500 * time.Second)
	// Commented out scanning prices for now
	flag.Parse()
	if *localhostFlag {

		err := http.ListenAndServe("localhost:5000",
			r)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else {
		go func() {

			err := http.ListenAndServeTLS(":443",
				"/etc/letsencrypt/live/scanuk.co/fullchain.pem",
				"/etc/letsencrypt/live/scanuk.co/privkey.pem",
				r)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()
		err := http.ListenAndServe(":80",
			r)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
	//log.Fatal(http.ListenAndServe(":8080", r))
	/*
	   Certificate is saved at: /etc/letsencrypt/live/scanuk.co/fullchain.pem
	   Key is saved at:         /etc/letsencrypt/live/scanuk.co/privkey.pem
	*/
}
