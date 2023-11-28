package nofy

import (
	"encoding/base64"

	"go.slabs.io/nofy-go/client"
	"go.slabs.io/nofy-go/customer"
	"go.slabs.io/nofy-go/upload"
)

type API struct {
	Backend  *client.Backend
	Customer *customer.Client
	Upload   *upload.Client
}

func New(k string) *API {
	// base64 encode the key
	key := base64.StdEncoding.EncodeToString([]byte(k + ":"))

	b := &client.Backend{
		Key: key,
	}

	return &API{
		Backend:  b,
		Customer: &customer.Client{B: b},
		Upload:   &upload.Client{B: b},
	}
}
