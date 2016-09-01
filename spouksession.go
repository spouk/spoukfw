package spoukfw

//import "fmt"

func newSpoukSession(spoukmux *Spoukmux) *SpoukSession {
	return &SpoukSession{
		Spoukmux:spoukmux,
		Paginate:NewSpoukPaginate(spoukmux),
		Csrf :NewSpoukCSRF(spoukmux.config.CSRFActiveMinute, spoukmux.config.CSRFSalt),
		Flasher:newSpoukFlasher(),
		Mail:newSpoukMail(),
		Convert:NewSpoukConverter(),
	}
}

type (
//фабрика сессий, с сохранением состояния между реквестами, статичный инстанс
	SpoukSession struct {
		Spoukmux      *Spoukmux
		Paginate      *SpoukPaginate
		Csrf          *SpoukCSRF
		Flasher       *SpoukFlasher
		Mail          *SpoukMail
		Convert       *SpoukConverter
		SessionObject SpoukSessionObject
		//Conf          *StackConfig
		//SpoukUploader *SpoukUploader
	}


//---------------------------------------------------------------------------
//  сессионный динамический объект, интерфейс для него
//---------------------------------------------------------------------------
//	SpoukSessionObject interface {
//		NewSpoukSessionObject(spoukcarry *SpoukCarry, mail *SpoukMail, flasher *SpoukFlasher, csrf *SpoukCSRF,
//		spoukPaginate *SpoukPaginate, spoukmux *Spoukmux, conver *SpoukConverter, data *DataSO) *SessionObject
//		InitSpoukSessionObject()
//	}
	SpoukSessionObject interface {
		NewSpoukSessionObject(s *SpoukSession, c *SpoukCarry) interface{}
		InitSpoukSessionObject(s interface{})
	}
	//динамическая структура данных, передается по контексту, должна быть реализована разработчиком приложение испоьзуемый фреймворк + сессии

	DataSO map[string]interface{}
)
var (
	ErrorSession = "[spouksession][error] `%s`\n"
)

//---------------------------------------------------------------------------
//  DataSO
//---------------------------------------------------------------------------
func NewDataSo() *DataSO {
	d := make(DataSO)
	d.setDefaultSection()
	return &d
}
func (s DataSO) setDefaultSection() {
	for _, x := range []string{"admin", "public", "user", "flash", "stack", "seo"} {
		s[x] = make(map[string]interface{})
	}
}
func (s DataSO) checkSection(section string) {
	for _, v := range s {
		if v == section {
			return
		}
	}
	s[section] = make(map[string]interface{})
	return
}
func (s DataSO) SetUserData(key string, value interface{}) {
	s["user"].(map[string]interface{})[key] = value
}
func (s DataSO) GetUserData(key string) (interface{}) {
	return s["user"].(map[string]interface{})[key]
}
func (s DataSO) Set(section, key string, value interface{}) {
	s[section].(map[string]interface{})[key] = value
}
func (s DataSO) Get(section, key string) (interface{}) {
	return s[section].(map[string]interface{})[key]
}
