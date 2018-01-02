# AWS ARN
[![GoDoc](https://godoc.org/github.com/jcmturner/awsarn?status.svg)](https://godoc.org/github.com/jcmturner/awsarn) [![Go Report Card](https://goreportcard.com/badge/github.com/jcmturner/awsarn)](https://goreportcard.com/report/github.com/jcmturner/awsarn)

The AWS ARN library provides Go representation for AWS ARNs and useful functions to parse and validate ARN strings.

Reference: https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html

To get the package, execute:
```
go get github.com/jcmturner/awsarn
```
To import this package, add the following line to your code:
```go
import "github.com/jcmturner/awsarn"

```

## Usage
For full API documentation refer to the GoDoc link above.

### Parse
Creating an ARN instance from an ARN string:
```go
import "github.com/jcmturner/awsarn"

arnStr := "arn:aws:iam::012345678912:user/testuser"
arn, err := awsarn.Parse(arnStr, http.DefaultClient)

fmt.Printf(`Partition: %s\n
Service %s\n
Region %s\n
AccountID %s\n
ResourceType %s\n
Resource %s\n`,
        arn.Service,
        arn.Region,
        arn.AccountID,
        arn.ResourceType,
        arn.Resource,
        )
```

### String
Get the AWS ARN in string format:
```go
arnStr := arn.String()
fmt.Println(arnStr)
```

### Validate
```go
arnStr := "arn:aws:iam::012345678912:user/testuser"
if awsarn.Valid(arnStr, http.DefaultClient) {
        fmt.Println("ARN valid")
} else {
        fmt.Println("ARN invalid")
}
```

## Online Reference Data
Validation can use the AWS [ip-ranges.json]("https://ip-ranges.amazonaws.com/ip-ranges.json") document as reference data 
for the valid values of AWS regions. AWS bring new regions online and therefore this is the best, most up to date source 
of valid regions.

To disable the use of this online resource pass a nil pointer where *http.Client is required as an argument.

If you need to define the use of an HTTP proxy to reach the internet configure the transport on an http.Client passed 
as an argument.
```go
proxyURL, err := url.Parse("http://proxyName:proxyPort")
client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}

arnStr := "arn:aws:iam::012345678912:user/testuser"
if awsarn.Valid(arnStr, client) {
        fmt.Println("ARN valid")
} else {
        fmt.Println("ARN invalid")
}
```