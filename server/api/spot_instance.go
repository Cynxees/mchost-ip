package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/types/known/timestamppb"

	models "mchost-spot-instance/server/models"
	pb "mchost-spot-instance/server/pb"
)

func (s *Server) CreateTemplate(ctx context.Context, request *pb.CreateTemplateRequest) (*pb.GetTemplateResponse, error) {

  s.Logger.Info("Creating Template")
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
			Id: uint64(spotInstanceTemplate.Id),
			UserId: uint64(spotInstanceTemplate.UserId),
			Name: spotInstanceTemplate.Name,
			Status: spotInstanceTemplate.Status,
			InstanceType: spotInstanceTemplate.InstanceType,
			CreatedAt: timestamppb.New(spotInstanceTemplate.CreatedAt),
			UpdatedAt: timestamppb.New(spotInstanceTemplate.UpdatedAt),
		},
	}, nil
}

func (s *Server) LaunchSpotFleet(ctx context.Context, request *pb.LaunchTemplateRequest) (*ec2.RequestSpotFleetOutput, error) {

  result, err := s.GetTemplate(ctx, &pb.GetTemplateRequest{SpotInstanceTemplateId: request.SpotInstanceTemplateId});
  template := result.Template

  if(template.Status == "ACTIVE") {
    return nil, errors.New("Template is already active")
  }

	client := s.AWSManager.EC2Client

	spotRequestInput := &ec2.RequestSpotFleetInput{
		SpotFleetRequestConfig: &types.SpotFleetRequestConfigData{

			IamFleetRole:   aws.String("arn:aws:iam::071412439153:role/aws-ec2-spot-fleet-tagging-role"),
			TargetCapacity: aws.Int32(1),

			InstanceInterruptionBehavior: types.InstanceInterruptionBehaviorStop,
			LaunchSpecifications: []types.SpotFleetLaunchSpecification{
				{
					ImageId:      aws.String("ami-0f35248300a04b419"),
					InstanceType: types.InstanceTypeT3Micro,
					KeyName:      aws.String("minecraft-server"),
					SecurityGroups: []types.GroupIdentifier{
						{
							GroupId: aws.String("sg-086ac6894e23bbcd2"),
						},
					},

					IamInstanceProfile: &types.IamInstanceProfileSpecification{
						Arn: aws.String("arn:aws:iam::071412439153:instance-profile/EC2-S3-FullAccess"),
					},
				},
			},
		},
	}

	fleetRequest, err := client.RequestSpotFleet(ctx, spotRequestInput)
	if err != nil {
		s.Logger.Fatal(err)
	}

  if err := s.Db.Model(&models.SpotInstanceTemplate{}).
    Where("id = ?", request.SpotInstanceTemplateId).
    UpdateColumns(map[string] interface{}{
      "fleet_request_id": fleetRequest.SpotFleetRequestId,
      "status": "ACTIVE",
    }).Error; err != nil {
    return nil, err
    }

	s.Logger.Info(fleetRequest)

	delay := time.Now().Add(20 * time.Second).Unix()
	err = s.Redis.ZAdd(ctx, "spot_instance_queue", redis.Z{
		Score: float64(delay),
		Member: fleetRequest.SpotFleetRequestId,
	}).Err()
	if err != nil {
		return nil, errors.New("Failed to push to queue");
	}

	return fleetRequest, nil
}

