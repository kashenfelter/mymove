package models_test

import (
	"context"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestElectronicOrdersRevisionValidateAndCreate() {
	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       ordersmessages.IssuerAirForce,
		OrdersNumber: "8675309",
	}
	verrs, err := models.CreateElectronicOrder(context.Background(), suite.DB(), &order)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	rev := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       ordersmessages.AffiliationAirForce,
		Paygrade:          ordersmessages.RankE1,
		Status:            ordersmessages.StatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          ordersmessages.TourTypeAccompanied,
		OrdersType:        ordersmessages.OrdersTypeSeparation,
		HasDependents:     true,
	}

	verrs, err = suite.DB().ValidateAndCreate(&rev)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestElectronicOrdersRevisionValidations() {
	empty := ""
	revision := &models.ElectronicOrdersRevision{
		SeqNum:                -1,
		MiddleName:            &empty,
		NameSuffix:            &empty,
		Title:                 &empty,
		LosingUIC:             &empty,
		LosingUnitName:        &empty,
		LosingUnitCity:        &empty,
		LosingUnitLocality:    &empty,
		LosingUnitCountry:     &empty,
		LosingUnitPostalCode:  &empty,
		GainingUIC:            &empty,
		GainingUnitName:       &empty,
		GainingUnitCity:       &empty,
		GainingUnitLocality:   &empty,
		GainingUnitCountry:    &empty,
		GainingUnitPostalCode: &empty,
		HhgTAC:                &empty,
		HhgSDN:                &empty,
		HhgLOA:                &empty,
		NtsTAC:                &empty,
		NtsSDN:                &empty,
		NtsLOA:                &empty,
		PovShipmentTAC:        &empty,
		PovShipmentSDN:        &empty,
		PovShipmentLOA:        &empty,
		PovStorageTAC:         &empty,
		PovStorageSDN:         &empty,
		PovStorageLOA:         &empty,
		UbTAC:                 &empty,
		UbSDN:                 &empty,
		UbLOA:                 &empty,
	}

	var expErrors = map[string][]string{
		"electronic_order_id":      {"ElectronicOrderID can not be blank."},
		"seq_num":                  {"-1 is not greater than -1."},
		"given_name":               {"GivenName can not be blank."},
		"middle_name":              {"MiddleName can not be blank."},
		"family_name":              {"FamilyName can not be blank."},
		"name_suffix":              {"NameSuffix can not be blank."},
		"paygrade":                 {"Paygrade can not be blank."},
		"affiliation":              {"Affiliation can not be blank."},
		"title":                    {"Title can not be blank."},
		"date_issued":              {"DateIssued can not be blank."},
		"status":                   {"Status can not be blank."},
		"tour_type":                {"TourType can not be blank."},
		"orders_type":              {"OrdersType can not be blank."},
		"losing_ui_c":              {"LosingUIC can not be blank."},
		"losing_unit_name":         {"LosingUnitName can not be blank."},
		"losing_unit_city":         {"LosingUnitCity can not be blank."},
		"losing_unit_locality":     {"LosingUnitLocality can not be blank."},
		"losing_unit_postal_code":  {"LosingUnitPostalCode can not be blank."},
		"losing_unit_country":      {"LosingUnitCountry can not be blank."},
		"gaining_ui_c":             {"GainingUIC can not be blank."},
		"gaining_unit_name":        {"GainingUnitName can not be blank."},
		"gaining_unit_city":        {"GainingUnitCity can not be blank."},
		"gaining_unit_locality":    {"GainingUnitLocality can not be blank."},
		"gaining_unit_postal_code": {"GainingUnitPostalCode can not be blank."},
		"gaining_unit_country":     {"GainingUnitCountry can not be blank."},
		"hhg_t_a_c":                {"HhgTAC can not be blank."},
		"hhg_s_d_n":                {"HhgSDN can not be blank."},
		"hhg_l_o_a":                {"HhgLOA can not be blank."},
		"nts_t_a_c":                {"NtsTAC can not be blank."},
		"nts_s_d_n":                {"NtsSDN can not be blank."},
		"nts_l_o_a":                {"NtsLOA can not be blank."},
		"pov_shipment_t_a_c":       {"PovShipmentTAC can not be blank."},
		"pov_shipment_s_d_n":       {"PovShipmentSDN can not be blank."},
		"pov_shipment_l_o_a":       {"PovShipmentLOA can not be blank."},
		"pov_storage_t_a_c":        {"PovStorageTAC can not be blank."},
		"pov_storage_s_d_n":        {"PovStorageSDN can not be blank."},
		"pov_storage_l_o_a":        {"PovStorageLOA can not be blank."},
		"ub_t_a_c":                 {"UbTAC can not be blank."},
		"ub_l_o_a":                 {"UbLOA can not be blank."},
		"ub_s_d_n":                 {"UbSDN can not be blank."},
	}

	suite.verifyValidationErrors(revision, expErrors)
}

