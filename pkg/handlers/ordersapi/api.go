package ordersapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/ordersapi"
	"github.com/transcom/mymove/pkg/gen/ordersapi/ordersoperations"
	"github.com/transcom/mymove/pkg/gen/ordersmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

// NewOrdersAPIHandler returns a handler for the Orders API
func NewOrdersAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the ordersAPIMux
	ordersSpec, err := loads.Analyzed(ordersapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	ordersAPI := ordersoperations.NewMymoveAPI(ordersSpec)
	ordersAPI.GetOrdersHandler = GetOrdersHandler{context}
	ordersAPI.GetOrdersByIssuerAndOrdersNumHandler = GetOrdersByIssuerAndOrdersNumHandler{context}
	ordersAPI.IndexOrdersForMemberHandler = IndexOrdersForMemberHandler{context}
	ordersAPI.PostRevisionHandler = PostRevisionHandler{context}
	ordersAPI.PostRevisionToOrdersHandler = PostRevisionToOrdersHandler{context}
	return ordersAPI.Serve(nil)
}

func payloadForElectronicOrderModel(order models.ElectronicOrder) (*ordersmessages.Orders, error) {
	var revisionPayloads []*ordersmessages.Revision
	for _, revision := range order.Revisions {
		payload, err := payloadForElectronicOrdersRevisionModel(revision)
		if err != nil {
			return nil, err
		}
		revisionPayloads = append(revisionPayloads, payload)
	}

	ordersPayload := &ordersmessages.Orders{
		UUID:      strfmt.UUID(order.ID.String()),
		OrdersNum: order.OrdersNumber,
		Edipi:     order.Edipi,
		Issuer:    order.Issuer,
		Revisions: revisionPayloads,
	}
	return ordersPayload, nil
}

func payloadForElectronicOrdersRevisionModel(revision models.ElectronicOrdersRevision) (*ordersmessages.Revision, error) {
	seqNum := int64(revision.SeqNum)
	revisionPayload := &ordersmessages.Revision{
		SeqNum: &seqNum,
		Member: &ordersmessages.Member{
			GivenName:   revision.GivenName,
			MiddleName:  revision.MiddleName,
			FamilyName:  revision.FamilyName,
			Suffix:      revision.NameSuffix,
			Affiliation: revision.Affiliation,
			Rank:        revision.Paygrade,
			Title:       revision.Title,
		},
		Status:        revision.Status,
		DateIssued:    (*strfmt.DateTime)(&revision.DateIssued),
		NoCostMove:    revision.NoCostMove,
		TdyEnRoute:    revision.TdyEnRoute,
		TourType:      revision.TourType,
		OrdersType:    revision.OrdersType,
		HasDependents: &revision.HasDependents,
		LosingUnit: &ordersmessages.Unit{
			Uic:        revision.LosingUIC,
			Name:       revision.LosingUnitName,
			City:       revision.LosingUnitCity,
			Locality:   revision.LosingUnitLocality,
			Country:    revision.LosingUnitCountry,
			PostalCode: revision.LosingUnitPostalCode,
		},
		GainingUnit: &ordersmessages.Unit{
			Uic:        revision.GainingUIC,
			Name:       revision.GainingUnitName,
			City:       revision.GainingUnitCity,
			Locality:   revision.GainingUnitLocality,
			Country:    revision.GainingUnitCountry,
			PostalCode: revision.GainingUnitPostalCode,
		},
		ReportNoEarlierThan: (*strfmt.Date)(revision.ReportNoEarlierThan),
		ReportNoLaterThan:   (*strfmt.Date)(revision.ReportNoLaterThan),
		PcsAccounting: &ordersmessages.Accounting{
			Tac: revision.HhgTAC,
			Sdn: revision.HhgSDN,
			Loa: revision.HhgLOA,
		},
		NtsAccounting: &ordersmessages.Accounting{
			Tac: revision.NtsTAC,
			Sdn: revision.NtsSDN,
			Loa: revision.NtsLOA,
		},
		PovShipmentAccounting: &ordersmessages.Accounting{
			Tac: revision.PovShipmentTAC,
			Sdn: revision.PovShipmentSDN,
			Loa: revision.PovShipmentLOA,
		},
		PovStorageAccounting: &ordersmessages.Accounting{
			Tac: revision.PovStorageTAC,
			Sdn: revision.PovStorageSDN,
			Loa: revision.PovStorageLOA,
		},
		UbAccounting: &ordersmessages.Accounting{
			Tac: revision.UbTAC,
			Sdn: revision.UbSDN,
			Loa: revision.UbLOA,
		},
		Comments: revision.Comments,
	}
	return revisionPayload, nil
}
