package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/qedus/osmpbf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var wg sync.WaitGroup

func main() {
	// command arg parse
	mongoURI := flag.String("uri", "mongodb://root@localhost:27017", "Mongo connection uri. (something look like \"mongodb://foo:bar@localhost:27017\")")
	pbfPath := flag.String("pbf", "default.osm.pbf", "OSM pbf file path.")
	dbName := flag.String("db", "map", "Mongodb database to import data to.")
	collectionName := flag.String("collection", "item", "Mongodb collection to import data to.")
	flag.Parse()

	pbfFile, err := os.Open(*pbfPath)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(*mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil) // assert server connection is okay
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Start dumping")

	collection := client.Database(*dbName).Collection(*collectionName)

	pbfDecoder := osmpbf.NewDecoder(pbfFile)
	pbfDecoder.Start(8)
	dataChan := make(chan bson.M, BUFFER_SIZE)
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go pbf2bson(pbfDecoder, dataChan)
	}
	go mongoWriter(collection, dataChan)
	wg.Wait()
	close(dataChan)
	fmt.Println("Done!")
}
