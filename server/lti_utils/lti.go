// Package lti_utils provides ...
package lti_utils

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
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

type AuthHelper struct {
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
	IsAdmin            bool
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
func (l *AuthHelper) LTILaunch(w http.ResponseWriter, r *http.Request) {
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
		IsAdmin:            false,
	}

	if valid {
		// The JWT token carries all the user information!
		jwtCookie, err := l.createJWTCookie(&userInfoFromRequest)
		if err != nil {
			log.Println(err)
		}
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
func (l *AuthHelper) DummyLTILaunch(w http.ResponseWriter, r *http.Request) {
	userInfoFromRequest := LTIUserInfos{
		ID:                 "200",
		PersonFamilyName:   "Testerson",
		PersonGivenName:    "Test",
		PersonPrimaryEmail: "test@example.com",
		PersonFullName:     "Test Testerson",
	}

	// The JWT token carries all the user information!
	jwtCookie, err := l.createJWTCookie(&userInfoFromRequest)
	if err != nil {
		log.Println(err)
	}
	http.SetCookie(w, jwtCookie)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func (l *AuthHelper) AdminLoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// get the user infos from the JWT Token
	userInfos, err := getUserInfos(&ctx)
	if err != nil {
		log.Println(err)
	}

	// enable the user to be admin
	userInfos.IsAdmin = true

	log.Println("User", userInfos, "just authenticated as admin!")

	// re-encode the JWT Token with the secret
	jwtCookie, err := l.createJWTCookie(userInfos)
	if err != nil {
		log.Println(err)
	}

	http.SetCookie(w, jwtCookie)
	w.Write([]byte(jwtCookie.Value))
}

func getUserInfos(ctxPtr *context.Context) (*LTIUserInfos, error) {
	if ctxPtr == nil {
		return nil, fmt.Errorf("WTF, how is the context for this request nil")
	}

	_, claims, err := jwtauth.FromContext(*ctxPtr)
	if err != nil {
		return nil, err
	}

	var userInfos LTIUserInfos
	err = json.Unmarshal([]byte(claims["user"].(string)), &userInfos)
	if err != nil {
		return nil, err
	}

	return &userInfos, nil
}

func (l *AuthHelper) createJWTCookie(ltiUser *LTIUserInfos) (*http.Cookie, error) {
	userInfosJSON, err := json.Marshal(ltiUser)
	if err != nil {
		return nil, err
	}

	// Create the JWT Token for the User so he can access our application
	jwtClaims := map[string]interface{}{"user": string(userInfosJSON)}
	jwtauth.SetExpiryIn(jwtClaims, time.Hour)
	_, tokenString, err := l.TokenAuth.Encode(jwtClaims)
	if err != nil {
		return nil, err
	}
	jwtCookie := &http.Cookie{Name: "jwt", Value: tokenString, HttpOnly: false, Path: "/"}
	return jwtCookie, nil
}
