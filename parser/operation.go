package parser

import "encoding/xml"

type Operation int

const (
	Merge Operation = iota
	Replace
	Create
	Delete
	Remove
)

func (o Operation) String() string {
	return [...]string{
		"merge",
		"replace",
		"create",
		"delete",
		"remove",
	}[o]
}

func (o Operation) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: o.String(),
	}, nil
}
