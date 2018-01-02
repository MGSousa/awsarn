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

func AWSIPRangeDoc() (awsip AWSIPRanges, err error) {
	req, e := http.NewRequest("GET", IPRangesURL, nil)
	if e != nil {
		err = e
		return
	}
	req.Header.Set("Content-Type", "application/json")
	var HTTPClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, e := HTTPClient.Do(req)
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

func AWSRegions() ([]string, error) {
	awsip, err := AWSIPRangeDoc()
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
