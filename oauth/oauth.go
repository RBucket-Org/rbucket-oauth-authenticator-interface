package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RBucket-Org/rbucket-oauth-authenticator-interface/oauth/rest_errors"
	"github.com/mercadolibre/golang-restclient/rest"
)

/*
This package contains the authenticator methods where the given access token is get
 authenticated and the request is passed in the form of the bool value to get the
auth situation
*/

//SET the constant public, client-id, caller-id, access-token
const (
	headerXPublic   = "X-Public"
	headerXClientID = "X-Client-Id"
	headerXCallerID = "X-Caller-Id"

	paramAccessToken = "access_token"
)

//setting the default URL to call the accessToken for verification
var (
	oauthRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 200 * time.Millisecond,
	}
)

//accessToken
type accessToken struct {
	ID       string `json:"id"`
	UserID   int64  `json:"user_id"`
	ClientID int64  `json:"client_id"`
}

// IsPublic : Check whether the request is public or private
func IsPublic(request *http.Request) bool {
	if request == nil {
		return true
	}
	//when condition is true then it is public
	return request.Header.Get(headerXPublic) == "true"
}

//GetCallerID : this method the fetches the callerID or UserID
func GetCallerID(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	//fetch the callerID from the request header
	callerID := request.Header.Get(headerXCallerID)
	//convert it into the int64
	userID, err := strconv.ParseInt(callerID, 10, 64)
	if err != nil {
		return 0
	}

	return userID //or CallerID
}

//GetClientID : this method fetches the clientID
func GetClientID(request *http.Request) int64 {
	if request == nil {
		return 0
	}

	//fetch the clientID from the request header
	id := request.Header.Get(headerXClientID)
	//convert it into the int64
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}

	return clientID
}

//clean the request header
func cleanRequestHeader(request *http.Request) {
	if request == nil {
		return
	}

	//delete the header
	request.Header.Del(headerXCallerID)
	request.Header.Del(headerXClientID)
}

//get the accessToken from the accessToken Database
func getAccessToken(accessTokenID string) (*accessToken, rest_errors.RestError) {
	response := oauthRestClient.Get(fmt.Sprintf("/oauth/access_token/%s", accessTokenID))
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalServerError("invalid rest client response when trying to get the access token")
	}

	//response status code
	if response.StatusCode > 299 {
		restErr, err := rest_errors.NewRestErrorFromBytes(response.Bytes())
		if err != nil {
			return nil, rest_errors.NewInternalServerError("cannot unmarshal the given response")
		}
		return nil, restErr
	}

	//if there is not accessToken then return the accessToken
	var at accessToken
	if err := json.Unmarshal(response.Bytes(), &at); err != nil {
		return nil, rest_errors.NewInternalServerError("cannot the unmarshal the response result")
	}

	return &at, nil
}

//AuthenticateRequest : this method clean the request header and authenticate the restAPI
func AuthenticateRequest(request *http.Request) rest_errors.RestError {
	if request == nil {
		return rest_errors.NewBadRequestError("empty request")
	}

	//remove all the header file
	cleanRequestHeader(request)

	//get the accessToken
	accessTokenID := strings.TrimSpace(request.URL.Query().Get(paramAccessToken))
	if accessTokenID == "" {
		return rest_errors.NewBadRequestError("Invalid tokenID")
	}

	//get the accessToken INFO by passing the accessTokenID
	accessToken, err := getAccessToken(accessTokenID)
	if err != nil {
		if err.Status() == http.StatusNotFound {
			return rest_errors.NewBadRequestError("accessToken info not found")
		}
		return err
	}

	//add the header to the given request
	request.Header.Add(headerXClientID, fmt.Sprintf("%v", accessToken.ClientID))
	request.Header.Add(headerXCallerID, fmt.Sprintf("%v", accessToken.UserID))
	return nil
}
