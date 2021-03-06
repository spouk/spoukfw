package spoukfw

import (
	"strconv"
	"reflect"
	"time"
	"strings"
)

const (
	defConverter = "[spoukconverter] `%s`\n"
	prefixLogConverter  = "[spoukconverter][logger]"
	ErrorValueNotValidConvert = "Value not valid for converting"

)
var (
	acceptTypes []interface{} = []interface{}{
		"", 0, int64(0),
	}
)
type (
	SpoukConverter struct {
		logger *SpoukLogger
		value   interface{}
		result  interface{}
		stockFu map[string]func()
	}

)
func NewSpoukConverter() *SpoukConverter {
	f := &SpoukConverter{
		stockFu:make(map[string]func()),
	}
	f.logger = NewSpoukLogger(prefixLogConverter, nil)
	f.stockFu["string"] = f.stringToInt
	f.stockFu["string"] = f.stringToInt64
	return f
}
func (c *SpoukConverter) StrToInt() (*SpoukConverter) {
	if f, exists := c.stockFu["string"]; exists {
		f()
	}
	return c
}
func (c *SpoukConverter) StrToInt64() (*SpoukConverter) {
	if f, exists := c.stockFu["string"]; exists {
		f()
	}
	return c
}
//---------------------------------------------------------------------------
//  String to Int64
//---------------------------------------------------------------------------
func (c *SpoukConverter) stringToInt64() {
	c.stringToInt()
	if c.result != nil {
		c.result = int64(c.result.(int))
	} else {
		c.result = nil
	}
}
//---------------------------------------------------------------------------
//  String to Int
//---------------------------------------------------------------------------
func (c *SpoukConverter) stringToInt() {
	if r, err := strconv.Atoi(c.value.(string)); err != nil {
		c.logger.Printf(makeErrorMessage(defConverter, err.Error()).Error())
		c.result = nil
	} else {
		c.result = r
	}
}
//---------------------------------------------------------------------------
//  возвращает результат конвертации
//---------------------------------------------------------------------------
func (c *SpoukConverter) Result() interface{} {
	return c.result
}
//---------------------------------------------------------------------------
//  инциализация вводным значением
//---------------------------------------------------------------------------
func (c *SpoukConverter) Value(value interface{}) (*SpoukConverter) {
	if c.checkValue(value) {
		c.value = value
		return c
	}
	return nil
}
//---------------------------------------------------------------------------
//  проверка типа поступившего значения на возможность конвертации
//---------------------------------------------------------------------------
func (c *SpoukConverter) checkValue(value interface{}) bool {
	tValue := reflect.TypeOf(value)
	for _, x := range acceptTypes {
		if tValue == reflect.TypeOf(x) {
			return true
		}
	}
	c.logger.Printf(makeErrorMessage(defConverter, ErrorValueNotValidConvert).Error())
	return false
}

func (c *SpoukConverter) DirectStringtoInt64(v string) int64 {
	if res, err := strconv.Atoi(v); err != nil {
		c.logger.Printf(makeErrorMessage(defConverter, err.Error()).Error())
		return 0
	} else {
		return int64(res)
	}
}
//---------------------------------------------------------------------------
//  time convert
//---------------------------------------------------------------------------
func (c *SpoukConverter) ConvertHTMLDatetoUnix(date string) (int64, error) {
	result, err := time.Parse("2006-01-02", date)
	if err == nil {
		return result.Unix(), err
	} else {
		return 0, err
	}
}
func (c *SpoukConverter) ConvertUnixTimeToString(unixtime int64) string {
	return time.Unix(unixtime, 0).String()
}
//convert timeUnix->HTML5Datatime_local(string)
func (c *SpoukConverter) TimeUnixToDataLocal(unixtime int64) string {
	tmp_result := time.Unix(unixtime,0).Format(time.RFC3339)
	g := strings.Join(strings.SplitAfterN(tmp_result,":",3)[:2],"")
	return g[:len(g)-1]
}
//convert HTML5Datatime_local(string)->TimeUnix
func (c *SpoukConverter) DataLocalToTimeUnix(datatimeLocal string) int64 {
	r,_ := time.Parse(time.RFC3339, datatimeLocal+":00Z")
	return r.Unix()
}