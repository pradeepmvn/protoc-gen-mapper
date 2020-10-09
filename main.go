// Main package to start the generation process
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/pradeepmvn/protoc-gen-mapper/generator"
)

func main() {
	// new mapper gen
	mGen := generator.New()

	//parse request
	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Panic(err, "reading input")
	}

	if err := proto.Unmarshal(data, mGen.Request); err != nil {
		log.Panic(err, "parsing input proto")
	}
	mGen.Generate()

	marshalled, err := proto.Marshal(mGen.Response)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(marshalled)
}
