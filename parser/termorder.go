package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type TermOrder int

const (
	Standard TermOrder = iota
	Legacy
)

func (t TermOrder) String() string {
	return [...]string{
		"standard",
		"legacy",
	}[t]
}


func (t TermOrder) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string

	switch t {
	case 0:
		s = "standard"
	case 1:
		s = "legacy"
	default:
		return errors.New("unknown action int: " + fmt.Sprint(t))
	}

	return e.EncodeElement(s, start)
}
