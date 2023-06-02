package searchservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"rainier/internal/config"
)

type NameSplitterService struct {
	hConfig *config.AppConfig
}

type SplitName struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

//CTXHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func (ns *NameSplitterService) Initialize(hConfig *config.AppConfig) {
	ns.hConfig = hConfig
}

// Takes URL encoded string that maybe has a name in it and attempts to decode it
func (ls *NameSplitterService) FindName(nameSearch string) (SplitName, error) {

	var splitName SplitName

	//Name Splitter service will split the url encoded string, no need to figure anything out here

	nameSpliterServiceURI := ls.hConfig.NameHelper + "/v1/name?name=" + url.QueryEscape(nameSearch)

	// Maybe look into putting this into a goroutine with a channel to make this "asynchronous" but I dont really see the reason if there is only one call.
	// Good practice includes adding a timeout for the request so it does not hang up routines for an extended period of time
	client := http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Get(nameSpliterServiceURI)

	// Usual culprits here includes search layer is not avaible.
	if err != nil {
		return splitName, fmt.Errorf("Error connecting to name splitter service with URI [%s].  Is the service up? The error is: %s", nameSpliterServiceURI, err)
	}

	// Response is nil if no response is returned, not always an error.  If response is nil, the ioutil.ReadAll will fail. Exception should be hit before!
	if response == nil {
		return splitName, fmt.Errorf("Error connecting to name splitter service with URI [%s].  Is the service up?", nameSpliterServiceURI)
	}

	// Read the response body, covert it to byte array
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return splitName, fmt.Errorf("Response data from service could not be read for URI [%s].  The error returned is: %s", nameSpliterServiceURI, err)
	}

	err = json.Unmarshal(responseData, &splitName)
	if err != nil {
		return splitName, fmt.Errorf("Error unmarshalling JSON responsed data for URI [%s].  Please confirm the response contains the correct data or has not changed.  The error returned is: %s", nameSpliterServiceURI, err)
	}

	return splitName, nil

}
