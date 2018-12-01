package models

import "sort"

type Value struct{
	Value1 int `json:"value1" orm:"column(value1);"`
	Value2 int `json:"value2" orm:"column(value2);"`
}

type StringValue struct{
	Value1 	string `json:"value1" orm:"column(value1);"`
	Value2	string `json:"value2" orm:"column(value2);"`
}

type DropDownListStr struct {
	Label  string  `json:"label"`
	Value  int  `json:"value"`
}

type DropDownListStr2 struct {
	Id  		string  `json:"id"`
	Name  		string  `json:"name"`
	Price		float64	`json:"price" orm:"column(price)"`
}

type FloatStruct struct {
	Value 	float64		`json:"value" orm:"column(value)"`
}


type ValueSet struct{
	Value1 int 		`json:"value1" orm:"column(value1)"`
	Value2 int 		`json:"value2" orm:"column(value2)"`
	Value3 int 		`json:"value3" orm:"column(value3)"`
	Value4 string 	`json:"value4" orm:"column(value4)"`
	Value5 string	`json:"value5" orm:"column(value5)"`
	Value6 int8		`json:"value6" orm:"column(value6)"`
	Value7 int8		`json:"value7" orm:"column(value7)"`
	Value8 int64	`json:"value8" orm:"column(value8)"`
	Value9 int64	`json:"value9" orm:"column(value9)"`
	Value10	string	`json:"value10" orm:"column(value10)"`
}


type uint32Slice []uint32

func (c uint32Slice) Len() int{
	return len(c)
}

func (c uint32Slice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c uint32Slice) Less(i, j int) bool{
	return c[i] < c[j]
}

func SortUint32(ii []uint32)  {
	sort.Sort(uint32Slice(ii))
}




