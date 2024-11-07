package api

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func (s *Server) RegisterSpotInstance(ctx context.Context) (*any, error) {

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		s.Logger.Fatal(err)
	}

	client := ec2.NewFromConfig(cfg)

	spotRequestInput := &ec2.RequestSpotInstancesInput{
		InstanceCount: aws.Int32(1), // Number of instances
		LaunchSpecification: &types.RequestSpotLaunchSpecification{
			ImageId:      aws.String("ami-xxxxxxxx"),   // Replace with a valid AMI ID
			InstanceType: types.InstanceTypeT2Micro,    // Replace with your preferred instance type
			KeyName:      aws.String("your-key-pair"),  // Replace with your key pair name
			SecurityGroups: []string{
				"default", // Replace with your security group
			},
			SubnetId: aws.String("subnet-xxxxxxxx"), // Replace with a valid Subnet ID
			IamInstanceProfile: &types.IamInstanceProfileSpecification{
				Arn: aws.String("arn:aws:iam::123456789012:instance-profile/your-iam-instance-profile"),
			},
		},
		SpotPrice: aws.String("0.02"), // Set the maximum price you're willing to pay per hour
		Type:      types.SpotInstanceTypeOneTime,
	}

	result, err := client.RequestSpotInstances(context.TODO(), spotRequestInput)
	if err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info(result)
	return nil, nil
}
