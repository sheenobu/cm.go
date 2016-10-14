package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"
	"net/http"

	"strconv"

	"golang.org/x/net/context"

	"github.com/bmizerany/pat"
	"github.com/codegangsta/negroni"
)

var mux = http.NewServeMux()
var router = pat.New()

// initHTTP starts up the various routers for hosting our static web interface
// and the API to our datamodel.
func initHTTP() {
	n := negroni.New()
	n.UseHandler(mux)

	gopath := os.Getenv("GOPATH")
	www := "/src/github.com/sheenobu/cm.go/demos/cm-webframeworks/www"

	mux.Handle("/api/", router)
	mux.Handle("/", http.FileServer(http.Dir(gopath+www)))

	router.Get("/api/frameworks/:id", http.HandlerFunc(frameworkGet))
	router.Del("/api/frameworks/:id", http.HandlerFunc(frameworkDelete))

	router.Get("/api/frameworks", http.HandlerFunc(frameworksRead))
	router.Post("/api/frameworks", http.HandlerFunc(frameworksPost))

	router.Get("/api/", http.HandlerFunc(infoRoute))

	fmt.Printf("listening on port 8888\n")
	err := http.ListenAndServe(":8888", n)
	if err != nil {
		panic(err)
	}
}

// infoRoute simply returns OK
func infoRoute(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "OK",
	}

	body, _ := json.Marshal(data)

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write(body)
}

// FrameworksPayload is a wrapper object for passing along pagination information
// as well as the framework list.
type FrameworksPayload struct {
	Frameworks  []Framework
	PageCount   int
	CurrentPage int
}

// frameworkRead is the GET /api/framework operation, sending out JSON
// and supporting pagination
func frameworksRead(w http.ResponseWriter, r *http.Request) {

	perPage := r.URL.Query().Get("perPage")
	if perPage == "" {
		perPage = "3"
	}

	perPageInt, err := strconv.Atoi(perPage)
	if err != nil {
		w.Header()["Content-Type"] = []string{"text/plain"}
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}

	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		w.Header()["Content-Type"] = []string{"text/plain"}
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}

	frameworks := FrameworksPayload{
		Frameworks:  make([]Framework, 0),
		PageCount:   0,
		CurrentPage: pageInt,
	}

	pager, err := Frameworks.Page(context.Background(), perPageInt)

	if err != nil {
		w.Header()["Content-Type"] = []string{"text/plain"}
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}

	frameworks.PageCount = pager.PageCount()

	for ok := true; ok && pageInt != pager.CurrentPage(); ok = pager.Next() {

	}

	frameworks.CurrentPage = pager.CurrentPage()

	err = pager.Apply(&frameworks.Frameworks)

	if err != nil {
		w.Header()["Content-Type"] = []string{"text/plain"}
		w.Write([]byte(fmt.Sprintf("Error: %s\n", err)))
		return
	}

	var body []byte

	body, _ = json.Marshal(frameworks)

	w.Header()["Content-Type"] = []string{"application/json"}
	w.Write(body)
}

// frameworkGet gets the framework object as JSON given the ID path parameter
func frameworkGet(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get(":id")

	ctx := context.Background()

	var frameworks []Framework

	err := Frameworks.Filter(Frameworks.ID.Eq(id)).Single(ctx, &frameworks)

	w.Header()["Content-Type"] = []string{"application/json"}

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(
			`{ "error": "%s" }`, err)))
		return
	}

	var body []byte
	w.Header()["Content-Type"] = []string{"application/json"}

	if len(frameworks) == 0 {
		w.WriteHeader(404)
		body = []byte(`{"error": "not found"}`)
	} else {
		body, _ = json.Marshal(frameworks[0])
	}

	w.Write(body)
}

// frameworksPost maps the POST operation and creates a framework object
// given the input JSON
func frameworksPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	body, _ := ioutil.ReadAll(r.Body)

	var framework Framework

	err := json.Unmarshal(body, &framework)

	w.Header()["Content-Type"] = []string{"application/json"}

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(
			`{ "error": "%s" }`, err)))
		return
	}

	err = Frameworks.Insert(ctx, framework)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(
			`{ "error": "%s" }`, err)))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(fmt.Sprintf(
		`{ }`)))

}

// frameworkDelete removes the framework given the ID as a path parameter
func frameworkDelete(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	id := r.URL.Query().Get(":id")

	err := Frameworks.Filter(Frameworks.ID.Eq(id)).Delete(ctx)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(
			`{ "error": "%s" }`, err)))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(
		`{ }`)))
}
