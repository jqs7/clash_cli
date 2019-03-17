package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/jqs7/clash_cli/model"
	"github.com/levigross/grequests"
)

func (c *Client) GetProxies() (*model.Proxies, error) {
	resp, err := grequests.Get(c.BaseURL+"/proxies", nil)
	if err != nil {
		return nil, err
	}
	rst := &model.Proxies{}
	if err := resp.JSON(rst); err != nil {
		return nil, err
	}
	return rst, nil
}

func (c *Client) UpdateProxy(group, proxy string) error {
	resp, err := grequests.Put(c.BaseURL+"/proxies/"+group, &grequests.RequestOptions{
		JSON: struct {
			Name string `json:"name"`
		}{Name: proxy},
	})
	if err != nil {
		return err
	}
	if !resp.Ok {
		return fmt.Errorf("%d: %s", resp.StatusCode, resp.String())
	}
	return nil
}

var (
	DelayTestError   = errors.New("503: An error occurred in the delay test")
	DelayTestTimeout = errors.New("408: Timeout")
)

func (c *Client) GetDelay(ctx context.Context, proxy string) (*model.Delay, error) {
	resp, err := grequests.Get(c.BaseURL+"/proxies/"+proxy+"/delay", &grequests.RequestOptions{
		Context: ctx,
		Params: map[string]string{
			"url":     "https://gstatic.com/generate_204",
			"timeout": "5000",
		},
	})
	if err != nil {
		return nil, err
	}
	if !resp.Ok {
		switch resp.StatusCode {
		case 503:
			return nil, DelayTestError
		case 408:
			return nil, DelayTestTimeout
		}
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, resp.String())
	}
	delay := &model.Delay{}
	if err := resp.JSON(delay); err != nil {
		return nil, err
	}
	return delay, nil
}
