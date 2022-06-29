package parser

type FlowData struct {
	Name            string     `yaml:"name"`
	Destination     string     `yaml:"destination"`
	Source          string     `yaml:"source"`
	Protocol        []Protocol `yaml:"protocol"`
	DestinationPort []int      `yaml:"destination-port"`
	SourcePort      []int      `yaml:"source-port"`
	Action          Action     `yaml:"action"`
}

type FlowSpecData struct {
	Flows []FlowData `yaml:"flows"`
}
