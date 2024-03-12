package api

import (
	"encoding/json"
	"fmt"
	"log"

	//"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/seal/scansearch/pkg/database"
	"github.com/seal/scansearch/pkg/routes"

	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/seal/scansearch/pkg/controllers"
	"github.com/seal/scansearch/pkg/serp"
	"github.com/seal/scansearch/pkg/types"
	"github.com/seal/scansearch/pkg/utils"
	"github.com/seal/scansearch/ui"
)

var (
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController
)

func GetRouter() *chi.Mux {
	AuthController = controllers.NewAuthController(database.Instance)
	AuthRouteController = routes.NewAuthRouteController(AuthController)
	UserController = controllers.NewUserController(database.Instance)
	UserRouteController = routes.NewRouteUserController(UserController)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		Debug:            true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Route("/api", func(r chi.Router) {

		r.Route("/", func(r chi.Router) {
			UserRouteController.UserRoute(r)
			AuthRouteController.AuthRouter(r)
		})
		r.Get("/tbs", TbsHandler)

		r.Post("/tbs", TbsPostHandler)
		r.Mount("/admin", AdminRouter())
	})
	r.HandleFunc("/*", indexHandler)
	// static files
	staticFS, err := fs.Sub(ui.StaticFiles, "dist")
	if err != nil {
		log.Println(err, "err here")
	}
	httpFS := http.FileServer(http.FS(staticFS))
	r.Handle("/assets/*", httpFS)
	r.Handle("/images/*", httpFS)
	// Below is a previous version
	/*r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		rawFile, err := ui.StaticFiles.ReadFile("dist/index.html")
		if err != nil {
			log.Println("err in dist/index", err)
		}
		log.Println("IndexHandler used")
		w.Write(rawFile)
	})*/
	/*
		workDir, _ := os.Getwd()
		filesDir := http.Dir(filepath.Join(workDir, "static"))
		FileServer(r, "/", filesDir)
	*/
	return r
}
func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		http.NotFound(w, r)
		return
	}

	if r.URL.Path == "/favicon.ico" {
		rawFile, err := ui.StaticFiles.ReadFile("dist/favicon.ico")
		if err != nil {
			log.Println(err, "error favicon")
		}
		w.Write(rawFile)
		return
	}

	rawFile, err := ui.StaticFiles.ReadFile("dist/index.html")
	if err != nil {
		log.Println("err in dist/index", err)
	}
	w.Write(rawFile)
}

/*
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}*/

// Delete later
var DeclaredParams = []string{"query", "brand", "category", "colour", "minprice", "maxprice"}

func TbsPostHandler(w http.ResponseWriter, r *http.Request) {
	TbsParams := []string{"query", "tbs", "minprice", "maxprice", "sort", "country"}
	RecievedParams := make(map[string]string)
	RecievedParams["country"] = "gb"
	for _, v := range TbsParams {
		if r.URL.Query().Get(v) != "" {
			RecievedParams[v] = r.URL.Query().Get(v)
		}
	}
	/*
		if v, ok := RecievedParams["country"]; !ok {
			RecievedParams["country"] = utils.CountryLoopup(r)
		} else {
			RecievedParams["country"] = v
		}
	*/
	var TbsArr []string
	err := json.NewDecoder(r.Body).Decode(&TbsArr)
	if err != nil {
		utils.Error(err)
		utils.HttpError(err, 500, w)
		return
	}
	for _, v := range TbsArr {
		RecievedParams["tbs"] = RecievedParams["tbs"] + v + ","
	}
	//RecievedParams["tbs"] = RecievedParams["tbs"] + ","
	RecievedParams["tbs"] = strings.TrimSuffix(RecievedParams["tbs"], ",")
	Response := serp.TbsIncluded(RecievedParams) // Marshal returned struct from Tbs response
	err = utils.ReturnJsonTbs(Response, w, 200)  // Pass struct into return json response func
	if err != nil {
		utils.Error(err)             // If err, I broke the json
		utils.HttpError(err, 500, w) // Internal
	}

}
func TbsHandler(w http.ResponseWriter, r *http.Request) {
	TbsParams := []string{"query", "tbs", "minprice", "maxprice", "sort", "country"}
	RecievedParams := make(map[string]string)
	RecievedParams["country"] = "gb"
	for _, v := range TbsParams {
		if r.URL.Query().Get(v) != "" {
			RecievedParams[v] = r.URL.Query().Get(v)
		}
	}
	// Check for country param, if no country loopup IP
	/*
		if v, ok := RecievedParams["country"]; !ok {
			RecievedParams["country"] = utils.CountryLoopup(r)
		} else {
			RecievedParams["country"] = v
		}
	*/
	// Check for query
	if _, ok := RecievedParams["query"]; !ok { // If no query reject req
		var Response types.TbsResponse
		Response.Success = false // Provide proper struct to respond
		Response.Message = "No query provided"
		err := utils.ReturnJsonTbs(Response, w, 200) // Pass struct into tbs func
		if err != nil {
			utils.Error(err) // error marshalling json
			utils.HttpError(err, 500, w)
			return
		}
		return
	}
	_, ok1 := RecievedParams["minprice"]
	_, ok2 := RecievedParams["maxprice"]
	if ok1 || ok2 {
		Response := serp.PriceIncluded(RecievedParams) // Response including min  / max price
		err := utils.ReturnJsonTbs(Response, w, 200)
		if err != nil {
			utils.Error(err)
			utils.HttpError(err, 500, w)
		}
		return
	}

	if _, ok := RecievedParams["tbs"]; ok { // If there's a tbs
		Response := serp.TbsIncluded(RecievedParams) // Marshal returned struct from Tbs response
		err := utils.ReturnJsonTbs(Response, w, 200) // Pass struct into return json response func
		if err != nil {
			utils.Error(err) // If err, I broke the json

			utils.HttpError(err, 500, w)
		}
		return
	}
	// If none are okay, return standard search
	Response := serp.Standard(RecievedParams)
	err := utils.ReturnJsonTbs(Response, w, 200)
	if err != nil {
		utils.Error(err)
		utils.HttpError(err, 500, w)
	}
}

