package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type UserData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	FirstName         string             `json:"first_name"`
	LastName          string             `json:"last_name"`
	Email             string             `json:"email"`
	EncryptedPassWord string             `json:"-"`
}

func NewUser(userData UserData) (*User, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         userData.FirstName,
		LastName:          userData.LastName,
		Email:             userData.Email,
		EncryptedPassWord: string(password),
	}, nil
}

func IsValidPassword(encPwd, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encPwd), []byte(pwd)) == nil
}
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}
func (s *UserData) Validate() map[string]string {
	errors := map[string]string{}
	if len(s.LastName) < 2 {
		errors["lastName"] = "Last name should be at least 2 characters!"
	}
	if len(s.FirstName) < 2 {
		errors["firstName"] = "First name should be at least 2 characters!"
	}
	if !isEmailValid(s.Email) {
		errors["lastName"] = "Email address not valid!"
	}
	if len(s.Password) < 8 {
		errors["password"] = "password should be at least 8 characters!"
	}
	return errors
}
func (s *UserData) ValidateUpdate() bson.M {
	update := bson.M{}
	if len(s.LastName) > 0 {
		update["lastName"] = s.LastName
	}
	if len(s.FirstName) > 0 {
		update["firstName"] = s.FirstName
	}
	return update
}
