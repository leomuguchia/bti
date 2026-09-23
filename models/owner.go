package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Owner struct {
	FullName         string `json:"full_name"`
	NationalID       string `json:"national_id"`
	PhoneNumber      string `json:"phone_number"`
	PhysicalAddress  string `json:"physical_address"`
	KRAPin           string `json:"kra_pin,omitempty"`
	IDPhotoRef       string `json:"id_photo_ref"`
	PassportPhotoRef string `json:"passport_photo_ref"`
}

func (o Owner) Validate() error {
	if o.FullName == "" {
		return errors.New("full name is required")
	}
	if o.NationalID == "" {
		return errors.New("national ID is required")
	}
	if o.PhoneNumber == "" {
		return errors.New("phone number is required")
	}
	if o.PhysicalAddress == "" {
		return errors.New("physical address is required")
	}
	if o.IDPhotoRef == "" {
		return errors.New("ID photo is required")
	}
	if o.PassportPhotoRef == "" {
		return errors.New("passport photo is required")
	}
	return nil
}

// Value lets database/sql write an Owner as JSON text.
func (o Owner) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}

// Scan lets database/sql read an Owner back out of the same column.
func (o *Owner) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("owner: unsupported scan type %T", src)
	}
	return json.Unmarshal(b, o)
}
