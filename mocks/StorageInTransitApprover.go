// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import apimessages "github.com/transcom/mymove/pkg/gen/apimessages"
import auth "github.com/transcom/mymove/pkg/auth"
import mock "github.com/stretchr/testify/mock"
import models "github.com/transcom/mymove/pkg/models"

import uuid "github.com/gofrs/uuid"
import validate "github.com/gobuffalo/validate"

// StorageInTransitApprover is an autogenerated mock type for the StorageInTransitApprover type
type StorageInTransitApprover struct {
	mock.Mock
}

// ApproveStorageInTransit provides a mock function with given fields: payload, shipmentID, session, storageInTransitID
func (_m *StorageInTransitApprover) ApproveStorageInTransit(payload apimessages.StorageInTransitApprovalPayload, shipmentID uuid.UUID, session *auth.Session, storageInTransitID uuid.UUID) (*models.StorageInTransit, *validate.Errors, error) {
	ret := _m.Called(payload, shipmentID, session, storageInTransitID)

	var r0 *models.StorageInTransit
	if rf, ok := ret.Get(0).(func(apimessages.StorageInTransitApprovalPayload, uuid.UUID, *auth.Session, uuid.UUID) *models.StorageInTransit); ok {
		r0 = rf(payload, shipmentID, session, storageInTransitID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.StorageInTransit)
		}
	}

	var r1 *validate.Errors
	if rf, ok := ret.Get(1).(func(apimessages.StorageInTransitApprovalPayload, uuid.UUID, *auth.Session, uuid.UUID) *validate.Errors); ok {
		r1 = rf(payload, shipmentID, session, storageInTransitID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*validate.Errors)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(apimessages.StorageInTransitApprovalPayload, uuid.UUID, *auth.Session, uuid.UUID) error); ok {
		r2 = rf(payload, shipmentID, session, storageInTransitID)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
