package main

import (
	"log"
	"net/http"
	"strings"
	httperrors "vobe-auth/cmd/catchers"
	"vobe-auth/cmd/routes"
	status "vobe-auth/cmd/status"
)

type Handler *func(string) *status.Status
type RouteMap map[string]Handler
type CatchMap map[int16]Handler

var handlers = make(RouteMap)
var catchers = make(CatchMap)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	var route = handlers[strings.ToLower(r.URL.Path[1:])]
	if route == nil {
		(*catchers[404])("Not found")
		return
	}

	var res = *(*handlers[strings.ToLower(r.URL.Path[1:])])("yeet")

	println(res.Success, res.Content)
}

func RegisterCatcher(code int16, handler Handler) bool {
	if catchers[code] != nil {
		return false
	}
	catchers[code] = handler
	return true
}

func RegisterRoute(route *string, handler Handler) bool {
	strings.ToLower(*route)
	if handlers[*route] != nil {
		return false
	}
	handlers[*route] = handler
	println("")
	return true
}

func main() {
	var notfound = httperrors.NotFound
	catchers[404] = &notfound

	var auth = routes.Auth
	handlers["auth"] = &auth

	http.HandleFunc("/", HandleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
