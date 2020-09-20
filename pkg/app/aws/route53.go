package aws

import (
	"encoding/json"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"os"
	"strconv"
)

func ListHostedZones() (*HostedZones, error) {
	out, err := bash.RunAndLogRead("aws route53 list-hosted-zones")
	if err != nil {
		return nil, err
	}

	var hz *HostedZones
	if err := json.Unmarshal(out, &hz); err != nil {
		return nil, errors.WithStack(err)
	}

	return hz, nil
}

func ListHostedZonesByProfile(profile string) (*HostedZones, error) {
	out, err := bash.RunAndLogRead("aws route53 list-hosted-zones --profile " + profile)
	if err != nil {
		return nil, err
	}

	var hz *HostedZones
	if err := json.Unmarshal(out, &hz); err != nil {
		return nil, errors.WithStack(err)
	}

	return hz, nil
}

func ListRecordsByHostedZone(hostedZone *HostedZone) (*ResourceRecordSets, error) {
	out, err := bash.RunAndLogRead("aws route53 list-resource-record-sets --hosted-zone-id " + hostedZone.ID)
	if err != nil {
		return nil, err
	}

	var rs *ResourceRecordSets
	if err := json.Unmarshal(out, &rs); err != nil {
		return nil, errors.WithStack(err)
	}

	return rs, nil
}

func ListRecordsByHostedZoneByProfile(hostedZone *HostedZone, profile string) (*ResourceRecordSets, error) {
	out, err := bash.RunAndLogRead("aws route53 list-resource-record-sets --hosted-zone-id " + hostedZone.ID + " --profile " + profile)
	if err != nil {
		return nil, err
	}

	var rs *ResourceRecordSets
	if err := json.Unmarshal(out, &rs); err != nil {
		return nil, errors.WithStack(err)
	}

	return rs, nil
}

func ListAllRecords() (*ResourceRecordSets, error) {
	hzs, err := ListHostedZones()
	if err != nil {
		return nil, err
	}
	rs := make([]*ResourceRecordSet, 0, 0)
	for _, hz := range hzs.Items {
		r, err := ListRecordsByHostedZone(hz)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r.Items...)
	}
	return &ResourceRecordSets{rs}, nil
}

func ListAllRecordsByProfile(profile string) (*ResourceRecordSets, error) {
	hzs, err := ListHostedZonesByProfile(profile)
	if err != nil {
		return nil, err
	}
	rs := make([]*ResourceRecordSet, 0, 0)
	for _, hz := range hzs.Items {
		r, err := ListRecordsByHostedZoneByProfile(hz, profile)
		if err != nil {
			return nil, err
		}
		for _, item := range r.Items {
			item.Profile = profile
			item.HostedZone = hz
			rs = append(rs, item)
		}
	}
	return &ResourceRecordSets{rs}, nil
}

type HostedZones struct {
	Items []*HostedZone `json:"HostedZones"`
}

type HostedZone struct {
	ID              string `json:"Id"`
	Name            string `json:"Name"`
	CallerReference string `json:"CallerReference"`
	Config          struct {
		Comment     string `json:"Comment"`
		PrivateZone bool   `json:"PrivateZone"`
	} `json:"Config"`
	ResourceRecordSetCount int `json:"ResourceRecordSetCount"`
}

type ResourceRecordSets struct {
	Items []*ResourceRecordSet `json:"ResourceRecordSets"`
}

func (s *ResourceRecordSets) RenderTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"HostedZone", "PrivateZone", "Name", "Type", "TTL", "Record", "Profile"})
	table.SetBorder(false)
	for _, item := range s.Items {
		for _, record := range item.ResourceRecords {
			table.Append([]string{item.HostedZone.Name, strconv.FormatBool(item.HostedZone.Config.PrivateZone), item.Name, item.Type, strconv.Itoa(item.TTL), record.Value, item.Profile})
		}
	}
	table.Render()
}

type ResourceRecordSet struct {
	Name            string `json:"Name"`
	Type            string `json:"Type"`
	TTL             int    `json:"TTL"`
	ResourceRecords []struct {
		Value string `json:"Value"`
	} `json:"ResourceRecords"`
	Profile    string
	HostedZone *HostedZone
}
