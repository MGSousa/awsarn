// Package arn provides utilities for manipulating Amazon Resource Names: http://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
package awsarn

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	DELIMITER = ":"
	MAXLENGHT = 6
)

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
	parts := strings.SplitN(arn, DELIMITER, MAXLENGHT)

	a := ARN{
		Partition: parts[1],
		Service:   parts[2],
		Region:    parts[3],
		AccountID: parts[4],
	}

	if strings.Contains(parts[5], DELIMITER) {
		r := strings.SplitN(parts[5], DELIMITER, 2)
		a.ResourceType = r[0]
		a.resourceSep = DELIMITER
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
	var (
		validationErrors int
	)
	t := NewTerminal()

	if strings.Contains(strings.SplitN(arn, "/", 2)[0], " ") {
		return false
	}
	if strings.Count(arn, DELIMITER) < 5 {
		fmt.Println("invalid ARN format: check syntax")
		validationErrors++
	}
	parts := strings.SplitN(arn, DELIMITER, MAXLENGHT)

	// 2nd field must be "aws" or start "aws-"
	if parts[1] != "aws" && parts[1] != "aws-cn" && parts[1] != "aws-us-gov" {
		t.highlight(parts, 1)
		fmt.Println(" [x] partition must be aws, aws-cn or aws-us-gov")
		validationErrors++
	}

	// service resource field must not be null
	if parts[2] == "" {
		t.highlight(parts, 2)
		fmt.Println(" [x] resource is not specified")
		validationErrors++
	}

	// check if service region is a valid region or a null or *
	if !ValidRegion(parts[3], client) && parts[3] != "" && parts[3] != "*" {
		t.highlight(parts, 3)
		fmt.Println(" [x] region is not valid")
		validationErrors++
	}

	// check if account number null/empty or *
	if parts[4] != "" && parts[4] != "*" {
		// if not then check if is has 12 digits
		if utf8.RuneCountInString(parts[4]) != 12 {
			t.highlight(parts, 4)
			fmt.Println(" [x] account ID number is not valid, needs to have 12 digits")
			validationErrors++
		}
		if _, err := strconv.Atoi(parts[4]); err != nil {
			fmt.Println(err)
			validationErrors++
		}
	}

	if validationErrors > 0 {
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
