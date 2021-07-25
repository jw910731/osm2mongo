package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const BUFFER_SIZE = 1024

func mongoWriter(coll *mongo.Collection, in chan bson.M) {
	var buffer [BUFFER_SIZE]interface{}
	size := 0
	for {
		var ok bool
		buffer[size], ok = <-in
		if ok {
			size++
			if size >= BUFFER_SIZE { // flush buffered data to database
				_, err := coll.InsertMany(context.TODO(), buffer[:])
				if err != nil {
					log.Println(err)
				}
				size = 0
			}
		} else {
			// final flush data into database
			_, err := coll.InsertMany(context.TODO(), buffer[:size])
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("Final flush success!\n")
			break
		}
	}
}
