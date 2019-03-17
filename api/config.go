package api

import (
	"fmt"

	T "github.com/Dreamacro/clash/tunnel"
	"github.com/jqs7/clash_cli/model"
	"github.com/levigross/grequests"
)

type Client struct {
	BaseURL string
}

func (c *Client) GetConfigs() (*model.Config, error) {
	resp, err := grequests.Get(c.BaseURL+"/configs", nil)
	if err != nil {
		return nil, err
	}
	rst := &model.Config{}
	if err := resp.JSON(rst); err != nil {
		return nil, err
	}
	return rst, nil
}

func (c *Client) UpdateMode(mode T.Mode) error {
	resp, err := grequests.Patch(c.BaseURL+"/configs", &grequests.RequestOptions{
		JSON: &model.Config{
			Mode: &mode,
		},
	})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("%d: %s", resp.StatusCode, resp.String())
	}
	return nil
}
