package models

import (
	"fmt"
	"strings"

	models "github.com/awareseven/mobilesniper/pkg/models/soap"
)

type ServiceType string
type DeviceType string

const (
	GetService              ServiceType = "GetService"
	SetService              ServiceType = "SetService"
	StateEventService       ServiceType = "StateEventService"
	ContextService          ServiceType = "ContextService"
	DescriptionEventService ServiceType = "DescriptionEventService"
	WaveformService         ServiceType = "WaveformService"
	LocalizationService     ServiceType = "LocalizationService"
)

const (
	MedicalDevice DeviceType = "MedicalDevice"
)

type SDCDevice struct {
	Name                     string
	FirmwareVersion          string
	SerialNumber             string
	Manufacturer             string
	ManufacturerUrl          string
	ModelName                string
	ModelUrl                 string
	PresentationUrl          string
	SupportedSdcTypes        []ServiceType
	DeviceTypes              []DeviceType
	EndpointReferenceAddress string
}

// Return human readable SDC device
func (d *SDCDevice) String() string {
	return fmt.Sprintf("%s (SN: %s | FW v%s) by %s", d.Name, d.SerialNumber, d.FirmwareVersion, d.Manufacturer)
}

// Create a SDCDevice from SOAP Get-request
func CreateSDCDevicebyGetResponse(soap models.ReceiveSOAPMessage) (SDCDevice, error) {

	retVal := SDCDevice{}

	getBody, ok := soap.SOAPBody.Payload.(models.GetResponseBody)
	if !ok {
		return retVal, fmt.Errorf("`soap` is not a 'Get' response message")
	}

	for _, sec := range getBody.GetMexMetadataSections() {
		if sec.DpwsThisDevice.DpwsFriendlyName != "" {
			// Metadata section contains device information
			retVal.Name = sec.DpwsThisDevice.DpwsFriendlyName
			retVal.FirmwareVersion = sec.DpwsThisDevice.DpwsFirmwareVersion
			retVal.SerialNumber = sec.DpwsThisDevice.DpwsSerialNumber
		} else if sec.DpwsThisModel.DpwsManufacturer != "" {
			// Metadata section contains model information
			retVal.Manufacturer = sec.DpwsThisModel.DpwsManufacturer
			retVal.ManufacturerUrl = sec.DpwsThisModel.DpwsManufacturerUrl
			retVal.ModelName = sec.DpwsThisModel.DpwsModelName
			retVal.ModelUrl = sec.DpwsThisModel.DpwsModelUrl
			retVal.PresentationUrl = sec.DpwsThisModel.DpwsPresentationUrl
		} else if sec.DpwsThisRelationship.DpwsHost.DpwsTypes != "" {
			// Metadata section contains relationship information
			retVal.EndpointReferenceAddress = sec.DpwsThisRelationship.DpwsHosted.WsaEndpointReference.WsaAddress
			retVal.DeviceTypes = parseDeviceTypes(sec.DpwsThisRelationship.DpwsHost.DpwsTypes)
			retVal.SupportedSdcTypes = parseServiceTypes(sec.DpwsThisRelationship.DpwsHosted.DpwsTypes)
		} else {
			// Metadata section contains location information
		}
	}

	return retVal, nil
}

// Parse string to detect SDC device type.
func parseDeviceTypes(types string) []DeviceType {

	retVal := []DeviceType{}

	if strings.Contains(types, string(MedicalDevice)) {
		retVal = append(retVal, MedicalDevice)
	}

	return retVal
}

// Parse string to detect SDC device type.
func parseServiceTypes(types string) []ServiceType {

	retVal := []ServiceType{}

	if strings.Contains(types, string(GetService)) {
		retVal = append(retVal, GetService)
	}

	if strings.Contains(types, string(SetService)) {
		retVal = append(retVal, SetService)
	}

	if strings.Contains(types, string(StateEventService)) {
		retVal = append(retVal, StateEventService)
	}

	if strings.Contains(types, string(ContextService)) {
		retVal = append(retVal, ContextService)
	}

	if strings.Contains(types, string(DescriptionEventService)) {
		retVal = append(retVal, DescriptionEventService)
	}

	if strings.Contains(types, string(WaveformService)) {
		retVal = append(retVal, WaveformService)
	}

	if strings.Contains(types, string(LocalizationService)) {
		retVal = append(retVal, LocalizationService)
	}

	return retVal
}
