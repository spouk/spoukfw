package spoukfw

import (
	"crypto/md5"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"time"
	"errors"
)

type (
	SpoukCSRF struct {
		TimeActive time.Duration
		TimeStart  time.Time
		Salt       string
		Key        *uuid.UUID
		ReadyKey   string
		Csrf_form  func() (*string, error)
		Csrf_head  func() (*string, error)
	}
)
func NewSpoukCSRF(minutesActive int, salt string) (*SpoukCSRF) {
	n := &SpoukCSRF{
		TimeActive: time.Duration(minutesActive) * time.Minute,
		TimeStart: time.Now(),
		Salt:salt,
	}
	u, err := uuid.NewV4()
	if err != nil {
		return nil
	}
	n.Key = u
	n.Csrf_form = n.wrapper(true, false)
	n.Csrf_head = n.wrapper(false, true)
	return n
}
func (c *SpoukCSRF) wrapper(form, head bool) (func() (*string, error)) {
	return func() (*string, error) {
		_tmptime := c.TimeStart.Add(c.TimeActive)
		if _tmptime.Before(time.Now()) {
			c.TimeStart = time.Now()
			u, err := uuid.NewV4()
			if err != nil {
				return nil, errors.New("[CSRF] [csrf_form] error generate uuid")
			}
			c.Key = u
		}
		r := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%v%v", c.Key, c.Salt))))
		c.ReadyKey = r
		var result string
		if form {
			result = fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s"> `, r)
		} else if head {
			result = fmt.Sprintf(`<meta id="csrf_token_ajax" content="%s" name="csrf_token_ajax" />`, r)
		}
		return &result, nil
	}
}
func (c SpoukCSRF) VerifyToken(s *SpoukCarry) bool {
	token := s.GetFormValue("csrf_token")
	if (token == c.ReadyKey) {
		return true
	}
	return false
}
func (c SpoukCSRF) VerifyTokenString(token string) bool {
	if (token == c.ReadyKey) {
		return true
	}
	return false
}


