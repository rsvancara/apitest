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

type PersonSearchService struct {
	hConfig *config.AppConfig
}

type PersonResults struct {
	ResultCount int       `json:"result_count"`
	Results     []Results `json:"results"`
}

type Results struct {
	ID           string `json:"id"`
	IsSuppressed bool   `json:"is_suppressed"`
	Fn           string `json:"fn"`
	Mn           string `json:"mn,omitempty"`
	Ln           string `json:"ln"`
	Uri          string
	Age          int
	Akas         []struct {
		Fn string `json:"fn"`
		Ln string `json:"ln"`
	} `json:"akas"`
	Dob struct {
		Y int `json:"y"`
		M int `json:"m"`
		D int `json:"d"`
	} `json:"dob,omitempty"`
	IsDead               bool                `json:"is_dead"`
	CurrentLocations     []CurrentLocations  `json:"current_locations"`
	HistoricLocations    []HistoricLocations `json:"historic_locations"`
	MobilePhonesCount    int                 `json:"mobile_phones_count"`
	NonmobilePhonesCount int                 `json:"nonmobile_phones_count"`
	FillScore            int                 `json:"fill_score"`
	Relatives            []Name              `json:"relatives"`
	Associates           []Name              `json:"associates"`
	Suffix               string              `json:"suffix,omitempty"`
}

type HistoricLocations struct {
	City  string `json:"city"`
	State string `json:"state"`
}

type CurrentLocations struct {
	City     string `json:"city"`
	State    string `json:"state"`
	Distance int    `json:"distance"`
	Nbh      string `json:"nbh"`
	MacroNbh string `json:"macro_nbh"`
}

type Name struct {
	Fn string `json:"fn"`
	Mn string `json:"mn"`
	Ln string `json:"ln"`
}

//CTXHandlerContext constructs a new HandlerContext,
//ensuring that the dependencies are valid values
func (ss *PersonSearchService) Initialize(hConfig *config.AppConfig) {
	ss.hConfig = hConfig
}

// Search by name
func (ss *PersonSearchService) GetPersonResultsByName(firstname string, lastname string, isFuzzy bool) (PersonResults, error) {

	var personResults PersonResults

	//log.Info().Msgf("First Name %s - Last Name %s", firstname, lastname)

	fuzzy := "false"
	if isFuzzy == true {
		fuzzy = "true"
	}

	personSearchServiceURI := ss.hConfig.SearchLayer + "/v1/search/people?fn=" + url.QueryEscape(firstname) + "&ln=" + url.QueryEscape(lastname) + "&city=&state=&fs=" + fuzzy

	//log.Info().Msgf("Searching...%s", personSearchServiceURI)

	// Maybe look into putting this into a goroutine with a channel to make this "asynchronous" but I dont really see the reason if there is only one call.
	// Good practice includes adding a timeout for the request so it does not hang up routines for an extended period of time
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Get(personSearchServiceURI)

	// Usual culprits here includes search layer is not avaible.
	if err != nil {
		return personResults, fmt.Errorf("Error connecting to search layer with URI [%s].  Is the service up? The error is: %s", personSearchServiceURI, err)
	}

	// Response is nil if no response is returned, not always an error.  If response is nil, the ioutil.ReadAll will fail. Exception should be hit before!
	if response == nil {
		return personResults, fmt.Errorf("Error connecting to search layer with URI [%s].  Is the service up?", personSearchServiceURI)
	}

	// Read the response body, covert it to byte array
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return personResults, fmt.Errorf("Response data from service could not be read for URI [%s].  The error returned is: %s", personSearchServiceURI, err)
	}

	// TODO: Remove
	//fmt.Println(string(responseData))

	err = json.Unmarshal(responseData, &personResults)
	if err != nil {
		return personResults, fmt.Errorf("Error unmarshalling JSON responsed data for URI [%s].  Please confirm the response contains the correct data or has not changed.  The error returned is: %s", personSearchServiceURI, err)
	}

	// Populate the URI and other bits into the PersonResults
	PersonDataBuilder(&personResults)

	return personResults, nil
}

