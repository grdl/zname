package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Client struct {
	ctx     context.Context
	route53 *route53.Client
}

func New() (*Client, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	return &Client{
		ctx:     ctx,
		route53: route53.NewFromConfig(cfg),
	}, nil

}

func (c *Client) GetZones() ([]types.HostedZone, error) {
	zones := make([]types.HostedZone, 0)

	params := &route53.ListHostedZonesInput{}

	var nextPage = true
	for nextPage {
		response, err := c.route53.ListHostedZones(c.ctx, params)
		if err != nil {
			return nil, err
		}

		params.Marker = response.NextMarker
		nextPage = response.IsTruncated

		zones = append(zones, response.HostedZones...)
	}

	return zones, nil
}

func (c *Client) GetRecords(zoneID string) ([]types.ResourceRecordSet, error) {
	records := make([]types.ResourceRecordSet, 0)

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zoneID,
	}

	var nextPage = true
	for nextPage {
		response, err := c.route53.ListResourceRecordSets(c.ctx, params)
		if err != nil {
			return nil, err
		}

		params.StartRecordIdentifier = response.NextRecordIdentifier
		params.StartRecordName = response.NextRecordName
		params.StartRecordType = response.NextRecordType
		nextPage = response.IsTruncated

		records = append(records, response.ResourceRecordSets...)
	}

	return records, nil
}
