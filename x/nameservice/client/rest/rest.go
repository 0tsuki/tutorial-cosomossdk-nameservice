package rest

import (
	"fmt"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

const (
	restName = "name"
)

// RegisterRoutes registers nameservice-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/names", storeName), namesHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}", storeName, restName), resolveNameHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/names/{%s}/whois", storeName, restName), whoisHandler(cliCtx, storeName)).Methods("GET")
}
