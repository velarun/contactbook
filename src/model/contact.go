package model

import (
	"github.com/jinzhu/gorm"
	"regexp"
	"log"
)

type Contact struct {
	ContactName string `json:"contact_name"`
	ContactEmail string `json:"contact_email"`
	PhoneNumber string `json:"phone_number"`
	UserId string `json:"user_id, omitempty"`
	gorm.Model
}

func (con *Contact) Validate(db *gorm.DB) (map[string]interface{}, bool) {

	resp := make(map[string] interface{})

	if con.ContactName == "" {
		resp["status"] = false
		resp["message"] = "Contact name is required."
		return resp, false
	} else if con.ContactName != "" {
		if len(con.ContactName) > 50 {
			resp["status"] = false
			resp["message"] = "Length of Contact name should be less 50 characters."
			return resp, false
		}
	}

	if con.ContactEmail == "" {
		resp["status"] = false
		resp["message"] = "Contact email is required."
		return resp, false
	} else if con.ContactEmail != "" {
		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
		if !re.MatchString(con.ContactEmail) {
			resp["status"] = false
			resp["message"] = "Contact email has incorrect pattern."
			return resp, false
		}
	}

	if con.PhoneNumber == "" {
		resp["status"] = false
		resp["message"] = "Phone number is required."
		return resp, false
	} else if con.PhoneNumber != "" {
		re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
		if !re.MatchString(con.PhoneNumber) {
			resp["status"] = false
			resp["message"] = "Phone number has incorrect pattern."
			return resp, false
		}
	}

	if con.UserId == "" {
		resp["status"] = false
		resp["message"] = "Username is required."
		return resp, false
	}

	user := &User{}
	db.Table("users").Where("username = ?", con.UserId).First(user)
	if user.Username == "" {
		resp["status"] = false
		resp["message"] = "User not found."
		return resp, false
	}

	resp["status"] = true
	resp["message"] = "All required fields(s) present"
	return resp, true
}

func (contact *Contact) Create(db *gorm.DB) (map[string]interface{}) {

	resp := make(map[string]interface{})
	con := &Contact{}
	db.Where("contact_email = ? AND user_id = ?", contact.ContactEmail, contact.UserId).First(con)
	if con.ContactEmail != "" {
		resp["status"] = false
		resp["message"] = "This contact already exists in user's contact list"
		return resp
	}

	db.Create(contact)

	log.Println("Create contact with Email:", contact.ContactEmail," for user ", contact.UserId, "Done.")
	resp["status"] = true
	resp["message"] = "Contact has been created."
	resp["contact"] = GetContact(contact.ID, db)
	return resp
}

func (contact *Contact) Delete(db *gorm.DB) (map[string]interface{}) {

	resp := make(map[string]interface{})
	con := &Contact{}
	db.Where("contact_email = ? AND user_id = ?", contact.ContactEmail, contact.UserId).First(con)
	if con.ContactEmail == "" {
		resp["status"] = false
		resp["message"] = "This contact not exists in user's contact list"
		return resp
	}

	db.Unscoped().Where("contact_email = ? AND user_id = ?", contact.ContactEmail, contact.UserId).Delete(con)

	log.Println("Delete contact with Email:", contact.ContactEmail," for user ", contact.UserId, "Done.")
	resp["status"] = true
	resp["message"] = "Contact has been deleted."
	return resp
}

func (contact *Contact) Update(db *gorm.DB) (map[string]interface{}) {

	resp := make(map[string]interface{})
	con := &Contact{}
	db.Where("contact_email = ? AND user_id = ?", contact.ContactEmail, contact.UserId).First(con)
	if con.ContactEmail == "" {
		resp["status"] = false
		resp["message"] = "This contact not exists in user's contact list"
		return resp
	}

	db.Table("contacts").Where("contact_email = ? AND user_id = ?", contact.ContactEmail, contact.UserId).Update(contact)

	log.Println("Update contact with Email:", contact.ContactEmail," for user ", contact.UserId, "Done.")
	resp["status"] = true
	resp["message"] = "Contact has been Updated."
	resp["contact"] = GetUpdatedContact(contact.UserId, contact.ContactEmail, db)
	return resp
}