package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/api/getJson", getJson)
	r.Run(":8080")
}

func getJson(c *gin.Context) {
	// SOAP request reqBody
	tplXml := `<?xml version="1.0" encoding="utf-8"?>
	<soap12:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
	  <soap12:Body>
		<TransData_Json xmlns="http://tempuri.org/">
		  <company>test</company>
		  <password>test1</password>
		  <json>%s</json>
		</TransData_Json>
	  </soap12:Body>
	</soap12:Envelope>`

	// SOAP endpoint URL
	endpointURL := "https://hctrt.hct.com.tw/EDI_WebService2/Service1.asmx"
	reqJson, _ := io.ReadAll(c.Request.Body)
	reqXml := fmt.Sprintf(tplXml, string(reqJson))

	// Create a new HTTP request
	req, err := http.NewRequest("POST", endpointURL, strings.NewReader(reqXml))
	if err != nil {
		log.Fatal(err)
	}

	// Set the SOAP headers and content type
	req.Header.Set("Content-Type", "application/soap+xml; charset=utf-8")

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the SOAP response into a struct
	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Text    string   `xml:",chardata"`
		Soap    string   `xml:"soap,attr"`
		Xsi     string   `xml:"xsi,attr"`
		Xsd     string   `xml:"xsd,attr"`
		Body    struct {
			Text                  string `xml:",chardata"`
			TransDataJsonResponse struct {
				Text                string `xml:",chardata"`
				Xmlns               string `xml:"xmlns,attr"`
				TransDataJsonResult string `xml:"TransData_JsonResult"`
			} `xml:"TransData_JsonResponse"`
		} `xml:"Body"`
	}

	type JsonResult []struct {
		Num        string `json:"Num"`
		Success    string `json:"success"`
		Edelno     string `json:"edelno"`
		Epino      string `json:"epino"`
		Erstno     string `json:"erstno"`
		Eqamt      string `json:"eqamt"`
		Image      any    `json:"image"`
		ErrMsg     string `json:"ErrMsg"`
		NewOutArea string `json:"NewOutArea"`
		Eqmny      string `json:"eqmny"`
		Code1      string `json:"CODE1"`
		Code2      string `json:"CODE2"`
		Code3      string `json:"CODE3"`
		Code4      string `json:"CODE4"`
		Code5      string `json:"CODE5"`
		Code7      string `json:"CODE7"`
		Areas      string `json:"AREAS"`
		Mdcode1    string `json:"MDCODE1"`
		Mdcode2    string `json:"MDCODE2"`
		Mdcode3    string `json:"MDCODE3"`
	}

	var soapResponse Envelope
	err = xml.Unmarshal(body, &soapResponse)
	if err != nil {
		log.Fatal(err)
	}

	j := JsonResult{}

	json.Unmarshal([]byte(soapResponse.Body.TransDataJsonResponse.TransDataJsonResult), &j)
	// Check for errors
	if soapResponse.Body.TransDataJsonResponse.TransDataJsonResult == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving Json"})
		return
	}

	// Return the Json from the SOAP response
	c.IndentedJSON(http.StatusOK, j)
}
