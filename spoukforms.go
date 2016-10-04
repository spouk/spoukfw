package spoukfw


import (
	"strconv"
	"reflect"
	"errors"
	"strings"
	"fmt"
	"log"
	"database/sql"
)

var (
//сообщения об ошибках в формах
	ErrorUsername = "Имя пользователя ошибочно"
	ErrorPassword = "Пароль ошибочен"
	ErrorEmail = "Почтовый адрес ошибочен"

//`placeholder` описания для формы
	PlaceUsername = "= имя пользователя = "
	PlacePassword = "= пароль ="
	PlaceEmail = "= почтовый адрес ="

//ошибки
	ParseErrorInt = errors.New("[parseform][error] ошибка парсинга `string`->`int64`")
	PTRFormError = errors.New("[baseform][error] Ошибка, дай мне указатель на структуру для записи")
	PTRFormErrorMethods = errors.New("[baseform][error] Ошибка, отсутствует реализация интерфейса методов для получения данных из формы")
	CSRFErrorValidate = "CSRF не валидное значение"

//название стилей для ошибок в формах полей
	ErrorStyleForm = "has-error"
	SuccessStyleForm = "has-success"

//сообщения для ошибки в формах при валидации формы
	ErrorMsgFormString = "- поле не может быть пустым -"
	ErrorMsgFormCheckbox = "- нажмите на чекбокс, если вы не робот -"
	ErrorMsgFormBool = "- сделайте отметку -"
	ErrorMsgFormSelect = "- не выбран ни один из элементов -"
)

type (
//структура для дефолтных значений
	DefaultForm struct {
		ErrorMsg          string
		ErrorClassStyle   string
		SuccessClassStyle string
		SuccessMessage    string
		Placeholder       string
	}

//	интерфейс для методов получения данных из формы,
//  	для полного абстрагирования от всякого говна, типа фреймворков
	MethodsForms interface {
		GetMultiple(name string) []string
		GetSingle(name string) string
	}

//	интерфейс для юзерских проверок
	UserForm interface {
		Validate(b *SpoukForm, stock ...interface{}) bool
	}
//	интерфейс для csrf валидации
	CSRFValidate interface {
		Validate() bool
	}
//базовая структура для всех форм
	SpoukForm struct {
		ParseWithInit bool
		Errors        map[string]string
		Class         map[string]string
		Desc          map[string]string
		Stack         map[string]interface{}
		CSRFValidate  CSRFValidate
		MethodsForms  MethodsForms
		DefaultForm   map[string]DefaultForm
	}
	StockValue struct {
		Name  string
		Value interface{}
	}
	Stocker struct {
		Stock map[string]interface{}
	}
)
//---------------------------------------------------------------------------
//  Стокер структура для внутреннего использования внутри либы
//---------------------------------------------------------------------------
func NewStocker() *Stocker {
	return &Stocker{Stock:make(map[string]interface{})}
}
//---------------------------------------------------------------------------
//  создания инстанса + основной функционал
//---------------------------------------------------------------------------
func NewSpoukForm(defaultForms map[string]DefaultForm, methodsForm MethodsForms, ParseWithInit bool) *SpoukForm {
	form := new(SpoukForm)
	form.ParseWithInit = ParseWithInit
	form.MethodsForms = methodsForm
	form.DefaultForm = defaultForms
	form.Errors = make(map[string]string)
	form.Class = make(map[string]string)
	form.Desc = make(map[string]string)
	form.Stack = make(map[string]interface{})
	return form
}

