// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import auth "github.com/transcom/mymove/pkg/auth"
import mock "github.com/stretchr/testify/mock"

import uuid "github.com/gofrs/uuid"

// StorageInTransitDeleter is an autogenerated mock type for the StorageInTransitDeleter type
type StorageInTransitDeleter struct {
	mock.Mock
}

// DeleteStorageInTransit provides a mock function with given fields: shipmentID, storageInTransitID, session
func (_m *StorageInTransitDeleter) DeleteStorageInTransit(shipmentID uuid.UUID, storageInTransitID uuid.UUID, session *auth.Session) error {
	ret := _m.Called(shipmentID, storageInTransitID, session)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID, uuid.UUID, *auth.Session) error); ok {
		r0 = rf(shipmentID, storageInTransitID, session)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
