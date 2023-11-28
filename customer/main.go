package customer

import (
	"go.slabs.io/nofy-go/client"
)

type Client struct {
	B *client.Backend
}

func (n *Client) GetAll() (*[]Customer, *client.Error) {
	var res client.Res[[]Customer]

	err := n.B.Call("GET", "/customer", nil, &res)
	if err != nil {
		return nil, err
	}

	if res.Error != "" {
		return nil, &client.Error{
			Code:    res.Error,
			Message: res.Message,
		}
	}

	return res.Data, nil
}