// Search by name, city and state
func (ss *PersonSearchService) GetPersonResultsByNameLocation(firstname string, lastname string, city string, state string, isFuzzy bool) (PersonResults, error) {
	var personResults PersonResults

	fuzzy := "false"
	if isFuzzy == true {
		fuzzy = "true"
	}

	personSearchServiceURI := ss.hConfig.SearchLayer + "/v1/search/people?fn=" + url.QueryEscape(firstname) + "&ln=" + url.QueryEscape(lastname) + "&city=" + url.QueryEscape(city) + "&state=" + url.QueryEscape(state) + "&fs=" + fuzzy

	//log.Info().Msgf("Searching...%s", personSearchServiceURI)

	// Maybe look into putting this into a goroutine with a channel to make this "asynchronous" but I dont really see the reason if there is only one call.
	// Good practice includes adding a timeout for the request so it does not hang up routines for an extended period of time
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	response, err := client.Get(personSearchServiceURI)

	// Usual culprits here includes search layer is not avaible.
	if err != nil {
		return personResults, fmt.Errorf("Error connecting to search layer with URI [%s].  Is the service up? The error is: %s", personSearchServiceURI, err)
	}

	// Response is nil if no response is returned, not always an error.  If response is nil, the ioutil.ReadAll will fail. Exception should be hit before!
	if response == nil {
		return personResults, fmt.Errorf("Error connecting to search layer with URI [%s].  Is the service up?", personSearchServiceURI)
	}

	// Read the response body, covert it to byte array
	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return personResults, fmt.Errorf("Response data from service could not be read for URI [%s].  The error returned is: %s", personSearchServiceURI, err)
	}

	fmt.Println(string(responseData))
	err = json.Unmarshal(responseData, &personResults)
	if err != nil {
		return personResults, fmt.Errorf("Error unmarshalling JSON responsed data for URI [%s].  Please confirm the response contains the correct data or has not changed.  The error returned is: %s", personSearchServiceURI, err)
	}

	// Populate the URI and other bits into the PersonResults
	PersonDataBuilder(&personResults)

	return personResults, nil
}

// Person Details Results
func (ss *PersonSearchService) GetPersonDetails(firstname string, lastname string, city string, state string, isFuzzy bool) (Results, error) {
	var result Results

	return result, nil
}

// Enhances person results with additions details
func PersonDataBuilder(pr *PersonResults) {

	for i, p := range pr.Results {

		// Add user URI for quick reference in templating
		currentCity := ""
		currentState := ""
		for _, l := range p.CurrentLocations {
			currentCity = l.City
			currentState = l.State
			break
		}

		// Dereference the pointer here with the index and then set the values...a nuance of go.
		r := &pr.Results[i]

		r.Uri = NameLocationIDURIGenerator(p.Fn, p.Ln, currentCity, currentState, p.ID)

		r.Age = getAge(p.Dob.M, p.Dob.D, p.Dob.Y)

	}
}

// Generates URI based on Name Location ID
func NameLocationIDURIGenerator(fname string, lname, city string, state string, id string) string {

	return fmt.Sprintf("/name/%s-%s/%s-%s/%s", fname, lname, city, state, id)
}

// Determines the approximate age to display on details
func getAge(month int, day int, year int) int {

	dob := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return age(dob)
}

// AgeAt gets the age of an entity at a certain time.
func ageAt(birthDate time.Time, now time.Time) int {
	// Get the year number change since the player's birth.
	years := now.Year() - birthDate.Year()

	// If the date is before the date of birth, then not that many years have elapsed.
	birthDay := getAdjustedBirthDay(birthDate, now)
	if now.YearDay() < birthDay {
		years -= 1
	}

	return years
}

// Age is shorthand for AgeAt(birthDate, time.Now()), and carries the same usage and limitations.
func age(birthDate time.Time) int {
	return ageAt(birthDate, time.Now())
}

// Gets the adjusted date of birth to work around leap year differences.
func getAdjustedBirthDay(birthDate time.Time, now time.Time) int {
	birthDay := birthDate.YearDay()
	currentDay := now.YearDay()
	if isLeap(birthDate) && !isLeap(now) && birthDay >= 60 {
		return birthDay - 1
	}
	if isLeap(now) && !isLeap(birthDate) && currentDay >= 60 {
		return birthDay + 1
	}
	return birthDay
}

// Works out if a time.Time is in a leap year.
func isLeap(date time.Time) bool {
	year := date.Year()
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}
	return false
}
