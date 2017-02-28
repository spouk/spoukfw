package spoukfw

import (
	"time"
	"errors"
	"fmt"
	"strings"
	"reflect"
	"os"
	"io/ioutil"
	"encoding/json"
)

const (
	def  = "[spoukconfig] `%v`"
	prefixLogConfig = "[spoukconfig]"
)
var (
	ErrorDBSNotReleased = errors.New("[configproxy][error] DBSInterface not released this instanse")
)

var (
	cFG_HTTPREADTIMEOUT = 5 * time.Second
	cFG_HTTPWRITETIMEOUT = 7 * time.Second
	cFG_ADDRESS = ":8090"
	cFG_TEMPLATEPATH = "templates/*.html"
	cFG_TEMPLATEDEBUG = true
	cFG_USESESSION = true
	cFG_COOKNAME = "spoukftw"
	cFG_HOSTNAME = ""
	cFG_SESSIONTIME = int64(30) //minutes
	cFG_COUNTONPAGE = int64(5) //minutes
	cFG_COUNTLINKS = int64(5) //minutes
	cFG_CSRFACTIVEMINUTE = 60 //minutes
	cFG_CSRFSALT = "somesaltforCSRF"//minutes
	cFG_UPLOADFilesPath = "files/"
	cFG_HTTP  = ":8090"
	cFG_HTTPS = ":8091"
)

type (
	Spoukconfig struct {
		//http timeouts
		HTTPReadTimeout  time.Duration
		HTTPWriteTimeout time.Duration
		//addr
		Address          string
		//ssl + no ssl
		HTTP 	         string
		HTTPS 		 string
		//template
		TemplatePath     string
		TemplateDebug    bool
		//session
		UseSession       bool
		SessionTime      int64
		//cookies
		CookName         string
		//host info
		Hostname         string
		//paginate
		CountOnPage      int64
		CountLinks       int64
		//csrf
		CSRFActiveMinute int
		CSRFSalt         string
		//Files download path default
		UPLOADFilesPath  string
		//---------------------------------------------------------------------------
		//  проксификатор для внешнего файла конфига, реализованного в .json формате
		//---------------------------------------------------------------------------
		ConfigProxy      *SpoukConfigProxy
		logger           *SpoukLogger
	}
//---------------------------------------------------------------------------
//  spoukProxyConfig 
//---------------------------------------------------------------------------
	DBSInterface  interface {
		WriteConfigDB() error
		ReadConfigFromDB() error
	}
	SpoukConfigProxy struct {
		ConfigFromFile interface{}
		CurrentConfig  interface{}
		ConfigFile     string
		StructFile     interface{}
		Mapper         *ProxyMapperConfig
		dbsinterface   DBSInterface
	}
	ProxyMapperValue struct {
		Name  string
		Value interface{}
	}
	ProxyMapperConfig  struct {
		Stock map[string]interface{}
	}
)

