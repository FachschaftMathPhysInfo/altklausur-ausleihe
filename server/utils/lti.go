// Package utils provides ...
package utils

import (
	"encoding/xml"
	"net/http"
)

type Cartridge struct {
	Title           string `xml:"blti:title"`
	Description     string `xml:"blti:description"`
	LaunchUrl       string `xml:"blti:launch_url"`
	SecureLaunchUrl string `xml:"blti:secure_launch_url"`
	Icon            string `xml:"blti:icon"`
}
type LTIConfig struct {
	Cartridge Cartridge `xml:"cartridge_basiclti_link"`
	Blti      string    `xml:"xmlns:blti,attr"`
}

func LTIConfigHandler(w http.ResponseWriter, r *http.Request) {
	x, err := xml.Marshal(
		LTIConfig{
			Blti: "http://www.imsglobal.org/xsd/imsbasiclti_v1p0",
			Cartridge: Cartridge{
				Title:           "AltklausurAusleihe",
				Description:     "AltklausurAusleihe",
				LaunchUrl:       "https://" + r.Host + "/distributor/lti_launch",
				SecureLaunchUrl: "https://" + r.Host + "/distributor/lti_launch",
				// FIXME (christian): add a logo
				// Icon:            "https://" + r.Host + "/logo192.png",
			},
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}
