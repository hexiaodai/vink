package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type CephInterface interface {
	ListOsds(ctx context.Context) ([]*Osd, error)
}

func NewCeph(ctx context.Context, addr string, username, password string) (CephInterface, error) {
	client := ceph{
		restyClient: resty.New().SetBaseURL(addr),
	}
	if err := client.setBearerToken(ctx, username, password); err != nil {
		return nil, fmt.Errorf("failed to set bearer token: %w", err)
	}
	return &client, nil
}

type ceph struct {
	restyClient *resty.Client
	bearerToken string
}

func makeURL(path ...string) string {
	return strings.Join(path, "/")
}

func (client *ceph) getRequest(ctx context.Context) *resty.Request {
	return client.restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		R().
		SetContext(ctx).
		SetError(&HTTPErrorResponse{})
}

func (client *ceph) getRequestWithBearerAuth(ctx context.Context) *resty.Request {
	return client.getRequest(ctx).
		SetAuthToken(client.bearerToken).
		SetHeader("Accept", "application/vnd.ceph.api.v1.0+json").
		SetHeader("Content-Type", "application/json")
}

func (client *ceph) setBearerToken(ctx context.Context, username, password string) error {
	type simple struct {
		Token string `json:"token"`
	}
	result := simple{}
	req, err := client.getRequestWithBearerAuth(ctx).
		SetBody(map[string]string{"username": username, "password": password}).
		SetResult(&result).
		Post(makeURL("api", "auth"))
	if err := checkForError(req, err); err != nil {
		return fmt.Errorf("failed to get bearer token: %w", err)
	}
	client.bearerToken = result.Token
	return nil
}

func (client *ceph) getOsdIDs(ctx context.Context) ([]int, error) {
	type simple struct {
		Osd int `json:"osd"`
	}
	result := []*simple{}
	req, err := client.getRequestWithBearerAuth(ctx).
		SetResult(&result).
		Get(makeURL("api", "osd"))
	if err := checkForError(req, err); err != nil {
		return nil, err
	}
	osds := make([]int, 0, len(result))
	for _, v := range result {
		osds = append(osds, v.Osd)
	}
	return osds, nil
}

func (client *ceph) ListOsds(ctx context.Context) ([]*Osd, error) {
	ids, err := client.getOsdIDs(ctx)
	if err != nil {
		return nil, err
	}
	output := make([]*Osd, 0, len(ids))
	for _, osd := range ids {
		result := Osd{}
		req, err := client.getRequestWithBearerAuth(ctx).
			SetResult(&result).
			Get(makeURL("api", "osd", strconv.Itoa(osd)))
		if err := checkForError(req, err); err != nil {
			return nil, err
		}
		output = append(output, &result)
	}
	return output, nil
}

type OsdMap struct {
	Osd int `json:"osd"`
	Up  int `json:"up"`
}

type OsdMetadata struct {
	Hostname             string `json:"hostname"`
	BluestoreBdevDevNode string `json:"bluestore_bdev_dev_node"`
	BluestoreBdevType    string `json:"bluestore_bdev_type"`
}

type Osd struct {
	OsdMap      *OsdMap      `json:"osd_map"`
	OsdMetadata *OsdMetadata `json:"osd_metadata"`
}

type HTTPErrorResponse struct {
	Status    string `json:"status,omitempty"`
	Detail    string `json:"detail,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func (e HTTPErrorResponse) String() string {
	var res strings.Builder
	if len(e.Status) > 0 {
		res.WriteString(e.Status)
	}
	if len(e.Detail) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.Detail)
	}
	if len(e.RequestID) > 0 {
		if res.Len() > 0 {
			res.WriteString(": ")
		}
		res.WriteString(e.RequestID)
	}
	return res.String()
}

func (e HTTPErrorResponse) NotEmpty() bool {
	return len(e.Status) > 0 || len(e.Detail) > 0 || len(e.RequestID) > 0
}

func checkForError(resp *resty.Response, err error) error {
	if err != nil {
		return err
	}

	if resp == nil {
		return fmt.Errorf("empty response")
	}

	if resp.IsError() {
		var msg string

		if e, ok := resp.Error().(*HTTPErrorResponse); ok && e.NotEmpty() {
			msg = fmt.Sprintf("%s: %s", resp.Status(), e)
		} else {
			msg = resp.Status()
		}

		return fmt.Errorf("%s", msg)
	}
	return nil
}
