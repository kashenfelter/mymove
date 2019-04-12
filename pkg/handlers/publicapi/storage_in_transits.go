package publicapi

import (
	"fmt"
	"time"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	sitop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/storage_in_transits"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func payloadForStorageInTransitModel(s *models.StorageInTransit) *apimessages.StorageInTransit {
	if s == nil {
		return nil
	}

	location := string(s.Location)
	status := string(s.Status)

	return &apimessages.StorageInTransit{
		ID:                  *handlers.FmtUUID(s.ID),
		ShipmentID:          *handlers.FmtUUID(s.ShipmentID),
		EstimatedStartDate:  handlers.FmtDate(s.EstimatedStartDate),
		Notes:               handlers.FmtStringPtr(s.Notes),
		WarehouseAddress:    payloadForAddressModel(&s.WarehouseAddress),
		WarehouseEmail:      handlers.FmtStringPtr(s.WarehouseEmail),
		WarehouseID:         handlers.FmtString(s.WarehouseID),
		WarehouseName:       handlers.FmtString(s.WarehouseName),
		WarehousePhone:      handlers.FmtStringPtr(s.WarehousePhone),
		Location:            &location,
		Status:              *handlers.FmtString(status),
		AuthorizationNotes:  handlers.FmtStringPtr(s.AuthorizationNotes),
		AuthorizedStartDate: handlers.FmtDatePtr(s.AuthorizedStartDate),
		ActualStartDate:     handlers.FmtDatePtr(s.ActualStartDate),
		OutDate:             handlers.FmtDatePtr(s.OutDate),
	}
}

func authorizeStorageInTransitRequest(db *pop.Connection, session *auth.Session, shipmentID uuid.UUID, allowOffice bool) (isUserAuthorized bool, err error) {
	if session.IsTspUser() {
		_, _, err := models.FetchShipmentForVerifiedTSPUser(db, session.TspUserID, shipmentID)

		if err != nil {
			return false, err
		}
		return true, nil
	} else if session.IsOfficeUser() {
		if allowOffice {
			return true, nil
		}
	} else {
		return false, models.ErrFetchForbidden
	}
	return false, models.ErrFetchForbidden
}

func processStorageInTransitInput(h handlers.HandlerContext, shipmentID uuid.UUID, payload apimessages.StorageInTransit) (models.StorageInTransit, error) {
	incomingLocation := *payload.Location
	var savedLocation models.StorageInTransitLocation

	if incomingLocation == "ORIGIN" {
		savedLocation = models.StorageInTransitLocationORIGIN
	} else {
		savedLocation = models.StorageInTransitLocationDESTINATION
	}

	var estimatedStartDate time.Time
	if payload.EstimatedStartDate != nil {
		estimatedStartDate = time.Time(*payload.EstimatedStartDate)
	}

	var warehouseName string
	if payload.WarehouseName != nil {
		warehouseName = *payload.WarehouseName
	}

	var warehouseAddress models.Address
	if payload.WarehouseAddress != nil {
		warehouseAddress = models.Address{
			StreetAddress1: *payload.WarehouseAddress.StreetAddress1,
			StreetAddress2: payload.WarehouseAddress.StreetAddress2,
			StreetAddress3: payload.WarehouseAddress.StreetAddress3,
			City:           *payload.WarehouseAddress.City,
			State:          *payload.WarehouseAddress.State,
			PostalCode:     *payload.WarehouseAddress.PostalCode,
			Country:        payload.WarehouseAddress.Country,
		}
	}

	newStorageInTransit := models.StorageInTransit{
		ShipmentID:         shipmentID,
		Location:           savedLocation,
		EstimatedStartDate: estimatedStartDate,
		Notes:              payload.Notes,
		WarehouseID:        *payload.WarehouseID,
		WarehouseName:      warehouseName,
		WarehouseAddressID: warehouseAddress.ID,
		WarehouseAddress:   warehouseAddress,
		WarehouseEmail:     payload.WarehouseEmail,
		WarehousePhone:     payload.WarehousePhone,
	}

	return newStorageInTransit, nil

}

