package views

import (
	"html/template"
	"log"

	"github.com/mashun4ek/webdevcalhoun/gallery/models"
)

const (
	// AlertLvError is used for error alerts
	AlertLvError   = "danger"
	AlertLvWarning = "warning"
	AlertLvInfo    = "info"
	AlertLvSuccess = "success"

	// AlertMsgGeneric - generic message alert for users
	AlertMsgGeneric = "Sorry! Something went wrong. Please try again, and contact us if the problem persists."
)

// Alert is used to render Bootstrap alerts in templates
type Alert struct {
	Level   string
	Message string
}

// Data is used in views
type Data struct {
	// use pionter so you can use nil
	Alert *Alert
	User  *models.User
	CSRF  template.HTML
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok {
		d.Alert = &Alert{
			Level:   AlertLvError,
			Message: pErr.Public(),
		}
	} else {
		log.Println(err)
		d.Alert = &Alert{
			Level:   AlertLvError,
			Message: AlertMsgGeneric,
		}
	}
}

func (d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvError,
		Message: msg,
	}
}

type PublicError interface {
	error
	Public() string
}
