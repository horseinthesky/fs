package device

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type DeviceData struct {
	Name string
	Site struct {
		Slug string
	}
	Role struct {
		Slug string
	}
	Tenant struct {
		Slug string
	}
}

type NetboxResponse struct {
	Results []DeviceData
}

type Device struct {
	Name   string
	Site   string
	Role   string
	Tenant string
}

func (n *NetboxResponse) getDevices() []*Device {
	devices := []*Device{}
	for _, deviceData := range n.Results {
		devices = append(devices, &Device{
			Name:   deviceData.Name,
			Site:   deviceData.Site.Slug,
			Role:   deviceData.Role.Slug,
			Tenant: deviceData.Tenant.Slug,
		})
	}

	return devices
}

func GetDeviceList(inventory string) ([]*Device, error) {
	response, err := http.Get(inventory)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	result := &NetboxResponse{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	devices := result.getDevices()

	// devices := []*Device{}
	// devices = append(devices, &Device{Name: "10.10.30.4"})
	// devices = append(devices, &Device{Name: "10.10.30.8"})

	return devices, nil
}