func patchStorageInTransitWithPayload(storageInTransit *models.StorageInTransit, payload *apimessages.StorageInTransit) {
	if *payload.Location == "ORIGIN" {
		storageInTransit.Location = models.StorageInTransitLocationORIGIN
	} else {
		storageInTransit.Location = models.StorageInTransitLocationDESTINATION
	}

	if payload.EstimatedStartDate != nil {
		storageInTransit.EstimatedStartDate = *(*time.Time)(payload.EstimatedStartDate)
	}

	storageInTransit.Notes = handlers.FmtStringPtrNonEmpty(payload.Notes)

	if payload.WarehouseID != nil {
		storageInTransit.WarehouseID = *payload.WarehouseID
	}

	if payload.WarehouseName != nil {
		storageInTransit.WarehouseName = *payload.WarehouseName
	}

	if payload.WarehouseAddress != nil {
		updateAddressWithPayload(&storageInTransit.WarehouseAddress, payload.WarehouseAddress)
	}

	storageInTransit.WarehousePhone = handlers.FmtStringPtrNonEmpty(payload.WarehousePhone)
	storageInTransit.WarehouseEmail = handlers.FmtStringPtrNonEmpty(payload.WarehouseEmail)
}

// IndexStorageInTransitHandler returns a list of Storage In Transit entries
type IndexStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to return a list of storage in transits using their shipment ID.
func (h IndexStorageInTransitHandler) Handle(params sitop.IndexStorageInTransitsParams) middleware.Responder {
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("Unauthorized User", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransits, err := models.FetchStorageInTransitsOnShipment(h.DB(), shipmentID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransitsList := make(apimessages.StorageInTransits, len(storageInTransits))

	for i, storageInTransit := range storageInTransits {
		storageInTransitsList[i] = payloadForStorageInTransitModel(&storageInTransit)
	}

	return sitop.NewIndexStorageInTransitsOK().WithPayload(storageInTransitsList)
}

// CreateStorageInTransitHandler creates a storage in transit entry and returns a payload.
type CreateStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to create a storage in transit and return it in a payload
func (h CreateStorageInTransitHandler) Handle(params sitop.CreateStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransit
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, false)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	newStorageInTransit, err := processStorageInTransitInput(h, shipmentID, *payload)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	newStorageInTransit.Status = models.StorageInTransitStatusREQUESTED

	verrs, err := models.SaveStorageInTransitAndAddress(h.DB(), &newStorageInTransit)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(&newStorageInTransit)

	return sitop.NewCreateStorageInTransitCreated().WithPayload(storageInTransitPayload)

}

// ApproveStorageInTransitHandler approves an existing Storage In Transit
type ApproveStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to approved, save the authorization notes that
// support that status, save the authorization date, and return the saved object in a payload.
func (h ApproveStorageInTransitHandler) Handle(params sitop.ApproveStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransitApprovalPayload
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	// Only office users are authorized to do this.

	if !session.IsOfficeUser() {
		return sitop.NewApproveStorageInTransitForbidden()
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error(fmt.Sprintf("Could not find existing storage in transit with id: %s ", storageInTransitID), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		h.Logger().Error(fmt.Sprintf("ShipmentID: ( %s ) provided does not match the storage in transit", shipmentID))
		return sitop.NewApproveStorageInTransitForbidden()
	}

	if storageInTransit.Status == models.StorageInTransitStatusDELIVERED {
		h.Logger().Error(fmt.Sprintf("Cannot approve storage in transit that's already delivered. ID: %s", storageInTransitID))
		return sitop.NewApproveStorageInTransitConflict()
	}

	storageInTransit.Status = models.StorageInTransitStatusAPPROVED
	storageInTransit.AuthorizationNotes = &payload.AuthorizationNotes
	storageInTransit.AuthorizedStartDate = (*time.Time)(&payload.AuthorizedStartDate)

	if verrs, err := h.DB().ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewApproveStorageInTransitOK().WithPayload(returnPayload)

}

// DenyStorageInTransitHandler denies an existing storage in transit
type DenyStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to denied, save the supporting authorization notes,
// and return the saved object in a payload.
func (h DenyStorageInTransitHandler) Handle(params sitop.DenyStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransitApprovalPayload
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Only office users are authorized to do this
	if !session.IsOfficeUser() {

		return sitop.NewDenyStorageInTransitForbidden()
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error(fmt.Sprintf("Could not find existing storage in transit with id: %s ", storageInTransitID), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		h.Logger().Error(fmt.Sprintf("ShipmentID: ( %s ) provided does not match the storage in transit", shipmentID))
		return sitop.NewDenyStorageInTransitForbidden()
	}

	if storageInTransit.Status == models.StorageInTransitStatusDELIVERED {
		h.Logger().Error(fmt.Sprintf("Cannot approve storage in transit that's already delivered. ID: %s", storageInTransitID))
		return sitop.NewDenyStorageInTransitConflict()
	}

	storageInTransit.Status = models.StorageInTransitStatusDENIED
	storageInTransit.AuthorizationNotes = &payload.AuthorizationNotes

	if verrs, err := h.DB().ValidateAndSave(storageInTransit); verrs.HasAny() || err != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewDenyStorageInTransitOK().WithPayload(returnPayload)

}

// InSitStorageInTransitHandler places storage in transit into in sit
type InSitStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to 'in SIT' and return the saved object in a payload.
func (h InSitStorageInTransitHandler) Handle(params sitop.InSitStorageInTransitParams) middleware.Responder {
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	inSitPayload := params.StorageInTransitInSitPayload

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Only TSPs are authorized to do this
	if !session.IsTspUser() {
		return sitop.NewInSitStorageInTransitForbidden()
	}

	// Make sure the TSP is authorized for the shipment
	_, _, err = models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)

	if err != nil {
		return sitop.NewInSitStorageInTransitForbidden()
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error(fmt.Sprintf("Could not find existing storage in transit with id: %s ", storageInTransitID), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Verify that the shipment we're getting matches what's in the storage in transit
	if shipmentID != storageInTransit.ShipmentID {
		h.Logger().Error(fmt.Sprintf("ShipmentID: ( %s ) provided does not match the storage in transit", shipmentID))
		return sitop.NewApproveStorageInTransitForbidden()
	}

	payloadActualStartDate := (time.Time)(inSitPayload.ActualStartDate)

	if !(storageInTransit.Status == models.StorageInTransitStatusAPPROVED) {
		h.Logger().Error(fmt.Sprintf("Cannot place into SIT without approval. ID: %s ", storageInTransitID))
		return sitop.NewInSitStorageInTransitConflict()
	}

	storageInTransit.Status = models.StorageInTransitStatusINSIT
	storageInTransit.ActualStartDate = &payloadActualStartDate

	if verrs, errs := h.DB().ValidateAndSave(storageInTransit); verrs.HasAny() || errs != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewInSitStorageInTransitOK().WithPayload(returnPayload)

}

// DeliverStorageInTransitHandler delivers an existing storage in transit
type DeliverStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to set the status for a storage in transit to delivered and return the saved object in a payload.
func (h DeliverStorageInTransitHandler) Handle(params sitop.DeliverStorageInTransitParams) middleware.Responder {
	// TODO: it looks like from the wireframes for the delivery status change form that this will also need to edit
	//  delivery address(es) and the actual delivery date.
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	fmt.Println(session.IsTspUser())

	// Only TSPs are authorized to do this
	if !session.IsTspUser() {
		fmt.Println("was false")
		return sitop.NewDeliverStorageInTransitForbidden()
	}

	// Make sure the TSP is authorized for the shipment
	_, _, err = models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)

	if err != nil {
		return sitop.NewDeliverStorageInTransitForbidden()
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error(fmt.Sprintf("Could not find existing storage in transit with id: %s ", storageInTransitID), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT) &&
		!(storageInTransit.Status == models.StorageInTransitStatusRELEASED) {
		h.Logger().Error(fmt.Sprintf("Cannot deliver if its not in sit or released. ID: %s", storageInTransitID))
		return sitop.NewDeliverStorageInTransitConflict()
	}

	storageInTransit.Status = models.StorageInTransitStatusDELIVERED

	if verrs, errs := h.DB().ValidateAndSave(storageInTransit); verrs.HasAny() || errs != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewDeliverStorageInTransitOK().WithPayload(returnPayload)

}

// ReleaseStorageInTransitHandler releases an existing storage in transit
type ReleaseStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle Handles the handling
// This is meant to set the status of a storage in transit to released, save the actual date that supports that,
// and return the saved object in a payload.
func (h ReleaseStorageInTransitHandler) Handle(params sitop.ReleaseStorageInTransitParams) middleware.Responder {
	// TODO: There may be other fields that have to be addressed here when we get to the frontend story for this.
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	payload := params.StorageInTransitOnReleasePayload

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)

	// Only TSPs are authorized to do this
	if !session.IsTspUser() {
		return sitop.NewReleaseStorageInTransitForbidden()
	}

	// Make sure the TSP is authorized for the shipment
	_, _, err = models.FetchShipmentForVerifiedTSPUser(h.DB(), session.TspUserID, shipmentID)

	if err != nil {
		return sitop.NewReleaseStorageInTransitForbidden()
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error(fmt.Sprintf("Could not find existing storage in transit with id: %s ", storageInTransitID), zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	// Make sure we're not releasing something that wasn't in SIT or in delivered status.
	// The latter is there so that we can 'undo' a mistaken deliver action.
	if !(storageInTransit.Status == models.StorageInTransitStatusINSIT) &&
		!(storageInTransit.Status == models.StorageInTransitStatusDELIVERED) {
		h.Logger().Error(fmt.Sprintf("Cannot release something from storage in transit that wasn't in it. ID: %s", storageInTransitID))
		return sitop.NewReleaseStorageInTransitConflict()
	}
	storageInTransit.Status = models.StorageInTransitStatusRELEASED
	storageInTransit.OutDate = (*time.Time)(&payload.ReleasedOn)

	if verrs, errs := h.DB().ValidateAndSave(storageInTransit); verrs.HasAny() || errs != nil {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	returnPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewReleaseStorageInTransitOK().WithPayload(returnPayload)

}

// PatchStorageInTransitHandler updates an existing Storage In Transit entry
type PatchStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to edit a storage in transit object based on provided parameters and return the saved object
// in a payload
func (h PatchStorageInTransitHandler) Handle(params sitop.PatchStorageInTransitParams) middleware.Responder {
	payload := params.StorageInTransit
	shipmentID, err := uuid.FromString(params.ShipmentID.String())
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)
	if err != nil {
		h.Logger().Error("Could not find existing SIT record", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	if storageInTransit.ShipmentID != shipmentID {
		h.Logger().Error("Shipment ID clash between endpoint URL and SIT record")
		return sitop.NewPatchStorageInTransitForbidden()
	}

	patchStorageInTransitWithPayload(storageInTransit, payload)

	verrs, err := models.SaveStorageInTransitAndAddress(h.DB(), storageInTransit)
	if err != nil || verrs.HasAny() {
		return handlers.ResponseForVErrors(h.Logger(), verrs, err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)

	return sitop.NewPatchStorageInTransitOK().WithPayload(storageInTransitPayload)
}

// GetStorageInTransitHandler gets a single Storage In Transit based on its own ID
type GetStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to fetch a single storage in transit using its shipment and object IDs
func (h GetStorageInTransitHandler) Handle(params sitop.GetStorageInTransitParams) middleware.Responder {
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	shipmentID, err := uuid.FromString(params.ShipmentID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, true)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransit, err := models.FetchStorageInTransitByID(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	storageInTransitPayload := payloadForStorageInTransitModel(storageInTransit)
	return sitop.NewGetStorageInTransitOK().WithPayload(storageInTransitPayload)

}

// DeleteStorageInTransitHandler deletes a Storage in Transit based on the provided ID
type DeleteStorageInTransitHandler struct {
	handlers.HandlerContext
}

// Handle handles the handling
// This is meant to delete a storage in transit object using its own shipment and object IDs
func (h DeleteStorageInTransitHandler) Handle(params sitop.DeleteStorageInTransitParams) middleware.Responder {
	storageInTransitID, err := uuid.FromString(params.StorageInTransitID.String())
	shipmentID, err := uuid.FromString(params.ShipmentID.String())

	if err != nil {
		h.Logger().Error("UUID Parsing", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	session := auth.SessionFromRequestContext(params.HTTPRequest)
	isUserAuthorized, err := authorizeStorageInTransitRequest(h.DB(), session, shipmentID, false)

	if isUserAuthorized == false {
		h.Logger().Error("User is unauthorized", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	err = models.DeleteStorageInTransit(h.DB(), storageInTransitID)

	if err != nil {
		h.Logger().Error("DB Query", zap.Error(err))
		return handlers.ResponseForError(h.Logger(), err)
	}

	return sitop.NewDeleteStorageInTransitOK()

}
