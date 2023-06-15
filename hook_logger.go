package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
)

func main() {
	handler := http.HandlerFunc(handleRequest)
	http.Handle("/", handler)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()
	_ = r.ParseForm()
	form := r.PostForm

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("error when reading request body", err)
		return
	}

	go logIncoming(*r, body, form)

	_, _ = w.Write(nil)
	return
}

func logIncoming(r http.Request, body []byte, form url.Values) {
	requestID := rand.Intn(10000)
	log.SetOutput(os.Stdout)
	log.Printf("incoming message %s %s %d\n", r.Method, r.Host, requestID)

	closeFile := func(outputFile *os.File) {
		_ = outputFile.Close()
	}

	if len(body) > 0 {
		outputFileRawBody, err := os.Create(fmt.Sprintf("raw_body_incoming_%s_%s_%d.json", r.Method, r.Host, requestID))

		if err != nil {
			log.Println("error when creating output file", err)
			return
		}
		defer closeFile(outputFileRawBody)

		_, _ = outputFileRawBody.Write(body)
	}

	if len(form) > 0 {

		for v := range form {
			outputFileParams, err := os.Create(fmt.Sprintf("incoming_param_%s_%s_%s_%d.json", v, r.Method, r.Host, requestID))
			defer closeFile(outputFileParams)
			if err != nil {
				log.Println("error when creating output file", err)
				return
			}
			_, _ = outputFileParams.Write([]byte(form.Get(v)))

		}
	}

}
