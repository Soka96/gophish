package models

import (
	"errors"
	"time"

	log "github.com/gophish/gophish/logger"
	"github.com/jinzhu/gorm"
)

// Template models hold the attributes for an email template to be sent to targets
type EchoEmail struct {
	Id           int64        `json:"id" gorm:"column:id; primary_key:yes"`
	UserId       int64        `json:"-" gorm:"column:user_id"`
	Name         string       `json:"name"`
	Subject      string       `json:"subject"`
	HTML         string       `json:"html" gorm:"column:html"`
	ModifiedDate time.Time    `json:"modified_date"`
	Attachments  []EchoAttach `json:"attachments"`
}

// ErrEchoEmailNameNotSpecified is thrown when a template name is not specified
var ErrEchoEmailNameNotSpecified = errors.New("Echo of mine name not specified")

// ErrEchoEmailMissingParameter is thrown when a needed parameter is not provided
var ErrEchoEmailMissingParameter = errors.New("Need to specify at least plaintext or HTML content")

// Validate checks the given template to make sure values are appropriate and complete
func (t *EchoEmail) Validate() error {
	switch {
	case t.Name == "":
		return ErrEchoEmailNameNotSpecified
	case  t.HTML == "":
		return ErrEchoEmailMissingParameter
	}
	if err := ValidateEchoEmail(t.HTML); err != nil {
		return err
	}
	
	return nil
}

// GetEchoEmails returns the emails owned by the given user.
func GetEchoEmails(uid int64) ([]EchoEmail, error) {
	ts := []EchoEmail{}
	err := db.Where("user_id=?", uid).Find(&ts).Error
	if err != nil {
		log.Error(err)
		return ts, err
	}
	for i := range ts {
		// Get Echo email Attachment
		err = db.Where("EchoEmail_id=?", ts[i].Id).Find(&ts[i].Attachments).Error
		if err == nil && len(ts[i].Attachments) == 0 {
			ts[i].Attachments = make([]EchoAttach, 0)
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error(err)
			return ts, err
			
		}
	}
	return ts, err
}

// GetEchoEmail returns the template, if it exists, specified by the given id and user_id.
func GetEchoEmail(id int64, uid int64) (EchoEmail, error) {
	t := EchoEmail{}
	err := db.Where("user_id=? and id=?", uid, id).Find(&t).Error
	if err != nil {
		log.Error(err)
		return t, err
	}

	// Get Attachments
	err = db.Where("EchoEmail_id=?", t.Id).Find(&t.Attachments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		return t, err
	}
	if err == nil && len(t.Attachments) == 0 {
		t.Attachments = make([]EchoAttach, 0)
	}
	return t, err
}

// GetEchoEmailByName returns the feedback email, if it exists, specified by the given name and user_id.
func GetEchoEmailByName(n string, uid int64) (EchoEmail, error) {
	t := EchoEmail{}
	err := db.Where("user_id=? and name=?", uid, n).Find(&t).Error
	if err != nil {
		log.Error(err)
		return t, err
	}

	// Get Attachments
	err = db.Where("EchoEmail_id=?", t.Id).Find(&t.Attachments).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		return t, err
	}
	if err == nil && len(t.Attachments) == 0 {
		t.Attachments = make([]EchoAttach, 0)
	}
	return t, err
}

// PostEchoEmail creates a new email in the database.
func PostEchoEmail(t *EchoEmail ) error {
	// Insert into the DB
	if err := t.Validate(); err != nil {
		return err
	}
	err := db.Save(t).Error
	if err != nil {
		log.Error(err)
		return err
	}

	// Save every attachment
	for i := range t.Attachments {
		t.Attachments[i].EchoEmailId= t.Id
		err := db.Save(&t.Attachments[i]).Error
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// PutEchoEmail edits an existing email in the database.
// Per the PUT Method RFC, it presumes all data for a email is provided.
func PutEchoEmail(t *EchoEmail) error {
	if err := t.Validate(); err != nil {
		return err
	}
	// Delete all attachments, and replace with new ones
	err := db.Where("EchoEmail_id=?", t.Id).Delete(&EchoAttach{}).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error(err)
		return err
	}
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	for i := range t.Attachments {
		t.Attachments[i].EchoEmailId = t.Id
		err := db.Save(&t.Attachments[i]).Error
		if err != nil {
			log.Error(err)
			return err
		}
	}

	// Save final email
	err = db.Where("id=?", t.Id).Save(t).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// DeleteEchoEmail deletes an existing email in the database.
// An error is returned if an email with the given user id and echo email id is not found.
func DeleteEchoEmail(id int64, uid int64) error {
	// Delete attachments
	err := db.Where("EchoEmail_id=?", id).Delete(&EchoAttach{}).Error
	if err != nil {
		log.Error(err)
		return err
	}

	// Finally, delete the email itself
	err = db.Where("user_id=?", uid).Delete(EchoEmail{Id: id}).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}