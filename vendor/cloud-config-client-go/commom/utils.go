package common

import (
	"encoding/json"
	"strings"
	"strconv"
	"time"
	"bytes"
	"log"
)

const (
	ONE int = 1 + iota
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
)

const (
	ONE_I int8 = 1 + iota
	TWO_I
	THREE_I
	FOUR_I
	FIVE_I
	SIX_I
	SEVEN_I
	EIGHT_I
	NINE_I
	TEN_I
)

const (
	ONE_B int64 = 1 + iota
	TWO_B
	THREE_B
	FOUR_B
	FIVE_B
	SIX_B
	SEVEN_B
	EIGHT_B
	NINE_B
	TEN_B
)


type Base struct {
	Page	string		`json:"page"`
	Rows	string		`json:"rows"`
	OrderBy string		`json:"orderBy"`
}


//create by cwj on 2017-09-16
func ParseInt(val interface{}) (int) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseInt(v, 10, 64)
			return int(i)
		case int64:
			return int(v)
		case int8:
			return int(v)
		case int32:
			return int(v)
		case int:
			return v
		case float64:
			return int(v)
		case float32:
			return int(v)
		case time.Time:
			return int(v.UnixNano())
		}
		return 0
	}
	return 0
}

//create by cwj on 2017-09-16
func ParseInt64(val interface{}) (int64) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseInt(v, 10, 64)
			return i
		case int64:
			return v
		case int8:
			return int64(v)
		case int32:
			return int64(v)
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case float32:
			return int64(v)
		case time.Time:
			return v.UnixNano()
		}
		return 0
	}
	return 0
}

//create by cwj on 2017-11-25
func ParseFloat64(val interface{}) (float64) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return 1
			} else {
				return 0
			}
		case string:
			i, _ := strconv.ParseFloat(v, 64)
			return i
		case int64:
			return float64(v)
		case int8:
			return float64(v)
		case int32:
			return float64(v)
		case int:
			return float64(v)
		case float64:
			return v
		case float32:
			return float64(v)
		case time.Time:
			return float64(v.UnixNano())
		}
		return 0
	}
	return 0
}

const EMPTY_STRING  = ""
const TIME_FORMAT_1  = "2006-01-02 15:04:05"
func ParseString(val interface{}) (string) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			if v {
				return "true"
			} else {
				return "false"
			}
		case string:
			return v
		case int64:
			return strconv.FormatInt(v, 10)
		case int8:
			return strconv.Itoa(int(v))
		case int32:
			return strconv.Itoa(int(v))
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64 )
		case float32:
			return strconv.FormatFloat(float64(v), 'f', -1, 32 )
		case time.Time:
			return v.Format(TIME_FORMAT_1)
		}
		return EMPTY_STRING
	}
	return EMPTY_STRING
}

//create by cwj on 2017-10-17
// parse string from array
func ParseStringFromArray(a []string, character string, empty string) string {
	var buff bytes.Buffer
	for i, e := range a {
		if e == EMPTY_STRING {
			e = empty
		}
		if i < len(a) - 1{
			buff.WriteString(ParseString(e) + character)
		}else {
			buff.WriteString(ParseString(e))
		}
	}
	return buff.String()
}

//create by cwj on 2017-10-17
// parse string from array
func ParseStringFromArrayWithQuote(a []string, character string) string {
	var buff bytes.Buffer
	for i, e := range a {
		if i < len(a) - 1{
			buff.WriteString("'" + ParseString(e) + "'" + character)
		}else {
			buff.WriteString("'" + ParseString(e) + "'")
		}
	}
	return buff.String()
}

func LinkString(character string, empty string, str ...string) string{
	return ParseStringFromArray(str, character, empty)
}

func ParseStringFromMap(m map[string]interface{}) string{
	s, _ := json.Marshal(m)
	return string(s)
}

func ParseMapFromString(s string) (map[string]interface{}, error) {
	var m = make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &m)
	return m, err
}


func SubString(str string, head int, tail int, empty string) string{
	if len(str) < head {
		return empty
	}
	if len(str) < tail{
		return str[head:]
	}
	if head > tail {
		return str[head:] + str[:tail]
	}
	return str[head: tail]
}

//create by cwj on 2017-10-17
//check errArray
//if there are error, return it
func CheckError(errArray ...error) (error){
	for _, e := range errArray{
		if e != nil {
			log.Println(e)
			return e
		}
	}
	return nil
}

//create by cwj on 2017-10-17
//check errArray
//if there are error, return it
func ChangeFromCamelCase(s string) string{
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			s = strings.Replace(s, string(char), "_"+strings.ToLower(string(char)), -1)
		}
	}
	return s
}

//create by cwj on 2017-10-17
//check string
//if there are Single quotation marks('),transfer it for prevent injection
func QuotationTransferred(s interface{}) string{
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	return str
}

func QuotationTransferredForLike(s interface{}) string{
	str := ParseString(s)
	str = strings.Replace(str, "'", "''", -1)
	str = strings.Replace(str, "%", "[%]", -1)
	str = strings.Replace(str, "_", "[_]", -1)
	return str
}

