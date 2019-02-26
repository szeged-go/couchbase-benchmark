package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/borosr/couchbase-banchmark/models"

	"gopkg.in/couchbase/gocb.v1"
)

const databaseURL = "couchbase://localhost"

func main() {
	cluster, _ := gocb.Connect(databaseURL)
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "banchmark",
		Password: "banchmark",
	})
	bucket, err := cluster.OpenBucket("banchmark", "")

	if err != nil {
		fmt.Println("Error when opening bucket:", err)
		return
	}

	bucket.Manager("", "").CreatePrimaryIndex("created_at", true, false)

	var innerWg = &sync.WaitGroup{}
	for j := 0; j < 100; j++ {
		// for i := 0; i < 100; i++ {
		innerWg.Add(1)
		go func(i int) {
			action(bucket, i)
			innerWg.Done()
		}(j)
		// }
	}
	innerWg.Wait()

}

func action(bucket *gocb.Bucket, index int) {

	userHandler := models.NewUserHandler(bucket)
	t := time.Now()
	_, _, err := userHandler.Create(models.User{
		Email:    "test_user2@spam4.me",
		Fullname: "Test User2",
		Password: "123456",
		Roles:    []string{"ADMIN"},
	})

	if err != nil {
		fmt.Println("User Save error:", err)
		return
	}

	fmt.Println(time.Since(t), " index: ", index)
}
