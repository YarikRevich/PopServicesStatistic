package main

import (
	"net/http"
	"github.com/my/repo/Services"
	"log"
	"github.com/my/repo/Mysql"
	"github.com/my/repo/Redis"
	"github.com/my/repo/Interfaces"
)

var (
	MysqlS = Mysql.Mysql{TableCon: "test"}
	RedisS = Redis.Redis{}
	MysqlDB = Interfaces.Database(&MysqlS)
	RedisDB = Interfaces.Database(&RedisS)
)

type ServerInter interface{
	Server() 
}

type Server struct{
	addr string
}

//type Handler struct{}

func (s Server)Server(){
	log.Println("Listening on " + s.addr + " ...")
	http.ListenAndServe(s.addr, nil)
}

func check(r *http.Request)bool{
	//Checks credentials

	result, err := r.Cookie("auth_id")
	if err != nil{
		return false
	}
	if RedisS.CheckIdDublesExist(result.Value){
		return true
	}
	return false
}

func auth(fr func(w http.ResponseWriter, r *http.Request), w http.ResponseWriter, r *http.Request){
	if !check(r){
		http.Redirect(w, r, "/auth", http.StatusMovedPermanently)
		return
	}
	fr(w, r)
}

func authMainPage(fr func(w http.ResponseWriter, r *http.Request, m Mysql.Mysql), w http.ResponseWriter, r *http.Request){
	if !check(r){
		http.Redirect(w, r, "/auth", http.StatusMovedPermanently)
		return
	}
	fr(w, r, MysqlS)
}

func loginPage(fc func(w http.ResponseWriter, r *http.Request, m Mysql.Mysql, redis Redis.Redis), w http.ResponseWriter, r *http.Request){
	if !check(r){
		fc(w, r, MysqlS, RedisS)
		return
	}
	http.Redirect(w, r, "/usability", 301)
}

func HandlerLoop(w http.ResponseWriter, r *http.Request){
	if r.URL.Path == "/"{
		http.Redirect(w, r, "/usability", 301)
	}
	if r.URL.Path == "/test" || r.URL.Path == "/test/"{
		auth(Services.SendCheckMessage, w, r)
	}
	if r.URL.Path == "/usability" || r.URL.Path == "/usability/"{
		authMainPage(Services.ShowMainPage, w, r)
	}
	if r.URL.Path == "/auth" || r.URL.Path == "/auth/"{
		loginPage(Services.ShowAuthPage, w, r)
	}
}

func MainRedirecter(){
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./Templates/static"))))
	http.HandleFunc("/", HandlerLoop)
}

func main(){
	RunDB()
	MainRedirecter()
	server := Server{
		addr: ":8000",
	}
	ServerInter := ServerInter(server)
	ServerInter.Server()
}

func RunDB(){
	//Runs all the databases, mean Mysql and Redis.

	var messageChan = make(chan string)  //Creates main channel for message exchange

	RedisDB.Connect(); MysqlDB.Connect()
	go RedisDB.StartDataProcessing(messageChan); go MysqlDB.StartDataProcessing(messageChan)
}