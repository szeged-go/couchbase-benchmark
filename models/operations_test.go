package models_test

import (
	"testing"

	"github.com/borosr/couchbase-banchmark/models"
	"gopkg.in/couchbase/gocb.v1"
)

const databaseURL = "couchbase://localhost"

func TestUserHandler(t *testing.T) {
	cluster, _ := gocb.Connect(databaseURL)
	cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "banchmark",
		Password: "banchmark",
	})
	bucket, error := cluster.OpenBucket("banchmark", "")

	if error != nil {
		t.Error("Error when opening bucket:", error)
		return
	}

	bucket.Manager("", "").CreatePrimaryIndex("", true, false)

	userHandler := models.NewUserHandler(bucket)

	savedUser, id, error := userHandler.Create(models.User{
		Email:    "test_user2@spam4.me",
		Fullname: "Test User2",
		Password: "123456",
		Roles:    []string{"ADMIN"},
	})

	if error != nil {
		t.Error("User save error:", error)
		return
	}
	t.Log("Saved user:", savedUser)

	existUser, error := userHandler.GetByKey(id)
	if error != nil {
		t.Error(error)
		return
	}
	t.Log("Exist user:", existUser)
}
