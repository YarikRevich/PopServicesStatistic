package Mysql

import (
	"net/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"encoding/json"
)

type Mysql struct{
	//Main struct for databases, saves db connector
	//in `db` field and the name of the table to create
	//in `TableCon` field. 

	db *sql.DB
	TableCon string
}

func (g *Mysql)Connect(){
	//Makes connection to mysql database and sets conn to db
	//variable in struct `DB`.

	db, err := sql.Open("mysql", "yaroslav:yariksun4002@/massage")
	if err != nil{
		panic(err)
	}
	g.db = db
	if g.checkTableExists(g.TableCon){
		log.Println("Connected to Mysql!")
		return
	}
	g.createTable(g.TableCon, []string{"service_name varchar(100)", "visites_num int"})
	log.Println("Connected to Mysql!")
}

func (g Mysql)createTable(tableName string, columns []string){
	//Creates passed table in `TableCon` field.

	//_, err := g.db.Exec("create table " + g.TableCon + " (service_name varchar(100), visites_num int)")
	_, err := g.db.Exec("create table " + tableName + " (id int primary key)")
	if err != nil{
		log.Fatalln(err)
	}
	for _, value := range(columns){
		_, err := g.db.Exec("alter table " + tableName + " add column " + value)
		if err != nil{
			log.Fatalln(err)
		}
	}
}

func (g Mysql)checkTableExists(tableName string)bool{
	//Checks whether important table exists
	//If it does returns `true`, but if doesn't does `false`.

	result, err := g.db.Query("show tables like '" + tableName + "'")
	if err != nil{
		log.Fatalln(err)
	}
	var res []uint8
	for result.Next(){
		err := result.Scan(&res)
		if err != nil{
			log.Fatalln(err)
		}
	}
	if len(res) != 0{
		return true
	}
	return false
}

func (g Mysql)checkSuchServiceExists(serviceName string)bool{
	//Checks whether record for important service exists
	//If it does returns `true`, but if it doesn't does `false`.

	result, err := g.db.Query("select service_name from " + g.TableCon + " where service_name = ?", serviceName)
	if err != nil{
		log.Fatalln(err)
	}
	var res []uint8
	for result.Next(){
		result.Scan(&res)
	}
	if len(res) != 0{
		return true
	}
	return false
}

func (g Mysql)createServiceRecord(serviceName string){
	//Creates start records for important service
	//which name passed in `serviceName` param.

	_, err := g.db.Exec("insert into " + g.TableCon + " set service_name = ?, visites_num = 0", serviceName)
	if err != nil{
		log.Fatalln(err)
	}
}

func (g Mysql)updateServiceRecord(serviceName string){
	//Updates the number of visites for important 
	//service which name passed in `serviceName` param. 

	_, err := g.db.Exec("update " + g.TableCon + " set visites_num = visites_num + 1 where service_name = ?", serviceName)
	if err != nil{
		log.Fatalln(err)
	}
}

func (g Mysql)StartDataProcessing(messageChan chan string){
	//Sets gotten data from redis

	for{
		msg := <-messageChan
		if len(msg) != 0{
			if !g.checkSuchServiceExists(msg){
				g.createServiceRecord(msg)
			}
			g.updateServiceRecord(msg)
		}
	}
}

func (g Mysql)GetTheNewestData(w http.ResponseWriter){
	result, err := g.db.Query("select * from " + g.TableCon)
	if err != nil{
		log.Fatalln(err)
	}
	staticInfo := map[string]int{}
	for result.Next(){
		var (
			serviceName string
			visitesNum int
		)
		result.Scan(&serviceName, &visitesNum)
		staticInfo[serviceName] = visitesNum
	}
	json.NewEncoder(w).Encode(staticInfo)
	
}

func (g Mysql)CheckAuth(form map[string]string)bool{
	if !g.checkTableExists("users"){
		g.createTable("users", []string{"email char(100)", "pass char(100)"})
	}
	result, err := g.db.Query("select * from users;")
	if err != nil{
		log.Fatalln(err)
	}
	var(
		id []uint8
		email []uint8
		pass []uint8
	)
	for result.Next(){
		result.Scan(&id, &email, &pass)
		if form["email"] == string(email) && form["pass"] == string(pass){
			return true
		}
	}
	return false
}
