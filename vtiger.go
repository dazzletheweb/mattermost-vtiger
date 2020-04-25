package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Person struct {
	Firstname string
	Lastname string
	Email string
	Phone string
	Mobile string
	Id string
}

type vTigerAccess struct {
	// The label in the UI is "User Name" but the database key is "username"...
	UserName string
	// See previous comment.
	AccessKey string
	BaseUrl string
}

type vTigerResult struct {
	Success bool
	Result []Person
}

func search(config vTigerAccess, keyword string) string {
	keyword = strings.ReplaceAll(keyword, "'", "")
	keyword = strings.ReplaceAll(keyword, "\\", "")
	if len(keyword) < 3 {
		return "The search keyword must be at least 3 characters long."
	}

	var client http.Client
	req, _ := http.NewRequest("GET", config.BaseUrl + "/webservice.php", nil)
	q := req.URL.Query()
	q.Add("operation", "getchallenge")
	q.Add("username", config.UserName)
	req.URL.RawQuery = q.Encode()
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var dat map[string]interface{}
	json.Unmarshal(body, &dat)
	result := dat["result"].(map[string]interface{})
	var token string = result["token"].(string)
	h := md5.New()
	io.WriteString(h, token)
	io.WriteString(h, config.AccessKey)
	var hash string = fmt.Sprintf("%x", h.Sum(nil))

	resp, _ = http.PostForm(config.BaseUrl + "/webservice.php",
		url.Values{"operation": {"login"}, "username": {config.UserName}, "accessKey": {hash}})
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &dat)
	result = dat["result"].(map[string]interface{})
	var sessionname string = result["sessionName"].(string)

	query := `
	  SELECT firstname, lastname, email, phone, mobile
	  FROM Contacts
	  WHERE firstname LIKE '%%%s%%' OR lastname LIKE '%%%s%%';
`
	query = fmt.Sprintf(query, keyword, keyword)
	req, _ = http.NewRequest("GET", config.BaseUrl + "/webservice.php", nil)
	q = req.URL.Query()
	q.Add("sessionName", sessionname)
	q.Add("operation", "query")
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()
	resp, _ = client.Do(req)

	defer resp.Body.Close()
	persons := vTigerResult{}
	body, _ = ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &persons)

	var table string = "|Name|Email|Phone|Mobile|\n"
	table += "|:-|:-|:-|:-|\n"
	for _, v := range persons.Result {
		parts := strings.FieldsFunc(v.Id, func(r rune) bool { return r == 'x' })
		name := fmt.Sprintf("[%s %s](%s)", v.Firstname, v.Lastname, config.BaseUrl + "/index.php?module=Contacts&view=Detail&record=" + parts[1])
		table += fmt.Sprintf("|%s|%s|%s|%s|\n", name, v.Email, v.Phone, v.Mobile)
	}
	return table
}