func NewSpoukconfig() *Spoukconfig {
	sp := &Spoukconfig{
		HTTPReadTimeout:cFG_HTTPREADTIMEOUT,
		HTTPWriteTimeout:cFG_HTTPWRITETIMEOUT,
		Address:cFG_ADDRESS,
		TemplatePath:cFG_TEMPLATEPATH,
		TemplateDebug:cFG_TEMPLATEDEBUG,
		UseSession:cFG_USESESSION,
		CookName:cFG_COOKNAME,
		Hostname:cFG_HOSTNAME,
		SessionTime:cFG_SESSIONTIME,
		CountLinks:cFG_COUNTLINKS,
		CountOnPage:cFG_COUNTONPAGE,
		CSRFActiveMinute:cFG_CSRFACTIVEMINUTE,
		CSRFSalt:cFG_CSRFSALT,
		UPLOADFilesPath:cFG_UPLOADFilesPath,
		HTTP:cFG_HTTP,
		HTTPS:cFG_HTTPS,
	}
	sp.logger = NewSpoukLogger(prefixLogConfig, nil)
	return sp
}
//---------------------------------------------------------------------------
//  SpoukConfig
//---------------------------------------------------------------------------
func (s *Spoukconfig) InjectExternConfigJson(fileconfigjson string, structFile interface{}, dbs DBSInterface) (error) {
	if cp, err := NewSpoukProxyConfig(fileconfigjson, structFile, dbs); err == nil {
		s.ConfigProxy = cp
		return nil
	} else {
		s.logger.Printf(err.Error())
		return err
	}
}
//---------------------------------------------------------------------------
//  SpoukProxyConfig
//---------------------------------------------------------------------------
func NewSpoukProxyConfig(fileconfigjson string, structFile interface{}, dbs DBSInterface) (*SpoukConfigProxy, error) {
	c := &SpoukConfigProxy{ConfigFile:fileconfigjson}
	if dbs != nil {
		c.dbsinterface = dbs
	}
	c.StructFile = reflect.New(reflect.TypeOf(structFile)).Interface()
	if err, status := c.readfromfile(); !status {
		return nil, err
	}
	c.Mapper = c.convertConfigToMapper()
	return c, nil
}
func (c *SpoukConfigProxy) readfromfile() (err error, status bool) {
	var path string
	if path, err = os.Getwd(); err != nil {
		return err, false
	}
	var data []byte
	if data, err = ioutil.ReadFile(path + "/" + c.ConfigFile); err != nil {
		return err, false
	}

	if err = json.Unmarshal(data, c.StructFile); err != nil {
		return err, false
	}
	return nil, true
}
func (c *SpoukConfigProxy) convertConfigToMapper() *ProxyMapperConfig {
	mapper := new(ProxyMapperConfig)
	mapper.Stock = make(map[string]interface{})

	for _, x := range c.converterStructtoMapper(c.StructFile) {
		mapper.Stock[x.Name] = x.Value
		//fmt.Printf("::: [proxy][mapper] `%30s`  :::  `%30v`\n", x.Name, x.Value)
	}
	return mapper
}
func (c *SpoukConfigProxy) converterStructtoMapper(form interface{}) []ProxyMapperValue {
	//парсит структуру для конвертации структуры в карту
	//локальные переменки
	stockFields := make([]ProxyMapperValue, 0)
	var mv reflect.Value
	var mt reflect.Type

	//form может поступать в 2 видах как указатель так и по значению, отсюда надо ветвить
	switch reflect.ValueOf(form).Kind() {
	case reflect.Ptr:
		mv = reflect.ValueOf(form).Elem()
		mt = reflect.TypeOf(form).Elem()
	default:
		mv = reflect.ValueOf(form)
		mt = reflect.TypeOf(form)
	}
	//дефолтное имя первичноый структуры
	var defaultStructName string
	nameStruct := strings.Split(fmt.Sprintf("%T", mv.Interface()), ".")
	if len(nameStruct) >= 2 {
		defaultStructName = nameStruct[1]
	}
	//defaultStructName = nameStruct
	//рекурсивно получает список всех филдов объекта для последующего отбора
	for x := 0; x < mv.NumField(); x++ {
		v := mv.Field(x)

		switch v.Kind() {
		case reflect.Struct:
			stockFields = append(stockFields, c.converterStructtoMapper(v.Interface())...)
		default:
			s := ProxyMapperValue{}
			name := mt.Field(x).Name
			s.Name = defaultStructName + "_" + name
			s.Value = v.Interface()
			stockFields = append(stockFields, s)
		}
	}
	return stockFields
}
func (c *SpoukConfigProxy) WriteConfigDB() error {
	if c.dbsinterface == nil {
		return ErrorDBSNotReleased
	}
	return c.dbsinterface.WriteConfigDB()
}
func (c *SpoukConfigProxy) ReadConfigFromDB() error {
	if c.dbsinterface == nil {
		return ErrorDBSNotReleased
	}
	return c.dbsinterface.ReadConfigFromDB()
}
func (c *SpoukConfigProxy) ShowType(s interface{}) {
	fmt.Printf("[proxyconfig] `%T`\n", s)
}


