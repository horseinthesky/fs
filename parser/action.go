package parser

import (
	"encoding/xml"
	"errors"
)

type Action int

const (
	Discard Action = iota
	Accept
)

func (a Action) String() string {
	return [...]string{
		"discard",
		"accept",
	}[a]
}

func (a *Action) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var action string

	if err := unmarshal(&action); err != nil {
		return err
	}

	switch action {
	case "accept":
		*a = Accept
	case "discard":
		*a = Discard
	default:
		return errors.New("unknown action: " + action)
	}

	return nil
}

func (a Action) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement("<"+a.String()+"/>", start)
}
