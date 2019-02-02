package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentLineItemDimensions is an object representing dimensions of a shipment line item
type ShipmentLineItemDimensions struct {
	ID        uuid.UUID             `json:"id" db:"id"`
	Length    unit.BaseQuantityInch `json:"length" db:"length"`
	Width     unit.BaseQuantityInch `json:"width" db:"width"`
	Height    unit.BaseQuantityInch `json:"height" db:"height"`
	CreatedAt time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
}
