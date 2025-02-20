package models

import (
	"encoding/xml"

	"github.com/google/uuid"
)

//
//	Get
//

func NewGetSOAPHeader(wsaTo string) SendSOAPHeader {
	return SendSOAPHeader{
		WsaAction:    "http://schemas.xmlsoap.org/ws/2004/09/transfer/Get",
		WsaTo:        wsaTo,
		WsaMessageId: uuid.NewString(),
	}
}

func NewGetSOAPBody() SendSOAPBody {
	// A Get message has an empty body.
	return SendSOAPBody{}
}

//
//	Get Response
//

type GetResponseBody struct {
	XMLName     xml.Name
	MexMetadata MexMetadata `xml:"Metadata"`
}

type MexMetadata struct {
	XMLName             xml.Name
	MexMetadataSections []MexMetadataSection `xml:"MetadataSection"`
}

type MexMetadataSection struct {
	XMLName                   xml.Name
	MexMetadataSectionDialect string               `xml:"Dialect,attr"`
	DpwsThisModel             DpwsThisModel        `xml:"ThisModel"`
	DpwsThisDevice            DpwsThisDevice       `xml:"ThisDevice"`
	DpwsThisRelationship      DpwsThisRelationship `xml:"Relationship"`
	MexLocation               string               `xml:"Location"`
}

type DpwsThisModel struct {
	XMLName             xml.Name `xml:"ThisModel"`
	DpwsManufacturer    string   `xml:"Manufacturer"`
	DpwsManufacturerUrl string   `xml:"ManufacturerUrl"`
	DpwsModelName       string   `xml:"ModelName"`
	DpwsModelUrl        string   `xml:"ModelNumber"`
	DpwsPresentationUrl string   `xml:"PresentationUrl"`
}

type DpwsThisDevice struct {
	XMLName             xml.Name `xml:"ThisDevice"`
	DpwsFriendlyName    string   `xml:"FriendlyName"`
	DpwsFirmwareVersion string   `xml:"FirmwareVersion"`
	DpwsSerialNumber    string   `xml:"SerialNumber"`
}

type DpwsThisRelationship struct {
	XMLName    xml.Name   `xml:"Relationship"`
	DpwsHost   DpwsHost   `xml:"Host"`
	DpwsHosted DpwsHosted `xml:"Hosted"`
}

type DpwsHost struct {
	XMLName              xml.Name
	WsaEndpointReference ReceiveWsaEndpointReference
	DpwsTypes            string `xml:"Types"`
}

type DpwsHosted struct {
	XMLName              xml.Name
	WsaEndpointReference ReceiveWsaEndpointReference
	DpwsTypes            string `xml:"Types"`
	DpwsServiceID        string `xml:"ServiceId"`
}

// Getter of the Get Response Body

func (g *GetResponseBody) GetMexMetadataSections() []MexMetadataSection {
	return g.MexMetadata.MexMetadataSections
}

func (g *MexMetadataSection) GetMexLocation() string {
	return g.MexLocation
}

func (g *GetResponseBody) GetAllMexLocation() []string {
	retVal := []string{}

	for _, section := range g.GetMexMetadataSections() {
		retVal = append(retVal, section.GetMexLocation())
	}

	return retVal
}
