// Package arn provides utilities for manipulating Amazon Resource Names: http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
package arn

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

//arn:partition:service:region:account-id:resource
//arn:partition:service:region:account-id:resourcetype/resource
//arn:partition:service:region:account-id:resourcetype:resource

// ARN is a struct representing structure of an ARN.
type ARN struct {
	Partition    string
	Service      string
	Region       string
	AccountID    string
	ResourceType string
	resourceSep  string
	Resource     string
}

// Parse validates an ARN string and returns an ARN type.
// An http.Client can be provided to validate the region against the latest regions queries from an AWS endpoint.
// To disable this HTTP request pass nil as the client and the region will be
// validated against a static set of well known AWS regions.
func Parse(arn string, client *http.Client) (ARN, error) {
	if !Valid(arn, client) {
		return ARN{}, errors.New("ARN not valid")
	}
	parts := strings.SplitN(arn, ":", 6)
	a := ARN{
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		AccountID: parts[4],
	}
	if strings.Contains(parts[5], ":") {
		r := strings.SplitN(parts[5], ":", 2)
		a.ResourceType = r[0]
		a.resourceSep = ":"
		a.Resource = r[1]
	} else if strings.Contains(parts[5], "/") {
		r := strings.SplitN(parts[5], "/", 2)
		a.ResourceType = r[0]
		a.resourceSep = "/"
		a.Resource = r[1]
	} else {
		a.Resource = parts[5]
	}
	return a, nil
}

// String returns the AWS standard string representation of an ARN type.
func (a *ARN) String() string {
	r := a.Resource
	if a.resourceSep != "" {
		r = a.ResourceType + a.resourceSep + a.Resource
	}
	return fmt.Sprintf("arn:%s:%s:%s:%s:%s",
		a.Partition,
		a.Service,
		a.Region,
		a.AccountID,
		r,
	)
}

// Valid checks the format and content of the ARN string are valid.
// The http.Client is required to check the region in the ARN is valid against an authoritative AWS data source.
// To disable the HTTP check for the region pass nil as the client and the region will be checked against
// a static set of well known AWS regions only.
func Valid(arn string, client *http.Client) bool {
	if strings.Contains(strings.SplitN(arn, "/", 2)[0], " ") {
		return false
	}
	if strings.Count(arn, ":") < 5 {
		return false
	}
	parts := strings.SplitN(arn, ":", 6)
	// 2nd field must be "aws" or start "aws-"
	if !strings.HasPrefix(parts[1], "aws-") && parts[1] != "aws" {
		return false
	}
	// 3rd field must not be null
	if parts[2] == "" {
		return false
	}
	// 4th valid region or null or *
	if !ValidRegion(parts[3], client) && parts[3] != "" && parts[3] != "*" {
		return false
	}
	// 5th account number (12 digit) or null or *
	if parts[4] != "" && utf8.RuneCountInString(parts[4]) != 12 && parts[4] != "*" {
		return false
	}
	if _, err := strconv.Atoi(parts[4]); parts[4] != "" && parts[4] != "*" && err != nil {
		return false
	}
	return strings.HasPrefix(arn, "arn:")
}

// ValidRegion checks the region is a valid AWS region.
// The http.Client is required to check the region against an authoritative AWS data source.
// To disable the HTTP check for the region pass nil as the client and the region will be checked against
// a static set of well known AWS regions only.
func ValidRegion(region string, client *http.Client) bool {
	regions, err := AWSRegions(client)
	if err != nil {
		return false
	}
	for _, r := range regions {
		if region == r {
			return true
		}
	}
	return false
}
