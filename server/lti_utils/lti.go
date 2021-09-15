// Package lti_utils provides ...
package lti_utils

import (
	"encoding/xml"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/henrixapp/go-lti"
	"gorm.io/gorm"
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

type LTIConnector struct {
	DB        *gorm.DB
	TokenAuth *jwtauth.JWTAuth
}

// export later?
type LTIUserInfos struct {
	// maybe this gets to be a satori.UUID
	ID                 string
	UserName           string
	PersonGivenName    string
	PersonPrimaryEmail string
	PersonFamilyName   string
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

// LTILaunch Performs a search in database  for a course with the given ID and returns a invitation.
func (l *LTIConnector) LTILaunch(w http.ResponseWriter, r *http.Request) {
	// Create a new LTIToolProvider
	ltiRequest, err := lti.NewLTIToolProvider(r)
	if err != nil {
		log.Println("err:", err)
		return
	}
	// Validate LTI request
	valid, err := ltiRequest.ValidateRequest(
		os.Getenv("LTI_SECRET_KEY"),
		true,
		false,
		true,
		// dont transform the path, this is the identity func
		func(path string) string {
			return path
		},
	)

	var res LTIUserInfos
	res.ID = ltiRequest.LTIHeaders.UserId
	res.PersonFamilyName = ltiRequest.LTIHeaders.LISPersonFamilyName
	res.PersonGivenName = ltiRequest.LTIHeaders.LISPersonGivenName
	res.PersonPrimaryEmail = ltiRequest.LTIHeaders.LISPersonPrimaryEmail

	if valid {
		var userInfos LTIUserInfos
		l.DB.First(&userInfos, res.ID)
		if l.DB.Error != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if userInfos.ID == "" {
			l.DB.Create(&res)
		}

		jwtClaims := map[string]interface{}{"ID": res.ID}
		jwtauth.SetExpiryIn(jwtClaims, 60*time.Minute)
		_, tokenString, _ := l.TokenAuth.Encode(jwtClaims)
		jwtCookie := &http.Cookie{Name: "jwt", Value: tokenString, HttpOnly: false, Path: "/"}
		http.SetCookie(w, jwtCookie)

		http.Redirect(w, r, "https://"+r.Host+"/", http.StatusMovedPermanently)
	} else {
		log.Println(err)
		// Redirect to return URL
		//returnUrl, _ := ltiRequest.CreateReturnURL()

		http.Error(w, "Couldn't validate your request.", http.StatusInternalServerError)
	}
}