func (suite *ModelSuite) TestCreateElectronicOrdersRevision() {
	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       ordersmessages.IssuerAirForce,
		OrdersNumber: "8675309",
	}
	verrs, err := models.CreateElectronicOrder(context.Background(), suite.DB(), &order)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	rev := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       ordersmessages.AffiliationAirForce,
		Paygrade:          ordersmessages.RankE1,
		Status:            ordersmessages.StatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          ordersmessages.TourTypeAccompanied,
		OrdersType:        ordersmessages.OrdersTypeSeparation,
		HasDependents:     true,
	}

	verrs, err = models.CreateElectronicOrdersRevision(context.Background(), suite.DB(), &rev)
	suite.NoError(err)
	suite.NoVerrs(verrs)
}

func (suite *ModelSuite) TestCreateElectronicOrdersRevision_Amendment() {
	order := models.ElectronicOrder{
		Edipi:        "1234567890",
		Issuer:       ordersmessages.IssuerAirForce,
		OrdersNumber: "8675309",
	}
	ctx := context.Background()
	verrs, err := models.CreateElectronicOrder(ctx, suite.DB(), &order)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	rev0 := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            0,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       ordersmessages.AffiliationAirForce,
		Paygrade:          ordersmessages.RankE1,
		Status:            ordersmessages.StatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          ordersmessages.TourTypeAccompanied,
		OrdersType:        ordersmessages.OrdersTypeSeparation,
		HasDependents:     true,
	}

	models.CreateElectronicOrdersRevision(ctx, suite.DB(), &rev0)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	rev1 := models.ElectronicOrdersRevision{
		ElectronicOrderID: order.ID,
		ElectronicOrder:   order,
		SeqNum:            1,
		GivenName:         "First",
		FamilyName:        "Last",
		Affiliation:       ordersmessages.AffiliationAirForce,
		Paygrade:          ordersmessages.RankE1,
		Status:            ordersmessages.StatusAuthorized,
		DateIssued:        time.Now(),
		NoCostMove:        false,
		TdyEnRoute:        false,
		TourType:          ordersmessages.TourTypeAccompanied,
		OrdersType:        ordersmessages.OrdersTypeSeparation,
		HasDependents:     true,
	}

	models.CreateElectronicOrdersRevision(ctx, suite.DB(), &rev1)
	suite.NoError(err)
	suite.NoVerrs(verrs)

	retrievedOrder, err := models.FetchElectronicOrderByID(suite.DB(), order.ID)
	suite.NoError(err)
	suite.Len(retrievedOrder.Revisions, 2)
	expectedIDs := []uuid.UUID{rev0.ID, rev1.ID}
	suite.Contains(expectedIDs, retrievedOrder.Revisions[0].ID)
	suite.Contains(expectedIDs, retrievedOrder.Revisions[1].ID)
	suite.NotEqual(retrievedOrder.Revisions[0].ID, retrievedOrder.Revisions[1].ID)
}
