package models

import (
	"errors"
	"time"

	s "github.com/me/todo-go-server/src/shared"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	ID           uint `gorm:"primary_key"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
	FirstName    string     `gorm:"column:firstname"`
	LastName     string     `gorm:"column:lastname"`
	Email        string     `gorm:"column:email;unique_index;not null" json:"email"`
	Bio          string     `gorm:"column:bio;size:1024"`
	Image        *string    `gorm:"column:image"`
	PasswordHash string     `gorm:"column:password;not null" json:"password"`
	Todos        []TodoModel
}

// func AutoMigrateUsers() {
// 	db := s.GetDB()
// 	db.AutoMigrate(&UserModel{})
// }

func (u *UserModel) SetPassword(password string) error {
	if len(password) == 0 {
		return errors.New("password should not be empty!")
	}
	bytePassword := []byte(password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.PasswordHash = string(passwordHash)
	return nil
}

// Database will only save the hashed string, you should check it by util function.
// 	if err := serModel.checkPassword("password0"); err != nil { password error }
func (u *UserModel) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.PasswordHash)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

func FindOneUser(condition interface{}) (UserModel, error) {
	db := s.GetDB()
	var model UserModel
	err := db.Where(condition).First(&model).Error
	return model, err
}

// You could input an UserModel which will be saved in database returning with error info
// 	if err := SaveOne(&userModel); err != nil { ... }
func SaveOne(data interface{}) error {
	db := s.GetDB()
	err := db.Save(data).Error
	return err
}

// You could update properties of an UserModel to database returning with error info.
//  err := db.Model(userModel).Update(UserModel{Email: "wangzitian0"}).Error
func (model *UserModel) Update(data interface{}) error {
	db := s.GetDB()
	err := db.Model(model).Update(data).Error
	return err
}
