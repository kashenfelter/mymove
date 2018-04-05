package handlers

import (
	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForAddressModel(a *models.Address) *internalmessages.Address {
	if a != nil {
		return &internalmessages.Address{
			StreetAddress1: swag.String(a.StreetAddress1),
			StreetAddress2: a.StreetAddress2,
			StreetAddress3: a.StreetAddress3,
			City:           swag.String(a.City),
			State:          swag.String(a.State),
			PostalCode:     swag.String(a.PostalCode),
			Country:        a.Country,
		}
	}
	return nil
}

func addressModelFromPayload(rawAddress *internalmessages.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	address := models.Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		StreetAddress3: rawAddress.StreetAddress3,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		PostalCode:     *rawAddress.PostalCode,
		Country:        rawAddress.Country,
	}
	return &address
}
