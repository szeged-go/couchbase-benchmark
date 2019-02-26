package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"

	"gopkg.in/couchbase/gocb.v1"
)

// UserHandler contain bucket
type UserHandler struct {
	bucket *gocb.Bucket
}

// NewUserHandler for init UserHandler
func NewUserHandler(bucket *gocb.Bucket) *UserHandler {
	return &UserHandler{
		bucket: bucket,
	}
}

// GetByKey user by ID
func (handler *UserHandler) GetByKey(key string) (User, error) {
	var user User
	_, error := handler.bucket.Get(key, &user)

	if error != nil {
		return User{}, error
	}

	return user, nil
}

// Create user from user struct
func (handler *UserHandler) Create(user User) (User, string, error) {

	user.CreatedAt = time.Now()

	id := xid.New().String()

	_, error := handler.bucket.Upsert(id, user, 0)

	if error != nil {
		return User{}, "", error
	}

	return user, id, nil

}

// Update user
func (handler *UserHandler) Update(user User, id string) (User, error) {

	if id == "" {
		return User{}, errors.New("Id not found")
	}

	var existUser User
	cas, err := handler.bucket.Get(id, &existUser)

	if err != nil {
		return User{}, err
	}

	if user.Email == "" {
		user.Email = existUser.Email
	}
	if user.Password == "" {
		user.Password = existUser.Password
	}
	if user.Fullname == "" {
		user.Fullname = existUser.Fullname
	}
	if user.Roles == nil {
		user.Roles = existUser.Roles[:]
	}
	if user.CreatedAt.IsZero() {
		user.CreatedAt = existUser.CreatedAt
	}

	user.UpdatedAt = time.Now()

	_, err = handler.bucket.Replace(id, user, cas, 0)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Delete by id
func (handler *UserHandler) Delete(id string) error {

	cas, err := handler.bucket.Get(id, nil)

	if err != nil {
		return err
	}

	_, err = handler.bucket.Remove(id, cas)

	if err != nil {
		return err
	}

	return nil

}

// Query by id
func (handler *UserHandler) Query(queryString string, args []interface{}) ([]User, error) {

	query := gocb.NewN1qlQuery(queryString)

	rows, err := handler.bucket.ExecuteN1qlQuery(query, args)

	if err != nil {
		fmt.Printf("Error when Query: %s\n", err)
		return nil, err
	}

	var users []User
	// var user User
	// for rows.Next(&user) {
	// fmt.Printf("User: %+v\n", user)
	// 	users = append(users, user)
	// }

	if err = rows.Close(); err != nil {
		fmt.Printf("Couldn't get all the rows: %s\n", err)
		return nil, err
	}

	return users, nil

}
