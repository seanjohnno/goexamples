package main

import (
	"github.com/seanjohnno/memcache"
	"io/ioutil"
	"fmt"
	"net/http"
)

var (
	Cache = memcache.CreateLRUCache(1024*4)
	Resource = "https://raw.githubusercontent.com/seanjohnno/goexamples/master/helloworld.txt"
)

func main() {
	fileContent, _ := GetHttpData(Resource)
	fmt.Println("Content:", string(fileContent))

	fileContent, _ = GetHttpData(Resource)
	fmt.Println("Content:", string(fileContent))
}

type HttpContent struct {
	Content []byte
}

func (this *HttpContent) Size() int {
	return len(this.Content)
}

func GetHttpData(URI string) ([]byte, error) {
	if cached, present := Cache.Get(URI); present {
		fmt.Println("Found in cache")
		return cached.(*HttpContent).Content, nil
	} else {
		fmt.Println("Not found in cache, making network request")

		// No error handling here to make example shorter
		resp, err := http.Get(URI)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)		
		Cache.Add(URI, &HttpContent{ Content: body})
		return body, err
	}
}