var (
	DefaultValues map[string]DefaultForm = map[string]DefaultForm{
		"Name" : DefaultForm{Placeholder:"=введите имя пользователя=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Username" : DefaultForm{Placeholder:"=введите имя пользователя=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Password" : DefaultForm{Placeholder:"=введите пароль =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Email" : DefaultForm{Placeholder:"=введите email =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Port" : DefaultForm{Placeholder:"=порт сервера=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"CatName" : DefaultForm{Placeholder:"=введите название категории=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Title" : DefaultForm{Placeholder:"=введите заголовок =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"MetaKeys" : DefaultForm{Placeholder:"=введите SEO слова =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"MetaDesc" : DefaultForm{Placeholder:"=введите SEO описание-сниппет =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"MetaRobot" : DefaultForm{Placeholder:"=введите занчения для SEO robot=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Message" : DefaultForm{Placeholder:"=введите текст сообщения=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Body" : DefaultForm{Placeholder:"=введите тело сообщения =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Link" : DefaultForm{Placeholder:"=введите ссылку-ключ =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"Age" : DefaultForm{Placeholder:"=введите ваш возраст=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},

		"UserInfoUsername": DefaultForm{Placeholder:"=введите имя пользователя=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"UserInfoPassword" : DefaultForm{Placeholder:"=введите пароль =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"UserEmail" : DefaultForm{Placeholder:"=введите email =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"CategoryName" : DefaultForm{Placeholder:"=введите название категории=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostTitle" : DefaultForm{Placeholder:"=введите заголовок =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostBody" : DefaultForm{Placeholder:"=введите тело сообщения =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostSeoMetaKeys" : DefaultForm{Placeholder:"=введите SEO слова =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostSeoMetaDesc" : DefaultForm{Placeholder:"=введите SEO описание-сниппет =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostSeoMetaRobot" : DefaultForm{Placeholder:"=введите занчения для SEO robot=", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"TagName" : DefaultForm{Placeholder:"=введите имя метки =", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поле не может быть пустым"},
		"PostCategoryID": DefaultForm{Placeholder:"", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"значение должно быть выбрано"},
		"PostUserID": DefaultForm{Placeholder:"", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"значение должно быть выбрано"},
		"Robot" : DefaultForm{Placeholder:"", ErrorClassStyle:"has-error", SuccessClassStyle:"ok", ErrorMsg:"поставьте отметку что вы не робот"},


	}
)
//---------------------------------------------------------------------------
//  functions
//---------------------------------------------------------------------------
func (b *SpoukForm)ResetForm() {
	b.Errors = make(map[string]string)
	b.Class = make(map[string]string)
	b.Desc = make(map[string]string)
	b.Stack = make(map[string]interface{})
}
func (b *SpoukForm)AddForm(name string, form DefaultForm) {
	b.DefaultForm[name] = form
}
func (b *SpoukForm)resetErrors() {
	b.Errors = make(map[string]string)
}
func (b *SpoukForm)InitForm(form UserForm) {

	//проводит первичную инциализацию формы = присваивает placeholder`s
	main := reflect.ValueOf(form)
	numfield := reflect.ValueOf(form).Elem().NumField()
	if main.Kind() != reflect.Ptr {
		log.Fatal(PTRFormError)
	}
	//провожу заполнение
	for x := 0; x < numfield; x++ {
		//получаем имя элемента структуры
		name := reflect.TypeOf(form).Elem().Field(x).Name
		//получаю placeholder из дефолтного стека
		if def, exists := b.DefaultForm[name]; exists {
			b.Desc[name] = def.Placeholder
		}

	}
	//провожу парсинг формы //post form
	if b.ParseWithInit {
		b.ParseForm(form)
	}
}
//func (b *SpoukForm)LoadForm(obj interface{}, form UserForm) {
//	//загружаю данные из принятой формы с данными в объект по таким же именам
//	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(form).Kind() != reflect.Ptr {
//		log.Fatal(PTRFormError)
//	}
//	//получаю стокер со списком филдов объекта из базы данных как правило
//	stocker := b.ParseFields(obj)
//	//заполняет форму данными из объекта, используется при обновлении объектов, как пример
//	mv := reflect.ValueOf(form).Elem()
//	mt := reflect.TypeOf(form).Elem()
//	//провожу заполнение
//	for x := 0; x < mv.NumField(); x++ {
//		//получаем имя элемента структуры
//		name := mt.Field(x).Name
//		//пробую получить элемент с таким же названием из объекта +
//		//if v, ok := stocker.Stock[name]; ok {
//		//
//		//	b.Stack[name] = v
//		//}
//	}
//}
func (b *SpoukForm)UpdateForm(form UserForm, obj interface{}) {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr || reflect.ValueOf(form).Kind() != reflect.Ptr {
		log.Fatal(PTRFormError)
	}
	//получаю стокер со списком филдов объекта из базы данных как правило
	stocker := b.ParseFields(obj)
	//заполняет форму данными из объекта, используется при обновлении объектов, как пример
	mv := reflect.ValueOf(form).Elem()
	mt := reflect.TypeOf(form).Elem()
	//провожу заполнение
	for x := 0; x < mv.NumField(); x++ {
		//получаем имя элемента структуры
		name := mt.Field(x).Name
		//пробую получить элемент с таким же названием из объекта +
		if v, ok := stocker.Stock[name]; ok {
			b.Stack[name] = v
		}
	}
}
func (b *SpoukForm) ParseFields(obj interface{}) *Stocker {
	//рекурсивно собираю все поля полученном объекте
	stocker := NewStocker()
	for _, x := range b.UpdateFormDeep(obj) {
		stocker.Stock[x.Name] = x.Value
		fmt.Printf("[stock] Name: `%v`     Value : `%v`\n", x.Name, x.Value)
	}
	return stocker
}
func (b *SpoukForm)UpdateFormDeep(form interface{}) []StockValue {
	//локальные переменки
	stockFields := make([]StockValue, 0)
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
	//рекурсивно получает список всех филдов объекта для последующего отбора
	for x := 0; x < mv.NumField(); x++ {
		v := mv.Field(x)

		switch v.Kind() {
		//case reflect.Struct:
		//
		//	stockFields = append(stockFields, b.UpdateFormDeep(v.Interface())...)
		default:
			s := StockValue{}
			name := mt.Field(x).Name
			s.Name = defaultStructName + name
			//s.Name = name
			s.Value = v.Interface()
			stockFields = append(stockFields, s)
		}
	}
	return stockFields
}
func (b *SpoukForm) ParseForm(obj interface{}) {
	//проверка на наличие реализации методово интерфейса для получения данных из формы
	if b.MethodsForms == nil {
		log.Fatal(PTRFormErrorMethods)
	}
	main := reflect.ValueOf(obj)
	numfield := reflect.ValueOf(obj).Elem().NumField()
	if main.Kind() != reflect.Ptr {
		log.Fatal(PTRFormError)
	}
	//перебор элементов структуры, получение их имен и получение данных из формы
	//с дальнейшим присваиванием записям структуры

	for x := 0; x < numfield; x++ {
		//получаем элемент структуры
		f := reflect.Indirect(reflect.ValueOf(obj)).Field(x)
		//получаем имя элемента структуры
		name := reflect.TypeOf(obj).Elem().Field(x).Name

		switch f.Type().Kind() {

		case reflect.Struct:
			value := f.Interface()
			val := strings.TrimSpace(b.MethodsForms.GetSingle(name))
			v := sql.NullString{String:val, Valid:true}
			switch value.(type) {
			case sql.NullString, sql.NullFloat64, sql.NullInt64:
				b.Stack[name] = v
			}

		case reflect.Slice, reflect.Array:
			//проводим общие для всех операции
			//получаю данные из формы
			//c.Request().ParseMultipartForm(0)
			//formList2 := c.Request().Form[name]

			formList := b.MethodsForms.GetMultiple(name)
			value := f.Interface()

			switch value.(type) {
			case []int64:
				tmp := []int64{}
				for _, v := range formList {
					if parInt, errPat := strconv.ParseInt(v, 10, 64); errPat == nil {
						tmp = append(tmp, parInt)
					}
				}
				//добавление данных в baseform.stack
				b.Stack[name] = tmp
				//структура готова, можно менять
				f.Set(reflect.ValueOf(&tmp).Elem())
			case []string:
				tmp := []string{}
				for _, v := range formList {
					tmp = append(tmp, v)
				}
				//добавление данных в baseform.stack
				b.Stack[name] = tmp
				//меняем
				f.Set(reflect.ValueOf(&tmp).Elem())
			}
		case reflect.Bool:
			val := strings.TrimSpace(b.MethodsForms.GetSingle(name))
			if val != "" {
				f.SetBool(true)
				b.Stack[name] = true
			}  else {
				f.SetBool(false)
				b.Stack[name] = false
			}

		case reflect.String:
			val := strings.TrimSpace(b.MethodsForms.GetSingle(name))
			f.SetString(val)
			b.Stack[name] = val

		case reflect.Int, reflect.Int64:
			value := strings.TrimSpace(b.MethodsForms.GetSingle(name))
			if value != "" {
				r, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					log.Printf("%s", ParseErrorInt)
					log.Printf("%s", err)
				} else {
					f.SetInt(r)
					b.Stack[name] = r
				}
			}
		case reflect.Float64, reflect.Float32:
			value := strings.TrimSpace(b.MethodsForms.GetSingle(name))
			if value != "" {
				r, err := strconv.ParseFloat(value, 64)
				if err != nil {
					log.Printf("%s", ParseErrorInt)
					log.Printf("%s", err)
				} else {
					f.SetFloat(r)
					b.Stack[name] = r
				}
			}
		}
	}
}
func (b *SpoukForm) ValidateForm(form UserForm) (status bool) {
	//проверка формы на валидность полей
	main := reflect.ValueOf(form)
	numfield := reflect.ValueOf(form).Elem().NumField()
	if main.Kind() != reflect.Ptr {
		log.Fatal(PTRFormError)
	}
	//обнуляю стек ошибок
	b.resetErrors()
	//проверка CSRF если есть объект реализовавший интерфейс
	if b.CSRFValidate != nil {
		//проверка валидности CSRF значения
		if b.CSRFValidate.Validate() == false {
			b.Errors["CSRF"] = CSRFErrorValidate
			return false
		}
	}
	//проверка на валидность пользовательской проверкой
	if form.Validate(b) == false {
		return false
	}

	//количество флагов = количество полей в форме
	var total int = numfield
	var countValidate int = 0
	//проверка дефолтных значений и полей
	for x := 0; x < numfield; x++ {
		f := reflect.Indirect(reflect.ValueOf(form)).Field(x)
		name := reflect.TypeOf(form).Elem().Field(x).Name
		ff := reflect.TypeOf(form).Elem().Field(x)

		var def *DefaultForm
		if do, exists := b.DefaultForm[name]; exists {
			def = &do
		}
		switch f.Type().Kind() {
		default:
			fmt.Printf("[reflect][validateform][ALERT] непроверяемый тип Name: `%v` Value: %v\n", name, f.Interface())

		case reflect.Float64, reflect.Float32:
			if ff.Tag != "" {
				total --
			} else {
				result := f.Interface().(float64)
				if result == 0 {
					//error
					if def != nil {
						b.Class[name] = def.ErrorClassStyle
						b.Errors[name] = def.ErrorMsg
					}

				} else {
					if def != nil {
						b.Class[name] = def.SuccessClassStyle
					}
					countValidate ++
				}
			}
		case reflect.Int64,reflect.Int32, reflect.Int16, reflect.Int:
			if ff.Tag != "" {
				total --
			} else {
				result := f.Interface().(int64)
				if result == 0 {
					//error
					if def != nil {
						b.Class[name] = def.ErrorClassStyle
						b.Errors[name] = def.ErrorMsg
					}

				} else {
					if def != nil {
						b.Class[name] = def.SuccessClassStyle
					}
					countValidate ++
				}
			}
		case reflect.Slice, reflect.Array:
			//разбор по типу списка
			value := f.Interface()
			switch value.(type) {
			case []int64:
				//проверка на метку необходимости проверки
				//если метка присутствует, проверка не нужна
				if ff.Tag != "" {
					total --
				} else {
					result := value.([]int64)
					if len(result) == 0 {
						//error
						if def != nil {
							b.Class[name] = def.ErrorClassStyle
							b.Errors[name] = def.ErrorMsg
						}
					} else {
						if def != nil {
							b.Class[name] = def.SuccessClassStyle
						}
						countValidate ++
					}
				}
			case []string:

				if ff.Tag != "" {
					total --
				} else {
					result := value.([]string)
					if len(result) == 0 {
						//error
						if def != nil {
							b.Class[name] = def.ErrorClassStyle
							b.Errors[name] = def.ErrorMsg
						}
					} else {
						if def != nil {
							b.Class[name] = def.SuccessClassStyle
						}
						countValidate ++
					}

				}
			}
		case reflect.String:
			tag := ff.Tag
			if tag != "" {
				total --
			} else {
				result := strings.TrimSpace(f.Interface().(string))
				if result == "" {
					//error
					if def != nil {
						b.Class[name] = def.ErrorClassStyle
						b.Errors[name] = def.ErrorMsg
					}
					status = false
				} else {
					if def != nil {
						b.Class[name] = def.SuccessClassStyle
					}

					countValidate ++
				}
			}

		case reflect.Bool:
			tag := ff.Tag
			if tag != "" {
				total --
			} else {
				result := f.Interface().(bool)
				if result == false {
					//error
					if def != nil {
						b.Class[name] = def.ErrorClassStyle
						b.Errors[name] = def.ErrorMsg
					}

				} else {
					if def != nil {
						b.Class[name] = def.SuccessClassStyle
					}
					countValidate ++
				}
			}
		}

	}

	//подведение итогов по валидности всей формы
	if total == countValidate {
		//fmt.Printf("[validateform] Total: %v   Numfield: %v   CountValidate: %v , Result: VALIDATE\n", total, numfield, countValidate)
		status = true
	} else {
		status = false
		//fmt.Printf("[validateform] Total: %v   Numfield: %v   CountValidate: %v , Result: NOT VALIDATE\n", total, numfield, countValidate)
	}
	return
}
func (b *SpoukForm)ConvertSliceINT64(name string) []int64 {
	v := b.Stack[name]
	result := []int64{}
	if v != nil {
		result = v.([]int64)
	}
	return result
}
func (b *SpoukForm)ConvertString(name string) string {
	v := b.Stack[name]
	var result string
	if v != nil {
		result = v.(string)
	}
	return result
}
func (b *SpoukForm)ConvertInt(name string) int64 {
	v := b.Stack[name]

	var result int64
	if v != nil {
		result = int64(v.(int))
	}
	return result
}
func (b *SpoukForm)ConvertBool(name string) bool {
	v := b.Stack[name]
	var result bool
	if v != nil {
		result = v.(bool)
	}
	return result
}

