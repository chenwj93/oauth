package main

import (
	"oauth/controllers"
	"testing"
	"oauth/models"
	"strconv"
	"time"
)

func TestSchedule(t *testing.T)  {

	for i := 1; i < 10; i++ {
		index := &models.Index{strconv.Itoa(i), strconv.Itoa(i), time.Now().Add(time.Second * time.Duration(i))}
		controllers.ScheduleList.PushBack(index)
	}

	go controllers.ScheduleLim(time.Second * 10)
	c := make(chan bool)
	<-c
}