package model

func GetContact(Id uint) *Contact {

	contact := &Contact{}
	GetConn().First(contact, Id)
	return contact
}

func GetUserContact(userId interface{}) *[]Contact {

	var contacts []Contact

	GetConn().Table("contacts").Where("user_id = ?", userId).Find(&contacts)
	return &contacts
}

func GetUser(userId interface{}) *User {

	var user User
	
	GetConn().Table("users").Where("id = ?", userId).Find(&user)
	return &user
}