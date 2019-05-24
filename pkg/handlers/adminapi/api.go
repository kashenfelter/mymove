package adminapi

import (
	"log"
	"net/http"

	"github.com/go-openapi/loads"

	adminops "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations"
	"github.com/transcom/mymove/pkg/gen/restapi"
	"github.com/transcom/mymove/pkg/handlers"
	officeuserservice "github.com/transcom/mymove/pkg/services/office_user"
)

// NewAdminAPIHandler returns a handler for the admin API
func NewAdminAPIHandler(context handlers.HandlerContext) http.Handler {

	// Wire up the handlers to the publicAPIMux
	adminSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	adminAPI := adminops.NewMymoveAPI(adminSpec)

	adminAPI.OfficeIndexOfficeUsersHandler = IndexOfficeUsersHandler{context, officeuserservice.NewOfficeUsersFetcher(context.DB())}

	return adminAPI.Serve(nil)
}
