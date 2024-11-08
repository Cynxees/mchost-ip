package api

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	models "mchost-spot-instance/server/models"
	pb "mchost-spot-instance/server/pb"
)

func (s *Server) CreateTemplate(ctx context.Context, request pb.CreateTemplateRequest) (*pb.GetTemplateResponse, error) {
	
	spotRequestInput := &ec2.RequestSpotFleetInput{
		SpotFleetRequestConfig: &types.SpotFleetRequestConfigData{

			IamFleetRole:   aws.String("arn:aws:iam::071412439153:role/aws-ec2-spot-fleet-tagging-role"),
			TargetCapacity: aws.Int32(1),

			InstanceInterruptionBehavior: types.InstanceInterruptionBehaviorStop,
			LaunchSpecifications: []types.SpotFleetLaunchSpecification{
				{
					ImageId:      aws.String("ami-0f35248300a04b419"),
					InstanceType: types.InstanceTypeT32xlarge,
					KeyName:      aws.String("minecraft-server"),
					SecurityGroups: []types.GroupIdentifier{
						{
							GroupId: aws.String("sg-086ac6894e23bbcd2"),
						},
					},
					SubnetId: aws.String("subnet-034d02bbb1909da28"),

					IamInstanceProfile: &types.IamInstanceProfileSpecification{
						Arn: aws.String("arn:aws:iam::071412439153:instance-profile/EC2-S3-FullAccess"),
					},
				},
			},
		},
	}

	if(spotRequestInput == nil) {
		panic("Something wong")
	}

	spotInstanceTemplate := &models.SpotInstanceTemplate{
		FleetRequestId: nil,
		InstanceId: nil,
		UserId: 1,
		Name: request.Name,
		Status: "PENDING",
		InstanceType: "t2.TODO",
	}

	if err := s.Db.Create(spotInstanceTemplate).Error; err != nil {
		return nil, err
	}
	
	return &pb.GetTemplateResponse{
		Error: false,
		Code: http.StatusOK,
		Message: "Success",
		Template: &pb.SpotInstanceTemplate{
			Id: uint64(spotInstanceTemplate.ID),
			FleetRequestId: *spotInstanceTemplate.FleetRequestId,
			InstanceId: *spotInstanceTemplate.InstanceId,
			UserId: uint64(spotInstanceTemplate.UserId),
			Name: spotInstanceTemplate.Name,
			Status: spotInstanceTemplate.Status,
			InstanceType: spotInstanceTemplate.InstanceType,
			CreatedAt: timestamppb.New(spotInstanceTemplate.CreatedAt),
			UpdatedAt: timestamppb.New(spotInstanceTemplate.UpdatedAt),

		},
	}, nil
}

func (s *Server) LaunchSpotFleet(ctx context.Context) (*ec2.RequestSpotFleetOutput, error) {

	client := s.AWSManager.EC2Client

	spotRequestInput := &ec2.RequestSpotFleetInput{
		SpotFleetRequestConfig: &types.SpotFleetRequestConfigData{

			IamFleetRole:   aws.String("arn:aws:iam::071412439153:role/aws-ec2-spot-fleet-tagging-role"),
			TargetCapacity: aws.Int32(1),

			InstanceInterruptionBehavior: types.InstanceInterruptionBehaviorStop,
			LaunchSpecifications: []types.SpotFleetLaunchSpecification{
				{
					ImageId:      aws.String("ami-0f35248300a04b419"),
					InstanceType: types.InstanceTypeT32xlarge,
					KeyName:      aws.String("minecraft-server"),
					SecurityGroups: []types.GroupIdentifier{
						{
							GroupId: aws.String("sg-086ac6894e23bbcd2"),
						},
					},
					SubnetId: aws.String("subnet-034d02bbb1909da28"),

					IamInstanceProfile: &types.IamInstanceProfileSpecification{
						Arn: aws.String("arn:aws:iam::071412439153:instance-profile/EC2-S3-FullAccess"),
					},
				},
			},
		},
	}

	result, err := client.RequestSpotFleet(ctx, spotRequestInput)
	if err != nil {
		s.Logger.Fatal(err)
	}

	s.Logger.Info(result)
	return result, nil
}
