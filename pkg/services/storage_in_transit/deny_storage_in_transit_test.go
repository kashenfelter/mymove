package storageintransit

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *StorageInTransitServiceSuite) TestDenyStorageInTransit() {
	shipment, sit, user := setupStorageInTransitServiceTest(suite)
	tspUser := testdatagen.MakeDefaultTspUser(suite.DB())
	session := auth.Session{
		ApplicationName: auth.TspApp,
		UserID:          *tspUser.UserID,
		IDToken:         "fake token",
		TspUserID:       tspUser.ID,
	}
	payload := apimessages.StorageInTransitDenialPayload{
		AuthorizationNotes: *handlers.FmtString("looks bad to me"),
	}

	denier := NewStorageInTransitDenier(suite.DB())

	// Should not work for a TSP user
	_, _, err := denier.DenyStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "FETCH_FORBIDDEN")

	// Should not work if the status is already delivered
	sit.Status = models.StorageInTransitStatusDELIVERED
	_, _ = suite.DB().ValidateAndSave(&sit)

	_, _, err = denier.DenyStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.Error(err, "WRITE_CONFLICT")

	// Happy path
	sit.Status = models.StorageInTransitStatusREQUESTED
	_, _ = suite.DB().ValidateAndSave(&sit)

	session = auth.Session{
		ApplicationName: auth.OfficeApp,
		UserID:          *user.UserID,
		IDToken:         "fake token",
		OfficeUserID:    user.ID,
	}

	actualStorageInTransit, verrs, err := denier.DenyStorageInTransit(payload, shipment.ID, &session, sit.ID)
	suite.NoError(err)
	suite.False(verrs.HasAny())
	suite.Equal(models.StorageInTransitStatusDENIED, actualStorageInTransit.Status)
	suite.Equal(payload.AuthorizationNotes, *actualStorageInTransit.AuthorizationNotes)
}
