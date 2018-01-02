package arn

import (
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"sort"
	"testing"
)

func TestAWSIPRangeDoc(t *testing.T) {
	d, err := AWSIPRangesDoc(http.DefaultClient)
	if err != nil {
		t.Fatalf("error getting IP ranges document: %v", err)
	}
	assert.NotZero(t, d.SyncToken, "SyncToken is empty")
	assert.NotZero(t, d.CreateDate, "CreateDate not set")
	assert.NotZero(t, d.Prefixes, "Prefixes is not populated")
}

func TestAWSRegions(t *testing.T) {
	var regions = []string{
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
	}
	r, err := AWSRegions(http.DefaultClient)
	if err != nil {
		t.Fatalf("error getting regions: %v", err)
	}

	trans := cmp.Transformer("Sort", func(in []string) []string {
		out := make([]string, len(in), len(in))
		copy(out, in) // Copy input to avoid mutating it
		sort.Strings(out)
		return out
	})
	assert.True(t, cmp.Equal(regions, r, trans), "Regions slice not as expected: %+v", r)
}
