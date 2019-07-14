package controller

import (
	"net/http"
	"model"
	"encoding/json"
	"log"
	"fmt"
	"net/url"
	"strconv"
	"os"

	pagination "github.com/gemcook/pagination-go"
)

var userid interface{}

type SearchCondition struct {
	Email *string
	Name *string
}

type contactsRepository struct {
	email string
	name string
}

type contactFetcher struct {
	repo *contactsRepository
}

func (fr *contactsRepository) GetContact(orders []*pagination.Order) []model.Contact {
	result := make([]model.Contact, 0)
	user := model.GetUser(userid)
	contacts := model.GetUserContact(user.Username)
	if fr.email != "" {
		log.Println("Search Pattern is Email for user ", user.Username)
		for _, f := range *contacts {
			if fr.email == f.ContactEmail {
				result = append(result, f)
			}
		}
	} else if fr.name != "" {
		log.Println("Search Pattern is Name for user ", user.Username)
		for _, f := range *contacts {
			if fr.name == f.ContactName {
				result = append(result, f)
			}
		}
	} else {
		log.Println("Search Pattern is nothing simply page display for user ", user.Username)
		for _, f := range *contacts {
			result = append(result, f)
		}
	}
		
	return result
}

func (ff *contactFetcher) ConditionCheck(cond *SearchCondition) {
	if cond.Email != nil {
		ff.repo.email = *cond.Email
	}
	if cond.Name != nil {
		ff.repo.name = *cond.Name
	}
}

func (ff *contactFetcher) Count(cond interface{}) (int, error) {
	if cond != nil {
		ff.ConditionCheck(cond.(*SearchCondition))
	}
	orders := make([]*pagination.Order, 0, 0)
	contacts := ff.repo.GetContact(orders)
	return len(contacts), nil
}

func (ff *contactFetcher) FetchPage(cond interface{}, input *pagination.PageFetchInput, result *pagination.PageFetchResult) error {
	if cond != nil {
		ff.ConditionCheck(cond.(*SearchCondition))
	}
	contacts := ff.repo.GetContact(input.Orders)
	var toIndex int
	toIndex = input.Offset + input.Limit
	if toIndex > len(contacts) {
		toIndex = len(contacts)
	}
	for _, contact := range contacts[input.Offset:toIndex] {
		*result = append(*result, contact)
	}
	return nil
}

func newSearchCondition(email, name string) *SearchCondition {
	return &SearchCondition{
		Email:  &email,
		Name: &name,
	}
}

func parseSearchIndex(queryStr string) *SearchCondition {
	u, err := url.Parse(queryStr)
	if err != nil {
		log.Println(err)
		email := ""
		name := ""
		return newSearchCondition(email, name)
	}
	query := u.Query()

	if s := query.Get("contact_email"); s != "" {
		name := ""
		return newSearchCondition(s, name)
	} else if s := query.Get("contact_name"); s != "" {
		email := ""
		return newSearchCondition(email, s)
	} 
	
	email := ""
	name := ""
	return newSearchCondition(email, name)
}

func newContactFetcher() *contactFetcher {
	return &contactFetcher{
		repo: &contactsRepository{},
	}
}

func newContactsRepository() *contactsRepository {
	return &contactsRepository{
		email:  "",
		name: "",
	}
}

func CreateContact(w http.ResponseWriter, r *http.Request) {

	contact := &model.Contact{}
	json.NewDecoder(r.Body).Decode(contact)

	if resp, ok := contact.Validate(); !ok {
		log.Println("Create Contact: Validation Failed")
		Respond(w, resp)
	} else {
		resp = contact.Create()
		Respond(w, resp)
	}
}

func DeleteContact(w http.ResponseWriter, r *http.Request) {

	contact := &model.Contact{}
	json.NewDecoder(r.Body).Decode(contact)

	if resp, ok := contact.Validate(); !ok {
		log.Println("Delete Contact: Validation Failed")
		Respond(w, resp)
	}else {
		resp = contact.Delete()
		Respond(w, resp)
	}
}

func UpdateContact(w http.ResponseWriter, r *http.Request) {

	contact := &model.Contact{}
	json.NewDecoder(r.Body).Decode(contact)

	if resp, ok := contact.Validate(); !ok {
		log.Println("Update Contact: Validation Failed")
		Respond(w, resp)
	}else {
		resp = contact.Update()
		Respond(w, resp)
	}
}

func SearchContact(w http.ResponseWriter, r *http.Request) {
	userid = r.Context().Value("user")
	p := pagination.ParseQuery(r.URL.RequestURI())
	cond := parseSearchIndex(r.URL.RequestURI())
	fetcher := newContactFetcher()
	pg := os.Getenv("pagelimit")
	pageLimit, _ := strconv.Atoi(pg)

	totalCount, totalPages, res, err := pagination.Fetch(fetcher, &pagination.Setting{
		Limit:  pageLimit,
		Page:   p.Page,
		Cond:   cond,
	})

	if err != nil {
		log.Println("Error During Search Contact:", err)
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(400)
		fmt.Fprintf(w, "something wrong: %v", err)
		return
	}

	log.Println("Pagination of Search Content, Total Pages:", strconv.Itoa(totalPages), ", Total Count of Data:", strconv.Itoa(totalCount), ", Limit per Page:", pageLimit)
	
	if len(res.Pages["active"]) == 0 {
		w.Header().Set("Content-Type", "text/html; charset=utf8")
		w.WriteHeader(404)
		fmt.Fprintf(w, "Content Not Found")
	} else {
		w.Header().Set("X-Total-Count", strconv.Itoa(totalCount))
		w.Header().Set("X-Total-Pages", strconv.Itoa(totalPages))
		w.Header().Set("Access-Control-Expose-Headers", "X-Total-Count,X-Total-Pages")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		resp, _ := json.Marshal(res.Pages["active"])
		w.Write(resp)
	}
}
