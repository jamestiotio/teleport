/*
Copyright 2023 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package services

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gravitational/trace"
	"github.com/stretchr/testify/require"

	"github.com/gravitational/teleport/api/types"
)

func TestNewAWSNodeFromEC2Instance(t *testing.T) {
	isBadParameterErr := func(tt require.TestingT, err error, i ...any) {
		require.True(tt, trace.IsBadParameter(err), "expected bad parameter, got %v", err)
	}

	makeEC2Instance := func(fn func(*ec2Types.Instance)) ec2Types.Instance {
		s := ec2Types.Instance{
			PrivateDnsName:   aws.String("my-private-dns.compute.aws"),
			InstanceId:       aws.String("i-123456789abcedf"),
			VpcId:            aws.String("vpc-abcd"),
			SubnetId:         aws.String("subnet-123"),
			PrivateIpAddress: aws.String("172.31.1.1"),
			Tags: []ec2Types.Tag{
				{
					Key:   aws.String("MyTag"),
					Value: aws.String("MyTagValue"),
				},
			},
		}
		fn(&s)
		return s
	}

	for _, tt := range []struct {
		name             string
		ec2Instance      ec2Types.Instance
		awsCloudMetadata *types.AWSInfo
		errCheck         require.ErrorAssertionFunc
		expectedServer   types.Server
	}{
		{
			name:        "valid",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: require.NoError,
			expectedServer: &types.ServerV2{
				Kind:    "node",
				Version: "v2",
				SubKind: "openssh-ec2-ice",
				Metadata: types.Metadata{
					Labels: map[string]string{
						"account-id":               "1234567889012",
						"region":                   "us-east-1",
						"MyTag":                    "MyTagValue",
						"teleport.dev/instance-id": "i-123456789abcedf",
					},
					Namespace: "default",
				},
				Spec: types.ServerSpecV2{
					Addr:     "172.31.1.1:22",
					Hostname: "my-private-dns.compute.aws",
					CloudMetadata: &types.CloudMetadata{
						AWS: &types.AWSInfo{
							AccountID:   "1234567889012",
							InstanceID:  "i-123456789abcedf",
							Region:      "us-east-1",
							VPCID:       "vpc-abcd",
							SubnetID:    "subnet-123",
							Integration: "myintegration",
						},
					},
				},
			},
		},
		{
			name: "instance metadata generated labels are not replaced by instance tags",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {
				i.Tags = append(i.Tags, ec2Types.Tag{
					Key:   aws.String("region"),
					Value: aws.String("evil"),
				})
				i.Tags = append(i.Tags, ec2Types.Tag{
					Key:   aws.String("account-id"),
					Value: aws.String("evil"),
				})
			}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: require.NoError,
			expectedServer: &types.ServerV2{
				Kind:    "node",
				Version: "v2",
				SubKind: "openssh-ec2-ice",
				Metadata: types.Metadata{
					Labels: map[string]string{
						"account-id":               "1234567889012",
						"region":                   "us-east-1",
						"MyTag":                    "MyTagValue",
						"teleport.dev/instance-id": "i-123456789abcedf",
					},
					Namespace: "default",
				},
				Spec: types.ServerSpecV2{
					Addr:     "172.31.1.1:22",
					Hostname: "my-private-dns.compute.aws",
					CloudMetadata: &types.CloudMetadata{
						AWS: &types.AWSInfo{
							AccountID:   "1234567889012",
							InstanceID:  "i-123456789abcedf",
							Region:      "us-east-1",
							VPCID:       "vpc-abcd",
							SubnetID:    "subnet-123",
							Integration: "myintegration",
						},
					},
				},
			},
		},
		{
			name: "missing ec2 private dns name",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {
				i.PrivateDnsName = nil
			}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name: "missing ec2 instance id",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {
				i.InstanceId = nil
			}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name: "missing ec2 vpc id",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {
				i.VpcId = nil
			}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name: "missing ec2 private ip address",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {
				i.PrivateDnsName = nil
			}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name:        "missing account id",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {}),
			awsCloudMetadata: &types.AWSInfo{
				Region:      "us-east-1",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name:        "missing region",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID:   "1234567889012",
				Integration: "myintegration",
			},
			errCheck: isBadParameterErr,
		},
		{
			name:        "missing integration name",
			ec2Instance: makeEC2Instance(func(i *ec2Types.Instance) {}),
			awsCloudMetadata: &types.AWSInfo{
				AccountID: "1234567889012",
				Region:    "us-east-1",
			},
			errCheck: isBadParameterErr,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewAWSNodeFromEC2Instance(tt.ec2Instance, tt.awsCloudMetadata)
			tt.errCheck(t, err)
			if err != nil {
				return
			}

			require.Empty(t, cmp.Diff(tt.expectedServer, s, cmpopts.IgnoreFields(types.ServerV2{}, "Metadata.Name")))
		})
	}
}
