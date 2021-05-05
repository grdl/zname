package zname

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type ZonesAPI interface {
	ListHostedZones(context.Context, *route53.ListHostedZonesInput, ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)
}

type RecordsAPI interface {
	ListResourceRecordSets(context.Context, *route53.ListResourceRecordSetsInput, ...func(*route53.Options)) (*route53.ListResourceRecordSetsOutput, error)
}

type Client struct {
	ctx        context.Context
	zonesAPI   ZonesAPI
	recordsAPI RecordsAPI
}

func NewFromConfig() (*Client, error) {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := route53.NewFromConfig(cfg)
	return New(client, client), nil
}

func New(zonesAPI ZonesAPI, recordsAPI RecordsAPI) *Client {
	return &Client{
		ctx:        context.TODO(),
		zonesAPI:   zonesAPI,
		recordsAPI: recordsAPI,
	}
}

func (c *Client) GetZones() ([]Zone, error) {
	zones := make([]Zone, 0)

	params := &route53.ListHostedZonesInput{}

	var nextPage = true
	for nextPage {
		response, err := c.zonesAPI.ListHostedZones(c.ctx, params)
		if err != nil {
			return nil, err
		}

		params.Marker = response.NextMarker
		nextPage = response.IsTruncated

		for _, hostedZone := range response.HostedZones {
			zone := Zone{
				ID:   *hostedZone.Id,
				Name: *hostedZone.Name,
			}

			zones = append(zones, zone)
		}
	}

	return zones, nil
}

func (c *Client) GetRecords(zoneID string) ([]Record, error) {
	records := make([]Record, 0)

	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zoneID,
	}

	var nextPage = true
	for nextPage {
		response, err := c.recordsAPI.ListResourceRecordSets(c.ctx, params)
		if err != nil {
			return nil, err
		}

		params.StartRecordIdentifier = response.NextRecordIdentifier
		params.StartRecordName = response.NextRecordName
		params.StartRecordType = response.NextRecordType
		nextPage = response.IsTruncated

		for _, rs := range response.ResourceRecordSets {
			if record := parseRecord(rs); record != nil {
				records = append(records, *record)
			}
		}
	}

	return records, nil
}

func parseRecord(rs types.ResourceRecordSet) *Record {
	// Ignore types other than A or CNAME
	if rs.Type != types.RRTypeA && rs.Type != types.RRTypeCname {
		return nil
	}

	record := &Record{
		Name: *rs.Name,
		Type: string(rs.Type),
	}

	if rs.AliasTarget != nil {
		record.Target = *rs.AliasTarget.DNSName
		return record
	}

	if rs.ResourceRecords != nil {
		// TODO: handle multiple targets
		record.Target = *rs.ResourceRecords[0].Value
	}

	return record
}
