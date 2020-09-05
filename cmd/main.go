package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"strings"
	"time"
	"vobe-auth/catchers"
	"vobe-auth/cmd/status"
	"vobe-auth/routes"
)

type RouteMap map[string]func(w *http.Request) status.Status
type CatchMap map[int16]func(msg string) status.Status

var routeHandler = make(RouteMap)
var catchHandler = make(CatchMap)
var client, ctx = mongoConnect("mongodb://mongo:27017")

var UserCollection = client.Database("vobe-auth").Collection("users")

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
	w.Header().Set("content-type", "application/json")
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

	catchHandler = CatchMap{
		404: catchers.NotFound,
		500: catchers.InternalServerError,
	}
	routeHandler = RouteMap{
		"signup": routes.Signup,
		"login":  routes.Login,
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

func mongoConnect(uri string) (*mongo.Client, *context.Context) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client, &ctx
}
