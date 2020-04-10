package nmap

type Status struct {
	State string `xml:"state,attr"`
}

type State struct {
	State string `xml:"state,attr"`
}

type Service struct {
	Name string `xml:"name,attr"`
}

type Port struct {
	PortId  int     `xml:"portid,attr"`
	State   State   `xml:"state"`
	Service Service `xml:"service"`
}

type Host struct {
	Status Status `xml:"status"`
	Ports  []Port `xml:"ports>port"`
}

type NmapExecution struct {
	Hosts []Host `xml:"host"`
}
