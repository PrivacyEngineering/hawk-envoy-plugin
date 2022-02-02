package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func rootHandler(_ http.ResponseWriter, request *http.Request) {
	var bodyBytes []byte

	if request.Body == nil {
		return
	}

	bodyBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Unable to read body. %v", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing body. %v", err)
		}
	}(request.Body)

	if len(bodyBytes) == 0 {
		return
	}

	log.Printf("BODY: %v", string(bodyBytes))

}
