package arn

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	IPRangesURL = "https://ip-ranges.amazonaws.com/ip-ranges.json"
	TimeFormat  = "2006-01-02-15-04-05" //"2017-12-21-20-12-10" Jan 2 15:04:05 2006 MST
)

// AWSIPRanges represents the AWS ip-ranges.json document.
type AWSIPRanges struct {
	SyncToken     string `json:"syncToken"`
	CreateDateStr string `json:"createDate"`
	CreateDate    time.Time
	Prefixes      []struct {
		IPPrefix string `json:"ip_prefix"`
		IP       net.IP
		Network  *net.IPNet
		Region   string `json:"region"`
		Service  string `json:"service"`
	} `json:"prefixes"`
	Ipv6Prefixes []struct {
		Ipv6Prefix string `json:"ipv6_prefix"`
		IP         net.IP
		Network    *net.IPNet
		Region     string `json:"region"`
		Service    string `json:"service"`
	} `json:"ipv6_prefixes"`
}

// AWSIPRangesDoc reads the ip-ranges.json document from https://ip-ranges.amazonaws.com/ip-ranges.json
// and returns an AWSIPRanges type.
func AWSIPRangesDoc(client *http.Client) (awsip AWSIPRanges, err error) {
	req, e := http.NewRequest("GET", IPRangesURL, nil)
	if e != nil {
		err = e
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if client.Timeout == 0 {
		client.Timeout = time.Second * 10
	}
	response, e := client.Do(req)
	if e != nil {
		err = e
		return
	}
	buf, e := ioutil.ReadAll(response.Body)
	if e != nil {
		err = e
		return
	}
	defer response.Body.Close()
	err = json.Unmarshal(buf, &awsip)
	if err != nil {
		return
	}
	awsip.CreateDate, err = time.Parse(TimeFormat, awsip.CreateDateStr)
	if err != nil {
		return
	}
	for i := range awsip.Prefixes {
		awsip.Prefixes[i].IP, awsip.Prefixes[i].Network, err = net.ParseCIDR(awsip.Prefixes[i].IPPrefix)
	}
	for i := range awsip.Ipv6Prefixes {
		awsip.Ipv6Prefixes[i].IP, awsip.Ipv6Prefixes[i].Network, err = net.ParseCIDR(awsip.Ipv6Prefixes[i].Ipv6Prefix)
	}
	return
}

// AWSRegions returns a slice of valid AWS region names.
// If the nil is passed for the client then a static list of well known regions is returned
// rather than making an HTTP request.
func AWSRegions(client *http.Client) ([]string, error) {
	if client == nil {
		return []string{
			"GLOBAL",
			"us-west-1",
			"ap-southeast-2",
			"us-east-2",
			"ap-northeast-1",
			"ap-northeast-2",
			"ap-south-1",
			"ap-southeast-1",
			"eu-central-1",
			"eu-west-1",
			"us-east-1",
			"sa-east-1",
			"us-gov-east-1",
			"us-west-2",
			"eu-west-2",
			"ca-central-1",
			"eu-west-3",
			"us-gov-west-1",
			"cn-north-1",
			"cn-northwest-1",
		}, nil
	}
	awsip, err := AWSIPRangesDoc(client)
	if err != nil {
		return []string{}, err
	}

	r := make([]string, 0, len(awsip.Prefixes))
	m := make(map[string]bool)

	for _, val := range awsip.Prefixes {
		if _, ok := m[val.Region]; !ok {
			m[val.Region] = true
			r = append(r, val.Region)
		}
	}
	return r, nil
}
