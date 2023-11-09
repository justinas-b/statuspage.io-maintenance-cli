package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func PointersOf(v interface{}) interface{} {
	in := reflect.ValueOf(v)
	out := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(in.Type().Elem())), in.Len(), in.Len())
	for i := 0; i < in.Len(); i++ {
		out.Index(i).Set(in.Index(i).Addr())
	}
	return out.Interface()
}

func (c *Client) newRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("OAuth %s", c.APIKey))
	return req
}

func (c *Client) doRequest(method string, url string, body io.Reader) *http.Response {
	req := c.newRequest(method, url, body)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

func parseResponse[T any](body io.Reader) (T, error) {

	var object T
	err := json.NewDecoder(body).Decode(&object)
	if err != nil {
		return object, fmt.Errorf("invalid API response %s: %w", body, err)
	}

	return object, nil
}
