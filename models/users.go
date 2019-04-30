package models

import (
	"errors"
	"fmt"

	"lenslocked.com/rand"

	"lenslocked.com/hash"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound        = errors.New("models: resource not found")
	ErrInvalidID       = errors.New("models: ID provided invalid")
	ErrInvalidEmail    = errors.New("models: incorrect email provided")
	ErrInvalidPassword = errors.New("models: incorrect password provided")
)

const userPwPepper = "secret-random-string-this-project"
const hmacSecretKey = "secret-hmac-key"

//User model including email, password and remember token used in cookie
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm: "-"`
	PasswordHash string `gorm:"not null"`
	Remember     string `gorm: "-"`
	RememberHash string `gorm:"not null;unique_index"`
}

//methods for querying for single users, interacting with users DB
//1 - user, nil
//2 - nil, errNotFound
//3 - nil, otherError
type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	//methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
	//used to close a DB connection
	Close() error

	//migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

//UserService is a set of methods to work with user model
type UserService interface {
	//Authenticate verifies provided email and password, returning user
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(connectionInfo string) (UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &userService{
		UserDB: &UserValidator{
			UserDB: ug,
		},
	}, nil
}

var _ UserService = &userService{}

type userService struct {
	UserDB
}

//Authenticate chceks password is correct for specified email address
func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+userPwPepper))

	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrInvalidPassword
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

var _ UserDB = &UserValidator{}

type UserValidator struct {
	UserDB
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	fmt.Println("CONN INFO ", connectionInfo)
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

//ByID looks up user by provided ID
//1 - user, nil
//2 - nil, errNotFound
//3 - nil, otherError
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := db.First(&user).Error
	return &user, err
}

//ByEmail returns user based on email address
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

//ByRemember looks up a user by remember token and returns that user
//this handles hashing of the token
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User

	rememberHash := ug.hmac.Hash(token)

	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Create creates provided user and backfills
//system fields
func (ug *userGorm) Create(user *User) error {
	pwBytes := []byte(user.Password + userPwPepper) // add pepper
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	//	fmt.Println("UPDATE ", user.Password, hashedBytes)
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
	//	return nil
}

//Update the provided user with all data in the provided user object
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

//delete user with specified id
func (ug *userGorm) Delete(id *User) error {
	//#######################################################
	//id confusion - is it uint or *User?
	/*	if id == 0 {
			return ErrInvalidID
		}
		user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error*/
	return nil
}

//Closes userService DB connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

//DestructiveReset drops users table and rebuilds it
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	fmt.Println("AUTOMIGRATE")
	return ug.AutoMigrate()
}

//Attempt to automatically migrate the users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
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
