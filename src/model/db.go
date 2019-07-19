package model

import (
	"github.com/jinzhu/gorm"
)

func GetContact(Id uint, db *gorm.DB) *Contact {

	contact := &Contact{}
	db.First(contact, Id)
	return contact
}

func GetUserContact(userId interface{}, db *gorm.DB) *[]Contact {

	var contacts []Contact

	db.Table("contacts").Where("user_id = ?", userId).Find(&contacts)
	return &contacts
}

func GetUpdatedContact(userId interface{}, email interface{}, db *gorm.DB) *Contact {

	contact := &Contact{}

	db.Table("contacts").Where("contact_email = ? AND user_id = ?", email, userId).Find(contact)
	return contact
}

func GetUser(userId interface{}, db *gorm.DB) *User {

	var user User
	
	db.Table("users").Where("id = ?", userId).Find(&user)
	return &user
}

func GetUserByUsername(username interface{}, db *gorm.DB) *User {

	var user User
	
	db.Table("users").Where("username = ?", username).Find(&user)
	return &user
}