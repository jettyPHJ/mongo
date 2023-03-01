package main

import (
	_ "ToMongoAndRedis/implement/dataToMongo"
	"time"
)

func main() {
	for {
		time.Sleep(time.Hour)
	}
}
