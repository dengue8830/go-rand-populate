package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
	Username             string `json:"username" bson:"username"`
}

var wg sync.WaitGroup

func populateChunk(id int) {
	defer wg.Done()
	seed := time.Now().UTC().UnixNano()
        nameGenerator := namegenerator.NewNameGenerator(seed)
	
	for i := 0; i < 100; i++ {
		user := &User{Username: nameGenerator.Generate()}
		mgm.Coll(user).Create(user)
	}
	fmt.Println("finished task", id)
}

func main() {
	fmt.Print("pid:", os.Getpid())
	start := time.Now()

	err := mgm.SetDefaultConfig(nil, "example-db", options.Client().ApplyURI("mongodb://root:example@localhost:27017"))
	if (err != nil) {
		fmt.Println("error connecting")
		panic("couldn't connect to db")
	}

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go populateChunk(i)
	}
	wg.Wait()
	fmt.Println("Finished in", time.Since(start))
}
