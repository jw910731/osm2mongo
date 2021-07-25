package main

import (
	"io"
	"log"

	"github.com/qedus/osmpbf"
	"go.mongodb.org/mongo-driver/bson"
)

// This function is designed to run in goroutine
// decoder param need to be setup and ready to decode before passed in
// error only report decode error other than io.EOF
func pbf2bson(decoder *osmpbf.Decoder, out chan bson.M) {
	defer wg.Done()
	for {
		raw, err := decoder.Decode()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			return
		}
		if ele, ok := raw.(*osmpbf.Node); ok {
			doc := bson.M{
				"_id":         ele.ID,
				"coordinates": bson.A{ele.Lon, ele.Lat},
				"type":        "Point",
				"time":        ele.Info.Timestamp,
			}
			for k, v := range ele.Tags {
				doc[k] = v
			}
			out <- doc
		}
	}
}
