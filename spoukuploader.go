package spoukfw

import (
	"sync"
	"fmt"
	"os"
	"io"
	"path/filepath"
	"mime/multipart"
)
const (
	SPOUKCARRYUPLOAD = "[spoukuploader][uploadfile][ERROR] %v\n"
	SPOUKCARRYUPLOADOK = "[spoukuploader][uploadfile][OK] `%v`\n"
)

type (
	FileInfo struct {
		Name string
		Path string
		Ext  string
		Size uint
	}
	SpoukUploader struct {
		Stock []FileInfo
	}
)
var (
	w sync.WaitGroup
)
func NewSpoukUploader() *SpoukUploader {
	return &SpoukUploader{}
}
//---------------------------------------------------------------------------
//  загрузка одиночного файла с участием AJAX
//---------------------------------------------------------------------------
func (u *SpoukUploader) UploadSingleAJAX(formName string, sr *SpoukCarry) *FileInfo {
	var f FileInfo

	err := sr.request.ParseForm()
	if err != nil {
		sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err.Error()))
	}
	_, fh, err_open := sr.request.FormFile(formName)
	if err_open != nil {
		sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_open.Error()))
		return nil
	}

	fin, err_inopen := fh.Open()
	if err_inopen != nil {
		sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_inopen.Error()))
		return nil
	} else {
		defer fin.Close()
		//создаю файл на локальной машине - приемный файл
		fout, err_fout := os.Create(sr.Config().UPLOADFilesPath  + fh.Filename)
		if err_fout != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fout.Error()))
			return nil
		} else {
			defer fout.Close()
			//копирую файл
			_, err_read := io.Copy(fout, fin)
			if err_read != nil {
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_read.Error()))
				return nil
			} else {
				//получаю данные по файлу
				info, err_fi := fout.Stat()
				if err_fi != nil {
					sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fi.Error()))
					return nil
				} else {
					f = FileInfo{Name: fh.Filename, Path:sr.Config().UPLOADFilesPath  + fh.Filename, Ext: filepath.Ext(info.Name()), Size: uint(info.Size())}
				}
				//success upload file
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOADOK, fh.Filename))
				return &f
			}
		}

	}
}
//---------------------------------------------------------------------------
// загрузка одиночного файла, функцию оптимально использовать как горутину при обработке `multiple`
//---------------------------------------------------------------------------
func (u *SpoukUploader) goUploadSingle(fh *multipart.FileHeader, ajax bool, sr *SpoukCarry) *FileInfo {
	var f FileInfo

	//не аякс, значит горутина
	if !ajax {
		defer func() {
			u.Stock = append(u.Stock, f)
			w.Done()
		}()
	}

	fin, err_inopen := fh.Open()
	if err_inopen != nil {
		sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_inopen.Error()))
		return nil
	} else {
		defer fin.Close()
		//создаю файл на локальной машине - приемный файл
		fout, err_fout := os.Create(sr.Config().UPLOADFilesPath  + fh.Filename)
		if err_fout != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fout.Error()))
			return nil
		} else {
			defer fout.Close()
			//копирую файл
			_, err_read := io.Copy(fout, fin)
			if err_read != nil {
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_read.Error()))
				return nil
			} else {
				//получаю данные по файлу
				info, err_fi := fout.Stat()
				if err_fi != nil {
					sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_fi.Error()))
					return nil
				} else {
					f = FileInfo{Name: fh.Filename, Path:sr.Config().UPLOADFilesPath  + fh.Filename, Ext: filepath.Ext(info.Name()), Size: uint(info.Size())}
				}
				//success upload file
				sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOADOK, fh.Filename))
				return &f
			}
		}

	}
}
//---------------------------------------------------------------------------
//  загрузка single/multiple формы без участия сторонних асинхронных методов типа ajax
//---------------------------------------------------------------------------
func (u *SpoukUploader) Upload(nameForm string, ajax bool, sr *SpoukCarry) (error) {
	//получаю список файлов мультиформ
	if !ajax {
		//в цикле открывая дескрипторы файлов для загрузки и файлы для принятия
		//запуск на каждый дескриптор горутину с синхронизацией загрузки
		//ajax == false
		err := sr.request.ParseMultipartForm(32 << 20)
		if err != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err.Error()))
			return err
		}
		formdata := sr.Request().MultipartForm
		listfiles := formdata.File[nameForm]
		fmt.Printf("LISTFILES===> %v\n", listfiles)

		w.Add(len(listfiles))
		for _, INfile := range listfiles {
			go u.goUploadSingle(INfile, false, sr)
		}
		//ожидаю завершение закачки
		w.Wait()
	} else {
		//ajax, запуск без синхрона, т.к. каждый выхов аякса дергает жту функцию, которая сама
		//запускается как горутина
		//sr.request.ParseMultipartForm(32 << 20)
		err := sr.request.ParseForm()
		if err != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err.Error()))
			return err
		}
		_, handler, err_form := sr.request.FormFile("uploadfile")
		if err_form != nil {
			sr.Spoukmux.logger.Error(fmt.Sprintf(SPOUKCARRYUPLOAD, err_form.Error()))
			return err_form
		}
		f := u.goUploadSingle(handler, true, sr)
		if f != nil {
			u.Stock = append(u.Stock, *f)
		}
	}
	fmt.Println("[sync] All uploading")
	fmt.Printf("StockFIles: %v\n", u.Stock)
	return nil
}
