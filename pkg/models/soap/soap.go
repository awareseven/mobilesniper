package models

import (
	"encoding/xml"
	"fmt"
	"net"

	"github.com/google/uuid"
)

func NewSOAPMessage(header SendSOAPHeader, body SendSOAPBody) SendSOAPMessage {
	return SendSOAPMessage{
		// default SDC standard 11073 XML attributes
		XMLNsSoap12: "http://www.w3.org/2003/05/soap-envelope",
		XMLNsWsa:    "http://www.w3.org/2005/08/addressing",
		XMLNsWsd:    "http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01",
		XMLNsWse:    "http://schemas.xmlsoap.org/ws/2004/08/eventing",
		XMLNsMdpws:  "http://standards.ieee.org/downloads/11073/11073-20702-2016",
		XMLNsDpws:   "http://docs.oasis-open.org/ws-dd/ns/dpws/2009/01",
		XMLNsMsg:    "http://standards.ieee.org/downloads/11073/11073-10207-2017/message",
		XMLNsPm:     "http://standards.ieee.org/downloads/11073/11073-10207-2017/participant",
		XMLNsXsi:    "http://www.w3.org/2001/XMLSchema-instance",
		XMLNsMex:    "http://schemas.xmlsoap.org/ws/2004/09/mex",
		XMLNsSdc:    "http://standards.ieee.org/downloads/11073/11073-20701-2018",
		SOAPHeader:  header,
		SOAPBody:    body,
	}
}

//
//	Hello
//

func NewHelloSOAPHeader() SendSOAPHeader {
	return SendSOAPHeader{
		WsaAction:    "http://docs.oasis-open.org/ws-dd/ns/discovery/2009/01/Hello",
		WsaTo:        "urn:docs-oasis-open-org:ws-dd:ns:discovery:2009:01",
		WsaMessageId: uuid.NewString(),
		WsdAppSequence: &SendWsdAppSequence{
			InstanceId:    "4005719049",
			MessageNumber: "1",
		},
	}
}

func NewHelloSOAPBody(ipAddr *net.TCPAddr) SendSOAPBody {
	var helloPayload struct {
		//! Forced order by standard
		// Misordering leads to a scheme validation error of
		// content model: '(EndpointReference,Types?,Scopes?,XAddrs?,MetadataVersion,)'
		XMLName              xml.Name `xml:"wsd:Hello"`
		WsaEndpointReference SendWsaEndpointReference
		WsdTypes             string `xml:"wsd:Types"`
		WsdScopes            string `xml:"wsd:Scopes"`
		WsdXAddrs            string `xml:"wsd:XAddrs"`
		WsdMetadataVersion   string `xml:"wsd:MetadataVersion"`
	}

	// this data is taken from the sdcX implementation
	// TODO: allow customization
	helloPayload.WsdTypes = "dpws:Device mdpws:MedicalDevice"
	helloPayload.WsdScopes = "sdc.mds.pkp:1.2.840.10004.20701.1.1 sdc.ctxt.loc:/sdc.ctxt.loc.detail/DWHL///F05//TKl?fac=DWHL&amp;poc=F05&amp;bed=TKl sdc.ctxt.pat:/http://www.somda.org/ids/SamplePatientId123 sdc.ctxt.wfl:/http://www.somda.org/ids/WORKFLOW sdc.ctxt.ens:/http://www.somda.org/ids/ENSEMBLE sdc.ctxt.opr:/http://www.somda.org/ids/OPERATOR sdc.ctxt.mns:/http://www.somda.org/ids/MEANS sdc.cdc.type:///130535 sdc.cdc.type:///130536 sdc.cdc.type:///130736 sdc.cdc.type:/urn:oid:1.3.6.1.4.1.3592.2.1.1.0//DN_VMD"

	helloPayload.WsdMetadataVersion = "1"

	epRef := fmt.Sprintf("urn:uuid:%s", uuid.NewString())
	helloPayload.WsaEndpointReference.WsaAddress = epRef
	helloPayload.WsdXAddrs = fmt.Sprintf("http://%s:%d/%s", &ipAddr.IP, ipAddr.Port, epRef)

	return SendSOAPBody{Payload: helloPayload}
}
