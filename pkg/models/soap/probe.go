package models

import (
	"encoding/xml"

	"github.com/google/uuid"
)

//
// Probe
//

func NewProbeSOAPHeader() SendSOAPHeader {
	return SendSOAPHeader{
		WsaAction:    "http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01/Probe",
		WsaTo:        "urn:docs-oasis-open-org:ws-dd:ns:discovery:2009:01",
		WsaMessageId: uuid.NewString(),
	}
}

func NewProbeSOAPBody() SendSOAPBody {
	var probePayload struct {
		XMLName   xml.Name `xml:"wsd:Probe"`
		WsdType   string   `xml:"wsd:Types"`
		WsdScopes string   `xml:"wsd:Scopes"`
	}

	// this data is taken from the sdcX implementation
	// TODO: allow customization
	probePayload.WsdType = "dpws:Device mdpws:MedicalDevice"
	probePayload.WsdScopes = "sdc.cdc.type:///130535 sdc.ctxt.loc:/sdc.ctxt.loc.detail/DWHL%2F%2F%2FF05%2F%2FTKl?fac=DWHL&amp;poc=F05&amp;bed=TKl"

	return SendSOAPBody{Payload: probePayload}
}

//
//	Probe
//

type ProbeBody struct {
	XMLName   xml.Name `xml:"Probe"`
	WsdType   string   `xml:"Types"`
	WsdScopes string   `xml:"Scopes"`
}

//
//	Probe Match
//

type ProbeMatchBody struct {
	XMLName         xml.Name
	WsdProbeMatches WsdProbeMatches `xml:"ProbeMatches"`
}

type WsdProbeMatches struct {
	XMLName       xml.Name
	WsdProbeMatch WsdProbeMatch `xml:"ProbeMatch"`
}

type WsdProbeMatch struct {
	XMLName              xml.Name
	WsaEndpointReference ReceiveWsaEndpointReference
	WsdTypes             string `xml:"Types"`
	WsdScopes            string `xml:"Scopes"`
	WsdXAddrs            string `xml:"XAddrs"`
	WsdMetadataVersion   string `xml:"MetadataVersion"`
}

// Getter of ProbeMatchs Body

func (p *ProbeMatchBody) GetXAddrs() string {
	return p.WsdProbeMatches.WsdProbeMatch.WsdXAddrs
}

func (p *ProbeMatchBody) GetTypes() string {
	return p.WsdProbeMatches.WsdProbeMatch.WsdTypes
}

func (p *ProbeMatchBody) GetScopes() string {
	return p.WsdProbeMatches.WsdProbeMatch.WsdScopes
}

func (p *ProbeMatchBody) GetMetadataVersion() string {
	return p.WsdProbeMatches.WsdProbeMatch.WsdMetadataVersion
}

func (p *ProbeMatchBody) GetEndpointReferenceAddress() string {
	return p.WsdProbeMatches.WsdProbeMatch.WsaEndpointReference.WsaAddress
}
