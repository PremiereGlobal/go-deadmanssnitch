package deadmanssnitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrorResponse represents the structure of a API error
type ErrorResponse struct {
	ErrorType   string       `json:"type"`
	ErrorString string       `json:"error"`
	Validations []Validation `json:"validations"`
}

// Validation contains the details of a API field validation error
type Validation struct {
	Attribute string `json:"attribute"`
	Message   string `json:"message"`
}

func (c *Client) do(method string, path string, body []byte) ([]byte, error) {
	request, err := http.NewRequest(method, fmt.Sprintf("%s/%s", c.apiBaseURL, path), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.SetBasicAuth(c.apiKey, "")

	resp, err := c.httpClient.Do(request)
	return c.checkResponse(resp, err)
}

func (c *Client) checkResponse(response *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, fmt.Errorf("Error calling the API endpoint: %v", err)
	}

	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading the API response: %v", err)
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		errorResponse := ErrorResponse{}
		err = json.Unmarshal(bodyBytes, &errorResponse)
		if err != nil {
			return nil, err
		}

		errorString := ""

		if response.StatusCode == http.StatusUnprocessableEntity {
			// Fields invalid
			for _, v := range errorResponse.Validations {
				errorString = errorString + fmt.Sprintf("%s: %s, ", v.Attribute, v.Message)
			}
		} else {
			// Generic error
			errorString = fmt.Sprintf("%s: %s", errorResponse.ErrorType, errorResponse.ErrorString)
		}

		return nil, fmt.Errorf("Error requesting %s %s, HTTP %d. %s", response.Request.Method, response.Request.URL, response.StatusCode, errorString)
	}

	return bodyBytes, nil
}