func (s *Server) GetTemplate (ctx context.Context, request *pb.GetTemplateRequest) (*pb.GetTemplateResponse, error) {
  
  template := &models.SpotInstanceTemplate{}
  if err := s.Db.Where("id = ?", request.SpotInstanceTemplateId).First(template).Error; err != nil {
    return nil, err
  }

  if template.FleetRequestId != nil {

    //TODO: maybe redundant
    fleetRequest, err := s.AWSManager.EC2Client.DescribeSpotFleetRequests(ctx, &ec2.DescribeSpotFleetRequestsInput{
      SpotFleetRequestIds: []string{*template.FleetRequestId},
    })
    if err != nil {
      return nil, err
    }

    if len(fleetRequest.SpotFleetRequestConfigs) > 0 {

      instances, err := s.AWSManager.EC2Client.DescribeSpotFleetInstances(ctx, &ec2.DescribeSpotFleetInstancesInput{
        SpotFleetRequestId: template.FleetRequestId,
      })
      if err != nil {
        return nil, err
      }

      if len(instances.ActiveInstances) > 0 {
        firstInstanceId := instances.ActiveInstances[0].InstanceId
        template.InstanceId = firstInstanceId
        s.Logger.Info("First instance ID:", *firstInstanceId)
      } else {
          s.Logger.Warn("No active instances found for fleet request:", *template.FleetRequestId)
      }
    }
  }

  return &pb.GetTemplateResponse{
    Error: false,
    Code: http.StatusOK,
    Message: "Success",
    Template: &pb.SpotInstanceTemplate{
      Id: uint64(template.ID),
      FleetRequestId: func() string {
          if template.FleetRequestId != nil {
              return *template.FleetRequestId
          }
          return ""
      }(),
      InstanceId: func() string {
          if template.InstanceId != nil {
              return *template.InstanceId
          }
          return ""
      }(),
      UserId: uint64(template.UserId),
      Name: template.Name,
      Status: template.Status,
      InstanceType: template.InstanceType,
      CreatedAt: timestamppb.New(template.CreatedAt),
      UpdatedAt: timestamppb.New(template.UpdatedAt),
    },
  }, nil
}

func (s *Server) StopTemplate (ctx context.Context, request *pb.StopTemplateRequest) (*pb.StopTemplateResponse, error) {
  
  template := &models.SpotInstanceTemplate{}
  if err := s.Db.Where("id = ?", request.SpotInstanceTemplateId).First(template).Error; err != nil {
    return nil, err
  }

  if template.Status != "ACTIVE" {
    return nil, errors.New("Template is not active")
  }

  if template.FleetRequestId == nil {
    return nil, errors.New("Fleet request ID is nil")
  }

  _, err := s.AWSManager.EC2Client.CancelSpotFleetRequests(ctx, &ec2.CancelSpotFleetRequestsInput{
    TerminateInstances: aws.Bool(true),
    SpotFleetRequestIds: []string{*template.FleetRequestId},
  })

  if err != nil {
    return nil, err
  }

	if err := s.uploadWorldToS3(ctx, *template.InstanceId); err != nil {
		s.Logger.Error("Failed to upload world to S3:", err)
		return nil, errors.New("Failed to upload world to S3")
	}

  if err := s.Db.Model(&models.SpotInstanceTemplate{}).
    Where("id = ?", request.SpotInstanceTemplateId).
    UpdateColumns(map[string] interface{}{
      "fleet_request_id": nil,
			"instance_id": nil,
      "status": "PENDING",
    }).Error; err != nil {
    return nil, err
  }
  
  return &pb.StopTemplateResponse{
    Error: false,
    Code: http.StatusOK,
    Message: "Success",
  }, nil
}

// uploadWorldToS3 connects to the EC2 instance and uploads the Minecraft world to S3
func (s *Server) uploadWorldToS3(ctx context.Context, instanceId string) error {
	// Fetch instance details to get public IP
	instanceDetails, err := s.AWSManager.EC2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceId},
	})
	if err != nil {
		return err
	}

	if len(instanceDetails.Reservations) == 0 || len(instanceDetails.Reservations[0].Instances) == 0 {
		return errors.New("instance not found")
	}
	instance := instanceDetails.Reservations[0].Instances[0]
	publicIP := instance.PublicIpAddress

	// SSH into the instance
	sshClient, err := s.connectSSH(*publicIP)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	// Upload world files to S3
	cmd := fmt.Sprintf("aws s3 sync /home/ubuntu/minecraft-prominence-server/world s3://mchost-%s --delete", *instance.Placement.AvailabilityZone)
	// cmd := fmt.Sprintf("aws s3 sync /homemchost-%s", *instance.Placement.AvailabilityZone)
	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	if err := session.Run(cmd); err != nil {
		return fmt.Errorf("failed to run S3 sync command: %v", err)
	}

	fmt.Println("S3 Upload Output:", stdoutBuf.String())
	return nil
}

// connectSSH creates an SSH connection to the EC2 instance
func (s *Server) connectSSH(host string) (*ssh.Client, error) {

	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", host, 22)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH: %w", err)
	}
	return client, nil
}