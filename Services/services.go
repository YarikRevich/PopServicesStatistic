package Services

import (
	"github.com/my/repo/Mysql"
	"github.com/my/repo/Redis"
	"log"
	"net/http"
	"encoding/json"
	"html/template"
	"math/rand"
	"strconv"
	"strings"
)

func SendCheckMessage(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"success": "ok",
	})
} 

func ShowMainPage(w http.ResponseWriter, r *http.Request, m Mysql.Mysql){
	if r.Method == "GET"{
		temp, err := template.ParseFiles("Templates/index.html")
		if err != nil{
			log.Fatalln(err)
		}
		temp.Execute(w, "")
		return
	}
	m.GetTheNewestData(w)
}

func CheckDubles(randomString string, redis Redis.Redis)bool{
	return redis.CheckIdDublesExist(randomString)
}

func CreateAuthRecord(randomString string, redis Redis.Redis){
	redis.CreateAuthRecord(randomString)
}

func CreateRandomId(redis Redis.Redis)string{
	randomNums := rand.Perm(10)
	randomString := make([]string, 10)
	for i := 0; i < 10; i++ {
		randomString[i] = strconv.Itoa(randomNums[i])
	}
	stringToReturn := strings.Join(randomString, "")
	if CheckDubles(strings.Join(randomString, ""), redis){
		return CreateRandomId(redis)
	}
	CreateAuthRecord(stringToReturn, redis)
	return stringToReturn
}

func LoginUser(w http.ResponseWriter, userId string){
	http.SetCookie(w, &http.Cookie{Name: "auth_id", Value: userId})
}

func ShowAuthPage(w http.ResponseWriter, r *http.Request, m Mysql.Mysql, redis Redis.Redis){
	r.ParseForm()
	form_storage := map[string]string{}
	for key, value := range r.Form{
		form_storage[key] = value[0]
	}
	if len(form_storage) != 0{
		if m.CheckAuth(form_storage){
			LoginUser(w, CreateRandomId(redis))
			http.Redirect(w, r, "/usability", http.StatusMovedPermanently)
		}
		return
	}
	temp, err := template.ParseFiles("Templates/auth.html")
	if err != nil{
		log.Fatalln(err)
	}
	temp.Execute(w, "")
}