// Package lti_utils provides ...
package lti_utils

import (
	"encoding/json"
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
	PersonFullName     string
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
				Icon:            "https://" + r.Host + "/favicon.png",
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

	userInfoFromRequest := LTIUserInfos{
		ID:                 ltiRequest.LTIHeaders.UserId,
		PersonFamilyName:   ltiRequest.LTIHeaders.LISPersonFamilyName,
		PersonGivenName:    ltiRequest.LTIHeaders.LISPersonGivenName,
		PersonPrimaryEmail: ltiRequest.LTIHeaders.LISPersonPrimaryEmail,
		PersonFullName:     ltiRequest.LTIHeaders.LISPersonFullName,
	}

	if valid {
		// The JWT token carries all the user information!
		userInfosJSON, _ := json.Marshal(userInfoFromRequest)

		// Create the JWT Token for the User so he can access our application
		jwtClaims := map[string]interface{}{"user": string(userInfosJSON)}
		jwtauth.SetExpiryIn(jwtClaims, time.Hour)
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

// DummyLTILaunch just returns an invitation
// You can obtain a JWT Token by executing the following command
// curl -X POST http://localhost:8081/distributor/lti_launch -I
func (l *LTIConnector) DummyLTILaunch(w http.ResponseWriter, r *http.Request) {
	userInfoFromRequest := LTIUserInfos{
		ID:                 "200",
		PersonFamilyName:   "Testerson",
		PersonGivenName:    "Test",
		PersonPrimaryEmail: "test@example.com",
		PersonFullName:     "Test Testerson",
	}

	// The JWT token carries all the user information!
	userInfosJSON, err := json.Marshal(userInfoFromRequest)
	if err != nil {
		log.Println(err)
	}

	// Create the JWT Token for the User so he can access our application
	jwtClaims := map[string]interface{}{"user": string(userInfosJSON)}
	jwtauth.SetExpiryIn(jwtClaims, time.Hour)
	_, tokenString, err := l.TokenAuth.Encode(jwtClaims)
	if err != nil {
		log.Println(err)
	}
	jwtCookie := &http.Cookie{Name: "jwt", Value: tokenString, HttpOnly: false, Path: "/"}
	http.SetCookie(w, jwtCookie)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
