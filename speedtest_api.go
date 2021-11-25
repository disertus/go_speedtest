package main

import (
	_ "fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Schedule struct{
	MacAddress string
	ScheduledTime time.Time
}


// todo: add mac-address parameter to the /schedule endpoint
// todo: create a dummy schedule for the local client to fetch
// todo: lay out some logic for future schedules

func main() {
	resp := gin.Default()
	var sched Schedule

	resp.GET("/schedule", func(c *gin.Context) {
		c.JSON(200, gin.H{
			sched.MacAddress: sched.ScheduledTime,
		})
	})
	resp.Run()
}