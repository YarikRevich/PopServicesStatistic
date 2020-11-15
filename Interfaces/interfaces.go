package Interfaces

type Database interface{
	//It is a main interface to work
	//with mysql db and redis key-value db
	//Has `Connect` method and `StartDataProcessing`
	//method with chan string param. 

	Connect()
	StartDataProcessing(chan string)
}