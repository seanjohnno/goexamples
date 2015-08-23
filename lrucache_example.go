package main

import (
	"github.com/seanjohnno/memcache"
	"io/ioutil"
)

var (
	cache = memcache.CreateLRUCache(1024*4)
)

func main() {

}

type FileContent struct {
	content []byte
}

func (this *FileContent) Size() int {
	return len(content)
}

func LoadFile(fileLoc string) []byte, error {
	if cached, present := cache.Retrieve(); present {
		return cached.(*FileContent).content, nil
	} else {
		if fileContent, err := ioutil.ReadFile(fileLoc); err == nil {
			cache.Add(&FileContent{ content: fileContent})
		}
		return fileContent, err
	}
}