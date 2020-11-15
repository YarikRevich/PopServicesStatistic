package Redis

import (
	"github.com/go-redis/redis/v8"
	"context"
	"log"
)

var ctx = context.Background()

type Redis struct{
	//Main structure for redis client.
	//It has connection field to save 
	//redis-client connection.

	connection *redis.Client
}


func (r *Redis)Connect(){
	//Connects to redis server.

	conn := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("Connected to Redis!")
	r.connection = conn
}


func (r Redis)StartDataProcessing(messageChan chan string){
	//Subscribes to important channel and starts to listen
	//to it. Pushes gotten string messages to `messageChan`.

	sub := r.connection.Subscribe(ctx, "test")
	for{
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil{
			panic(err)
		}
		messageChan <- msg.Payload
	}
}

func (r Redis)CheckIdDublesExist(StringToCheck string)bool{
	result := r.connection.LRange(ctx, "authed_users", 0, -1)
	list, err := result.Result()
	if err != nil{
		log.Fatalln(err)
	}
	for _, value := range(list){
		if value == StringToCheck{
			return true
		}
	}	 
	return false
}


func (r Redis)CreateAuthRecord(randomString string){
	r.connection.RPush(ctx, "authed_users", randomString)
}
