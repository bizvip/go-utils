/******************************************************************************
 * Copyright (c) Archer++ 2024.                                               *
 ******************************************************************************/

package obs

import (
	"encoding/xml"
	"fmt"
)

// ObsError defines error response from OBS
type ObsError struct {
	BaseModel
	Status   string
	XMLName  xml.Name `xml:"Error"`
	Code     string   `xml:"Code" json:"code"`
	Message  string   `xml:"Message" json:"message"`
	Resource string   `xml:"Resource"`
	HostId   string   `xml:"HostId"`
}

// Format print obs error's log
func (err ObsError) Error() string {
	return fmt.Sprintf("obs: service returned error: Status=%s, Code=%s, Message=%s, RequestId=%s",
		err.Status, err.Code, err.Message, err.RequestId)
}
