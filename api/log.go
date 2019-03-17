package api

import (
	"context"
	"io"

	"github.com/levigross/grequests"
)

func (c *Client) GetLogs(ctx context.Context) (io.ReadCloser, error) {
	resp, err := grequests.Get(c.BaseURL+"/logs", &grequests.RequestOptions{
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	return resp.RawResponse.Body, nil
}
