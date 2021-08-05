package cache

import (
	"context"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockZonesAPI func(context.Context, *route53.ListHostedZonesInput, ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)

func (m mockZonesAPI) ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
	return m(ctx, params, optFns...)
}

func TestAPI(t *testing.T) {
	tests := []struct {
		zonesAPI func() ZonesAPI
	}{
		{
			zonesAPI: func() ZonesAPI {
				return mockZonesAPI(func(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
					return &route53.ListHostedZonesOutput{
						HostedZones: []types.HostedZone{},
					}, nil
				})
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			client := New(test.zonesAPI(), nil)

			zones, err := client.GetZones()
			require.NoError(t, err)
			assert.Len(t, zones, 0)
		})
	}
}
