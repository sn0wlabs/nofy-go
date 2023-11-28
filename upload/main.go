package upload

import (
	"os"

	"go.slabs.io/nofy-go/client"
)

type Client struct {
	B *client.Backend
}

func (n *Client) WebP(image *os.File) (*WebPUpload, *client.Error) {
	var res client.Res[WebPUpload]

	err := n.B.CallMultipart("POST", "/upload/webp", image.Name(), image, &res)
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
