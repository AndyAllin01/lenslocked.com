package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound  = errors.New("models: resource not found")
	ErrInvalidID = errors.New("models: ID provided invalid")
)

const userPwPepper = "secret-random-string-this-project"

func NewUserService(connectionInfo string) (*UserService, error) {
	fmt.Println("CONN INFO ", connectionInfo)
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &UserService{
		db: db,
	}, nil

}

type UserService struct {
	db *gorm.DB
}

//ByID looks up user by provided ID
//1 - user, nil
//2 - nil, errNotFound
//3 - nil, otherError
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := db.First(&user).Error
	return &user, err
}

//ByEmail returns user based on email address
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//first will query using provided gorm DB and get first
//item returned and place into dst
//if nothing is found, ErrNotFound
func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

//Create creates provided user and backfills
//system fields
func (us *UserService) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper) // add pepper
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//	fmt.Println("UPDATE ", user.Password, hashedBytes)
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	//	fmt.Println(user)
	return us.db.Create(user).Error
	//	return nil
}

//Update the provided user with all data in the provided user object
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

//delete user with specified id
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

//Closes userService DB connection
func (us *UserService) Close() error {
	return us.db.Close()
}

//DestructiveReset drops users table and rebuilds it
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	fmt.Println("AUTOMIGRATE")
	return us.AutoMigrate()
}

//Attempt to automatically migrate the users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm: "-"`
	PasswordHash string `gorm:"not null"`
}
