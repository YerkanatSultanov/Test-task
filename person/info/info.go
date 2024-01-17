package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"test-task/person/entity"
)

func getAge(name string) (int, error) {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("error in get request: %s", err)
		return 0, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}

	var ageResponse struct {
		Age int `json:"age"`
	}
	if err := json.Unmarshal(body, &ageResponse); err != nil {
		return 0, err
	}
	return ageResponse.Age, nil
}

func getGender(name string) (string, error) {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	response, err := http.Get(url)
	if err != nil {
		log.Printf("error in get request: %s", err)
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var genderResponse struct {
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(body, &genderResponse); err != nil {
		return "", err
	}
	return genderResponse.Gender, nil
}

func getNationalityFromNationalize(name string) (string, error) {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var nationalityResponse map[string]interface{}
	err = json.Unmarshal(body, &nationalityResponse)
	if err != nil {
		return "", err
	}

	countries, ok := nationalityResponse["country"].([]interface{})
	if !ok || len(countries) == 0 {
		return "", fmt.Errorf("failed to extract nationality from Nationalize response")
	}

	sort.Slice(countries, func(i, j int) bool {
		probI, _ := countries[i].(map[string]interface{})["probability"].(float64)
		probJ, _ := countries[j].(map[string]interface{})["probability"].(float64)
		return probI > probJ
	})

	country, ok := countries[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("failed to extract nationality from Nationalize response")
	}

	nationality, ok := country["country_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract nationality from Nationalize response")
	}

	return nationality, nil
}

func GetAllAdditionalInfo(name string) (*entity.PersonInfo, error) {
	age, err := getAge(name)
	if err != nil {
		return &entity.PersonInfo{}, err
	}

	gender, err := getGender(name)
	if err != nil {
		return &entity.PersonInfo{}, err
	}

	nationality, err := getNationalityFromNationalize(name)
	if err != nil {
		return &entity.PersonInfo{}, err
	}

	personInfo := &entity.PersonInfo{
		Age:         age,
		Gender:      gender,
		Nationality: nationality,
	}

	return personInfo, nil
}
