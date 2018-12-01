package models

import (
	"time"
)

type Index struct {
	Abs 	string
	Token 	string
	Expire 	time.Time
}

type IndexSlice []*Index

func (c IndexSlice) Len() int{
	return len(c)
}

func (c IndexSlice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c IndexSlice) Less(i, j int) bool{
	return c[i].Expire.Before(c[j].Expire)
}
