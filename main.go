package main

import (
	"VK_Internship_Go/server"
	_ "github.com/lib/pq"
)

func main() {
	server.HttpServer()

}
