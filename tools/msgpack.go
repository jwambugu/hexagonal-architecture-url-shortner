package main

import (
	"bytes"
	"fmt"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	"github.com/vmihailenco/msgpack"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := fmt.Sprintf("http://localhost:%d", 8000)

	redirect := shortener.Redirect{}
	redirect.URL = "https://github.com/jwambugu"

	requestBody, err := msgpack.Marshal(&redirect)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(url, "application/x-msgpack", bytes.NewBuffer(requestBody))

	if err != nil {
		log.Fatal(err)
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	if err := msgpack.Unmarshal(responseBody, &redirect); err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", redirect)
}
