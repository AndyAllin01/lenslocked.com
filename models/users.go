package models

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"lenslocked.com/rand"

	"lenslocked.com/hash"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	ErrNotFound          = errors.New("models: resource not found")
	ErrIDInvalid         = errors.New("models: ID provided invalid")
	ErrInvalidEmail      = errors.New("models: incorrect email provided")
	ErrPasswordIncorrect = errors.New("models: incorrect password provided")
	ErrEmailRequired     = errors.New("Email address is required")
	ErrEmailInvalid      = errors.New("Email address is not valid")
	ErrEmailTaken        = errors.New("models: email address is already taken")
	ErrPasswordTooShort  = errors.New("models: password must be at least 8 characters")
	ErrPasswordRequired  = errors.New("models: password is required")
	ErrRememberTooShort  = errors.New("models: remember token must be at least 32 bytes")
	ErrRememberRequired  = errors.New("models: invlid remember token hassh")

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
	Delete(user uint) error
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

	hmac := hash.NewHMAC(hmacSecretKey)
	uv := newUserValidator(ug, hmac)
	return &userService{
		UserDB: uv,
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
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}
	return foundUser, nil
}

var _ UserDB = &UserValidator{}

func newUserValidator(udb UserDB, hmac hash.HMAC) *UserValidator {
	return &UserValidator{
		UserDB: udb,
		hmac:   hmac,
		//email matching regexp           bob99bob      @ email01     . com
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

type UserValidator struct {
	UserDB
	hmac       hash.HMAC
	emailRegex *regexp.Regexp
}

//ByEmail normalizes the email address before calling  ByEmail
//on the UserDB field
func (uv *UserValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

//Create creates provided user and backfills
//system fields
func (uv *UserValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.setRememberIfUnset,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

type userValFunc func(*User) error

//WOW!!!
//runUserValFunc runs every method for a particular type
func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

//Update will hash a remember token if provided
func (uv *UserValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.rememberMinBytes,
		uv.hmacRemember,
		uv.rememberHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvailable)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

//Delete a specified user
func (uv *UserValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThanZero)
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

//bcryptPassword hashes a users password with a predefined pepper
// and bcrypt if the password field is not ""
func (uv *UserValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + userPwPepper) // add pepper
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

func (uv *UserValidator) hmacRemember(user *User) error {
	if user.Remember == "" {
		return nil
	}
	user.RememberHash = uv.hmac.Hash(user.Remember)
	return nil
}

func (uv *UserValidator) setRememberIfUnset(user *User) error {
	if user.Remember != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	user.Remember = token
	return nil
}

func (uv *UserValidator) rememberMinBytes(user *User) error {
	if user.Remember == "" {
		return nil
	}
	n, err := rand.NBytes(user.Remember)
	if err != nil {
		return err
	}
	if n < 32 {
		return ErrRememberTooShort
	}
	return nil
}

func (uv *UserValidator) rememberHashRequired(user *User) error {
	if user.RememberHash == "" {
		return ErrRememberRequired
	}
	return nil
}

func (uv *UserValidator) idGreaterThanZero(user *User) error {
	if user.ID <= 0 {
		return ErrIDInvalid
	}
	return nil
}

//normalizeEmail removes spaces and sets email to lower case
func (uv *UserValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *UserValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *UserValidator) emailFormat(user *User) error {
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *UserValidator) emailIsAvailable(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrNotFound {
		return nil // email IS available
	}

	if err != nil {
		return err
	}
	//user found! If it has same id as this user, it is an update
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

//ByRemember hashes the remember token then calls ByRemember on
//the subsequent UserDB layer
func (uv *UserValidator) ByRemember(token string) (*User, error) {
	user := User{
		Remember: token,
	}
	if err := runUserValFuncs(&user, uv.hmacRemember); err != nil {
		return nil, err
	}
	return uv.UserDB.ByRemember(user.RememberHash)
}

func (uv *UserValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *UserValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *UserValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

var _ UserDB = &userGorm{}

func newUserGorm(connectionInfo string) (*userGorm, error) {
	fmt.Println("CONN INFO ", connectionInfo)
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	return &userGorm{
		db: db,
	}, nil
}

type userGorm struct {
	db *gorm.DB
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
func (ug *userGorm) ByRemember(rememberHash string) (*User, error) {
	var user User

	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//Create creates provided user and backfills
//system fields
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

//Update the provided user with all data in the provided user object
func (ug *userGorm) Update(user *User) error {
	return ug.db.Save(user).Error
}

//delete user with specified id
func (ug *userGorm) Delete(id uint) error {
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
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
