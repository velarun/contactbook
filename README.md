# contactbook

Pre-requistie:

Golang should be installed. 

Version >= 1.10.1

MySQL Should be installed.

Version >= 8.0.15

Setup:

set .env file:
Database name, machine, username, password and password(token).

Use make file:

> make dep
> make build
> make test

Sample API:

Create User:

curl -X POST \
  http://localhost:9000/user \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: c836b18a-1b06-4452-bf30-034194dd8903' \
  -H 'cache-control: no-cache' \
  -d '{ 
	"username": "test",
	"email": "test@testing.com",
	"password": "test"
}
'

Get token from response body.
Sample response:
{
    "account": {
        "ID": 2,
        "CreatedAt": "2019-07-15T07:18:10Z",
        "UpdatedAt": "2019-07-15T07:18:10Z",
        "DeletedAt": null,
        "username": "test",
        "email": "test@testing.com",
        "password": ""
    },
    "message": "Login success",
    "status": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyIjoyLCJVc2VybmFtZSI6InRlc3QifQ.KrirPAvJlVOCy8d67_BUntSkg4PI6zefwO2yetcPvlQ"
}

Login User:

curl -X POST \
  http://localhost:9000/login \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: d8de0edd-eb62-444b-bb37-d53d5a2d33f3' \
  -H 'cache-control: no-cache' \
  -d '{ 
	"username": "test",
	"password": "test"
}'

Delete User:

curl -X DELETE \
  http://localhost:9000/user \
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: c836b18a-1b06-4452-bf30-034194dd8903' \
  -H 'cache-control: no-cache' \
  -d '{ 
	"username": "test",
	"email": "test@testing.com",
	"password": "test"
}
'
Create Contact:

curl -X POST \
  http://localhost:9000/contact \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8b4a5e0e-6efe-4e6b-afc6-3cb5c6e791d8' \
  -H 'cache-control: no-cache' \
  -d '{
	"contact_name": "user11",
	"contact_email": "user11@hotmail.com",
	"phone_number": "+91 8897683200",
	"user_id": "test"
}'

Delete Contact:

curl -X DELETE \
  http://localhost:9000/contact \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8b4a5e0e-6efe-4e6b-afc6-3cb5c6e791d8' \
  -H 'cache-control: no-cache' \
  -d '{
	"contact_name": "user11",
	"contact_email": "user11@hotmail.com",
	"phone_number": "+91 8897683200",
	"user_id": "test"
}'

Update Contact:

curl -X PUT \
  http://localhost:9000/contact \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Content-Type: application/json' \
  -H 'Postman-Token: 8b4a5e0e-6efe-4e6b-afc6-3cb5c6e791d8' \
  -H 'cache-control: no-cache' \
  -d '{
	"contact_name": "user11",
	"contact_email": "user11@hotmail.com",
	"phone_number": "+91 8897683200",
	"user_id": "test"
}'

Search Contact:

curl -X GET \
  'http://localhost:9000/contact?contact_email=user1@yahoo.com' \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Postman-Token: 03c4e888-dcb0-4f17-9bcf-8bcc1cb05366' \
  -H 'cache-control: no-cache'

curl -X GET \
  'http://localhost:9000/contact?contact_name=user1' \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Postman-Token: 03c4e888-dcb0-4f17-9bcf-8bcc1cb05366' \
  -H 'cache-control: no-cache'

curl -X GET \
  'http://localhost:9000/contact?page=1' \
  -H "Authorization: Bearer xxxxxxxxxxxxxx"
  -H 'Postman-Token: 03c4e888-dcb0-4f17-9bcf-8bcc1cb05366' \
  -H 'cache-control: no-cache'
