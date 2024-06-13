package rapid7

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/shurcooL/graphql"
)

type GraphQLClient struct {
	http   *http.Client
	url    *url.URL
	client *graphql.Client
	ctx    context.Context
}

type transport struct {
	rt     http.RoundTripper
	apiKey string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("x-api-key", t.apiKey)
	return t.rt.RoundTrip(req)
}

func newGraphQLTransport(apiKey string) http.RoundTripper {
	return &transport{
		rt:     http.DefaultTransport,
		apiKey: apiKey,
	}
}

func NewGraphQLClient(region, apiKey string) (*GraphQLClient, error) {
	urlS := fmt.Sprintf("https://%s.api.insight.rapid7.com/graphql", strings.ToLower(region))
	u, err := url.Parse(urlS)
	if err != nil {
		return nil, err
	}
	hc := &http.Client{Transport: newGraphQLTransport(apiKey)}
	c := graphql.NewClient(u.String(), hc)
	ctx := context.TODO()
	gql := &GraphQLClient{
		http:   hc,
		client: c,
		ctx:    ctx,
		url:    u,
	}
	return gql, nil
}
