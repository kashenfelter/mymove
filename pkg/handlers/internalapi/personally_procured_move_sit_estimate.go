package internalapi

import (
	"reflect"
	"time"

	beeline "github.com/honeycombio/beeline-go"

	"github.com/transcom/mymove/pkg/models"

	"github.com/go-openapi/runtime/middleware"

	ppmop "github.com/transcom/mymove/pkg/gen/internalapi/internaloperations/ppm"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/unit"
)

// ShowPPMSitEstimateHandler returns PPM SIT estimate for a weight, move date,
type ShowPPMSitEstimateHandler struct {
	handlers.HandlerContext
}

// Handle calculates SIT charge and retrieves SIT discount rate.
// It returns the discount rate applied to relevant SIT charge.
func (h ShowPPMSitEstimateHandler) Handle(params ppmop.ShowPPMSitEstimateParams) middleware.Responder {
	ctx, span := beeline.StartSpan(params.HTTPRequest.Context(), reflect.TypeOf(h).Name())
	defer span.Send()

	engine := rateengine.NewRateEngine(h.DB(), h.Logger())
	sitZip3 := rateengine.Zip5ToZip3(params.DestinationZip)
	cwtWeight := unit.Pound(params.WeightEstimate).ToCWT()
	originalMoveDateTime := time.Time(params.OriginalMoveDate)

	lhDiscount, sitDiscount, err := models.PPMDiscountFetch(h.DB(),
		h.Logger(),
		params.OriginZip,
		params.DestinationZip,
		time.Time(params.OriginalMoveDate),
	)
	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error fetching personally procured move discount")
	}
	sitComputation, err := engine.SitCharge(cwtWeight, int(params.DaysInStorage), sitZip3, originalMoveDateTime, true)

	if err != nil {
		return h.RespondAndTraceError(ctx, err, "error computing sit charge")
	}

	// Swagger returns int64 when using the integer type
	sitCharge := int64(sitComputation.ApplyDiscount(lhDiscount, sitDiscount))

	ppmSitEstimate := internalmessages.PPMSitEstimate{
		Estimate: &sitCharge,
	}
	return ppmop.NewShowPPMSitEstimateOK().WithPayload(&ppmSitEstimate)
}