/*
func SearchPostHandle(w http.ResponseWriter, r *http.Request) {
	RecievedParams := make(map[string]string) // Create map of recieved params
	for _, v := range DeclaredParams { // Loop over params add only declared to map
		if r.URL.Query().Get(v) != "" {
			RecievedParams[v] = r.URL.Query().Get(v)
		}
	}
	SerpParams := make(map[string]string)               // Params to send to serp
	SerpParams["apikey"] = utils.EnvVariable("serpapi") // Add api key to params
	SerpParams["tbm"] = "shop"                          // Add tbm to params
	// If no query return error here
	if SerpParams["q"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No query provided") // Change to json
		return
	} else {
		SerpParams["q"] = RecievedParams["query"] // Needs query
	}
	// Set to UK for now
	SerpParams["gl"] = "uk" // Lines below change if not UK
	if RecievedParams["minprice"] != "" && RecievedParams["maxprice"] != "" {
		SerpParams["tbs"] = fmt.Sprintf("mr:1,price:1,ppr_min:%s,ppr_max:%s", RecievedParams["minprice"], RecievedParams["maxprice"])
	} else if RecievedParams["minprice"] != "" && RecievedParams["maxprice"] == "" {
		// min price only
		SerpParams["tbs"] = fmt.Sprintf("mr:1,price:1,ppr_min:%s", RecievedParams["minprice"])
	} else if RecievedParams["minprice"] == "" && RecievedParams["maxprice"] != "" {
		// max price only
		SerpParams["tbs"] = fmt.Sprintf("mr:1,price:1,ppr_max:%s", RecievedParams["maxprice"])
	}
	response:= serp.ApiRequest(SerpParams)
	var Search types.SearchResponse
	json.Unmarshal(response, &Search)
	u, err := json.Marshal(Search.ShoppingResults)
	if err != nil {
		utils.Error(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(u)
}
*/

func AdminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	r.Get("/", AdminAccounts)
	r.Get("/accounts", AdminAccounts)
	/*
		Methods to add:
		/dashboard -> returns tbd, probably total users, total daily searches, etc
	*/
	return r
}

func AdminAccounts(w http.ResponseWriter, r *http.Request) {
	// Example router, for now ignore
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Its ok "))
}

type AdminLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		// e.g.
		// if !isAdmin(r) {
		// 	http.Error(w, http.StatusText(401), 401)
		// 	return
		// }
		IsAdmin, _ := isAdmin(r)
		if !IsAdmin {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else {
			next.ServeHTTP(w, r)
		}

	})
}

func isAdmin(r *http.Request) (bool, AdminLoginResponse) {
	var RemoveLater AdminLoginResponse
	return true, RemoveLater
}
