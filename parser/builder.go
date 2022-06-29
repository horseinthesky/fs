package parser

import (
	"encoding/xml"
	"regexp"
)

type Match struct {
	XMLName         xml.Name   `xml:"match"`
	Protocol        []Protocol `xml:"protocol,omitempty"`
	Source          string     `xml:"source,omitempty"`
	Destination     string     `xml:"destination,omitempty"`
	SourcePort      []int      `xml:"source-port,omitempty"`
	DestinationPort []int      `xml:"destination-port,omitempty"`
}

type Route struct {
	Name  string `xml:"name"`
	Match Match
	Then  Action `xml:"then"`
}

type Config struct {
	XMLName xml.Name `xml:"config"`
	Flow    struct {
		XMLName   xml.Name  `xml:"flow"`
		Operation Operation `xml:"operation,attr"`
		Routes    []*Route  `xml:"route"`
		TermOrder TermOrder `xml:"term-order"`
	} `xml:"configuration>routing-options>flow"`
}

func normalizeXML(xml []byte) string {
	m := regexp.MustCompile("&lt;(.*)&gt;")
	return m.ReplaceAllString(string(xml), "<$1>")
}

func BuildConfig(data *FlowSpecData) (string, error) {
	cfg := &Config{}
	cfg.Flow.Operation = Replace

	var routes []*Route
	for _, flow := range data.Flows {
		routes = append(routes, &Route{
			Name: flow.Name,
			Match: Match{
				Protocol:        flow.Protocol,
				Destination:     flow.Destination,
				Source:          flow.Source,
				DestinationPort: flow.DestinationPort,
				SourcePort:      flow.SourcePort,
			},
			Then: flow.Action,
		},
		)
	}
	cfg.Flow.Routes = routes

	configXML, err := xml.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return "", err
	}

	fixedConfigXML := normalizeXML(configXML)

	return fixedConfigXML, nil
}
