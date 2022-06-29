package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
)

type Protocol int

const (
	TCP Protocol = iota
	UDP
	ICMP
)

func (p Protocol) String() string {
	return [...]string{
		"tcp",
		"udp",
		"icmp",
	}[p]
}

func (p *Protocol) UnmarshalYAML(unmarshal func(v interface{}) error) error {
	var protocol string

	if err := unmarshal(&protocol); err != nil {
		return err
	}

	switch protocol {
	case "tcp":
		*p = TCP
	case "udp":
		*p = UDP
	case "icmp":
		*p = ICMP
	default:
		return errors.New("unknown protocol: " + protocol)
	}

	return nil
}

func (p Protocol) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string

	switch p {
	case 0:
		s = "tcp"
	case 1:
		s = "udp"
	case 2:
		s = "icmp"
	default:
		return errors.New("unknown action int: " + fmt.Sprint(p))
	}

	return e.EncodeElement(s, start)
}
