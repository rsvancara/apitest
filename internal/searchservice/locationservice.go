package searchservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"rainier/internal/config"
)

type LocationService struct {
	hConfig *config.AppConfig
}

type LocationAutoComplete struct {
	State        string `json:"state"`
	StateCode    string `json:"state_code"`
	City         string `json:"city"`
	LocationType string `json:"location_type"`
}

type InterpretedLocation struct {
	LocationType string `json:"location_type"`
	StateCode    string `json:"state_code"`
	City         string `json:"city"`
	Zip          string `json:"zip"`
}

//CTXHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func (ls *LocationService) Initialize(hConfig *config.AppConfig) {
	ls.hConfig = hConfig
}

// Location Helper Service is used by auto complete to provide convenient lookup for users.  This function will flatten everything into a list of strings which works easier for
// client side javascript
func (ls *LocationService) LocationAutoComplete(location string) ([]string, error) {

	// Make zero length array, golang will also create "null" arrays if not initialized.  This creates wonky behaviour when value is serialized to json
	locations := make([]string, 0)

	LocationAutocompleteServiceURI := ls.hConfig.LocationHelper + "/v1/complete?prefix=" + strings.ReplaceAll(location, " ", "+")

	// Maybe look into putting this into a goroutine with a channel to make this "asynchronous" but I dont really see the reason if there is only one call.
	// Good practice includes adding a timeout for the request so it does not hang up routines for an extended period of time
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Get(LocationAutocompleteServiceURI)

	// Usual culprits here includes search layer is not avaible.
	if err != nil {
		return locations, fmt.Errorf("Error connecting to location helper service with URI [%s].  Is the service up? The error is: %s", LocationAutocompleteServiceURI, err)
	}

	// Response is nil if no response is returned, not always an error.  If response is nil, the ioutil.ReadAll will fail. Exception should be hit before!
	if response == nil {
		return locations, fmt.Errorf("Error connecting to location helper service with URI [%s].  Is the service up?", LocationAutocompleteServiceURI)
	}

	// Read the response body, covert it to byte array
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return locations, fmt.Errorf("Response data from service could not be read for URI [%s].  The error returned is: %s", LocationAutocompleteServiceURI, err)
	}

	// Hack to return if there are no results == {}
	if string(responseData) == "{}" {
		return locations, nil
	}

	var locationAutoComplete []LocationAutoComplete
	err = json.Unmarshal(responseData, &locationAutoComplete)
	if err != nil {
		return locations, fmt.Errorf("Error unmarshalling JSON response data for URI [%s].  Please confirm the response contains the correct data or has not changed.  The error returned is: %s", LocationAutocompleteServiceURI, err)
	}

	for _, val := range locationAutoComplete {
		if val.LocationType == "City" {
			locations = append(locations, val.City+", "+val.StateCode)
		}
		if val.LocationType == "State" {
			locations = append(locations, val.StateCode)
		}
	}

	return locations, nil

}

// Location Interpreter Service attempts to find a location and return some location data
func (ls *LocationService) LocationInterpreterService(searchLocation string) (InterpretedLocation, error) {

	var location InterpretedLocation

	locationInterpreterServiceURI := ls.hConfig.LocationHelper + "/v1/nonaddress?text=" + url.QueryEscape(searchLocation)

	// Maybe look into putting this into a goroutine with a channel to make this "asynchronous" but I dont really see the reason if there is only one call.
	// Good practice includes adding a timeout for the request so it does not hang up routines for an extended period of time
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Get(locationInterpreterServiceURI)

	// Usual culprits here includes search layer is not avaible.
	if err != nil {
		return location, fmt.Errorf("Error connecting to location helper service with URI [%s].  Is the service up? The error is: %s", locationInterpreterServiceURI, err)
	}

	// Response is nil if no response is returned, not always an error.  If response is nil, the ioutil.ReadAll will fail. Exception should be hit before!
	if response == nil {
		return location, fmt.Errorf("Error connecting to location helper service with URI [%s].  Is the service up?", locationInterpreterServiceURI)
	}

	// Read the response body, covert it to byte array
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return location, fmt.Errorf("Response data from service could not be read for URI [%s].  The error returned is: %s", locationInterpreterServiceURI, err)
	}

	// Hack to return if there are no results == {}
	if string(responseData) == "{}" {
		return location, nil
	}

	err = json.Unmarshal(responseData, &location)
	if err != nil {
		return location, fmt.Errorf("Error unmarshalling JSON response data for URI [%s].  Please confirm the response contains the correct data or has not changed.  The error returned is: %s", locationInterpreterServiceURI, err)
	}

	return location, nil

}
