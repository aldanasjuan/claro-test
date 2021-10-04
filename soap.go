package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"net/http"
)

const url = "http://56bf-190-148-157-118.ngrok.io:80/ValidarNitService/ValidarNitService"

const validarNit = `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
<Body>
	<validarNit xmlns="http://validadornitsoap.claro.com/">
		<nit xmlns="">%v</nit>
	</validarNit>
</Body>
</Envelope>`

const listarNit = `<Envelope xmlns="http://schemas.xmlsoap.org/soap/envelope/">
<Body>
	<listarNitValidos xmlns="http://validadornitsoap.claro.com/"/>
</Body>
</Envelope>`

func ValidateNIT(nit string) (*NitResponse, error) {
	payload := fmt.Sprintf(validarNit, nit) //unsafe string, should be escaped
	action := "urn:validarNit"
	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", action)
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result := new(NitXMLResponse)
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		return nil, err
	}
	return &result.Body.Response.Return, nil
}

type NitXMLResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName  xml.Name
		Response struct {
			XMLName xml.Name
			Return  NitResponse `xml:"return"`
		} `xml:"validarNitResponse"`
	} `xml:"Body"`
}

type NitResponse struct {
	XMLName    xml.Name `json:"xml_name,omitempty"`
	Age        int      `xml:"age,omitempty" json:"age,omitempty"`
	Birthdate  string   `xml:"birthdate,omitempty" json:"birthdate,omitempty"`
	Gender     string   `xml:"gender,omitempty" json:"gender,omitempty"`
	MiddleName string   `xml:"middle_name,omitempty" json:"middle_name,omitempty"`
	Name       string   `xml:"name,omitempty" json:"name,omitempty"`
	NIT        string   `xml:"nit,omitempty" json:"nit,omitempty"`
	Res        string   `xml:"res,omitempty" json:"res,omitempty"`
}

type NitRequest struct {
	NIT string `json:"nit,omitempty"`
}

/*
<age	age>0</age>
<birthdate	birthdate>No disponible</birthdate>
<gender	gender>No disponible</gender>
<middleName>No disponible</middleName>
<name>No disponible</name>
<nit>1972656-2</nit>
<res>correcto</res>
*/
