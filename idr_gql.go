package rapid7

import (
	"github.com/shurcooL/graphql"
)

var DefaultAssetCountMax int = 10_000

type AssetCountQuery struct {
	Organization struct {
		Assets struct {
			TotalCount graphql.Int `graphql:"totalCount"`
		} `graphql:"assets(first: $first)"`
	} `graphql:"organization(id: $id)"`
}

func (gql *GraphQLClient) AssetCount(orgID string) (*AssetCountQuery, error) {
	var res *AssetCountQuery
	err := gql.client.Query(gql.ctx, &res, map[string]any{"id": graphql.String(orgID), "first": graphql.Int(DefaultAssetCountMax)})
	if err != nil {
		return nil, err
	}
	return res, nil
}
