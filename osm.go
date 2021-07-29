package main

import (
	"io"
	"log"
	"strings"

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
				if idx := strings.Index(k, ":"); idx > 0 {
					if _, ok := doc[k[:idx]]; ok { // entry exist
						if mp, ok := doc[k[:idx]].(map[string]string); ok { // entry is map
							mp[k[idx+1:]] = v
						} else { // entry conflict => conflicted entry must be in string type
							mp := make(map[string]string)
							mp[k[idx+1:]] = v
							// reserve original element in "org"
							mp["org"] = doc[k[:idx]].(string) // may panic though this won't likely happend
							doc[k[:idx]] = mp
						}
					} else { // create new entry of map
						mp := make(map[string]string)
						mp[k[idx+1:]] = v
						doc[k[:idx]] = mp
					}

				} else {
					if mp, ok := doc[k].(map[string]string); ok { // same entry with map type exist
						// reserve original element in "org"
						mp["org"] = v
					} else {
						doc[k] = v
					}

				}
			}
			out <- doc
		}
	}
}
