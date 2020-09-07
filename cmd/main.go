package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"vobe-auth/catchers"
	"vobe-auth/cmd/db"
	"vobe-auth/cmd/status"
	"vobe-auth/routes"
)

type RouteMap map[string]func(w *http.Request) status.Status
type CatchMap map[int16]func(msg string) status.Status

var routeHandler = make(RouteMap)
var catchHandler = make(CatchMap)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println(fmt.Sprintf("[%s]\t<%s>\t%s", r.Method, r.RemoteAddr, r.URL.Path))

	path := strings.ToLower(r.URL.Path[1:])
	route := routeHandler[path]

	if route == nil {
		sendRes(w, catchHandler[404]("Resource not found"))
		return
	}

	sendRes(w, route(r))
}

func sendRes(w http.ResponseWriter, status status.Status) {
	header := w.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("content-type", "application/json")
	res, err := json.Marshal(status)

	if err != nil {
		iseCatcher := catchHandler[505]
		if iseCatcher != nil {
			w.WriteHeader(int(status.Code))
			_, err = fmt.Fprint(w, iseCatcher(err.Error()))
			return
		}
		fmt.Println("Catcher [500] is not defined -> \n", err)
		return
	}

	w.WriteHeader(int(status.Code))
	_, err = fmt.Fprint(w, string(res))
	if err != nil {
		fmt.Println("[ERROR] \n", err)
	}

}

func main() {
	fmt.Println("\n\nStating up")

	fmt.Println("Connecting to MongoDB")
	db.MongoConnect("mongodb://mongo:27017/", "vobe-auth")

	catchHandler = CatchMap{
		404: catchers.NotFound,
		500: catchers.InternalServerError,
	}
	routeHandler = RouteMap{
		"signup": routes.Signup,
		"signin": routes.Login,
		"auth":   routes.Auth,
	}

	fmt.Println("Loaded routes:")
	for k := range routeHandler {
		fmt.Printf("--> /%s\n", k)
	}
	fmt.Println("Loaded catchers:")
	for k := range catchHandler {
		fmt.Printf("--> %d\n", k)
	}

	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
