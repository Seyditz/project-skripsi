package utils

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type OAIResponse struct {
	Records []Record `xml:"ListRecords>record>metadata>dc>title"`
}

type Record struct {
	Title string `xml:",chardata"`
}

// FetchTitles fetches titles from the OAI-PMH endpoint and returns them as a slice of strings
func FetchTitles() ([]string, error) {
	// URL OAI-PMH endpoint
	url := "http://repository.upnvj.ac.id/cgi/oai2?verb=ListRecords&metadataPrefix=oai_dc"

	// Make HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse XML response
	var oaiResponse OAIResponse
	err = xml.Unmarshal(body, &oaiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	// Extract titles
	var titles []string
	for _, record := range oaiResponse.Records {
		titles = append(titles, record.Title)
	}

	return titles, nil
}
