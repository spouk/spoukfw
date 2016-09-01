package spoukfw

import (
	"fmt"
	"crypto/md5"
	"time"
	"sync"
)

type (
	SpoukFlasher struct {
		sync.RWMutex
		Key   string
		Stock map[string]*FlashMessage
	}
	FlashMessage struct {
		Status  string
		Message interface{}
	}
)
//---------------------------------------------------------------------------
//  functions
//---------------------------------------------------------------------------
func NewSpoukFlasher() *SpoukFlasher {
	n := &SpoukFlasher{
		Stock:make(map[string]*FlashMessage),
	}
	n.Key = n.generateKey()
	return n
}
func newSpoukFlasher() *SpoukFlasher {
	n := &SpoukFlasher{
		Stock:make(map[string]*FlashMessage),
	}
	n.Key = n.generateKey()
	return n
}
func (f *SpoukFlasher) generateKey() string {
	t := time.Now()
	return fmt.Sprintf("%x", md5.Sum([]byte(t.String() + "SoimeVowodfgldkfgjdlfkgj")))
}
func (f *SpoukFlasher)Set(status, section string, message interface{}) {
	nm := &FlashMessage{Status:status, Message:message}
	f.Lock()
	f.Stock[section] = nm
	f.Unlock()
}
func (f *SpoukFlasher)Get(section string) (*FlashMessage) {
	f.Lock()
	defer f.Unlock()
	if  result, exists := f.Stock[section]; exists {
		delete(f.Stock, section)
		return result
	}
	return nil
}
