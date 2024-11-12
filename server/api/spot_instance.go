package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"google.golang.org/protobuf/types/known/timestamppb"

	models "mchost-ip/server/models"
	pb "mchost-ip/server/pb"
	util "mchost-ip/server/lib/util"
)

func (s *Server) CreateIp(ctx context.Context, request *pb.CreateIpRequest) (*pb.CreateIpResponse, error) {
	s.Logger.Info("Allocating Ip")
	allocatedIp, err := s.AWSManager.EC2Client.AllocateAddress(ctx, &ec2.AllocateAddressInput{
		Domain: types.DomainTypeVpc,
	})

	if err != nil {
		return nil, err
	}

	allocatedIpJson, err := json.Marshal(allocatedIp)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("Creating Ip: ", string(allocatedIpJson))

	ip := &models.Ip{
		AllocationId: *allocatedIp.AllocationId,
		OwnerId:      int(request.UserId),
		InstanceId:   nil,
		Name:         request.Name,
		Type:         "elastic",
		Region:       s.AWSManager.Config.Region,
		Address:      *allocatedIp.PublicIp,
	}

	s.Logger.Info("Creating Ip in database")
	if err := s.Db.Create(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip created in database")
	return &pb.CreateIpResponse{
		Error:        false,
		Code:         http.StatusOK,
		Message:      "Success",
		AllocationId: ip.AllocationId,
		OwnerId:      uint64(ip.OwnerId),
		InstanceId:   util.SafeString(ip.InstanceId),
		Name:         ip.Name,
		Type:         ip.Type,
		Region:       ip.Region,
		Address:      ip.Address,
	}, nil
}

func (s *Server) GetIp(ctx context.Context, request *pb.GetIpRequest) (*pb.GetIpResponse, error) {
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	return &pb.GetIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Ip: &pb.Ip{
			Id:           uint64(ip.ID),
			AllocationId: ip.AllocationId,
			InstanceId:   util.SafeString(ip.InstanceId),
			OwnerId:      uint64(ip.OwnerId),
			Name:         ip.Name,
			Type:         ip.Type,
			Region:       ip.Region,
			Address:      ip.Address,
			CreatedAt:    timestamppb.New(ip.CreatedAt),
			UpdatedAt:    timestamppb.New(ip.UpdatedAt),
		},
	}, nil
}


func (s *Server) DeleteIp(ctx context.Context, request *pb.DeleteIpRequest) (*pb.DeleteIpResponse, error){

	s.Logger.Info("Releasing Ip")
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	_, err := s.AWSManager.EC2Client.ReleaseAddress(ctx, &ec2.ReleaseAddressInput{
		AllocationId: &ip.AllocationId,
	})
	if err != nil {
		return nil, err
	}

	s.Logger.Info("Deleting Ip from database")
	if err := s.Db.Delete(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip deleted from database")
	return &pb.DeleteIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}, nil
}

func (s *Server) ReserveIp(ctx context.Context, request *pb.ReserveIpRequest) (*pb.ReserveIpResponse, error){

	s.Logger.Info("Associating Ip with spot instance")
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	id := int(request.SpotInstanceTemplateId)
	ip.SpotInstanceTemplateId = &id

	s.Logger.Info("Updating Ip in database")
	if err := s.Db.Save(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip updated in database")
	return &pb.ReserveIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		EipAllocationId: ip.AllocationId,
	}, nil
}

func (s *Server) UnreserveIp(ctx context.Context, request *pb.UnreserveIpRequest) (*pb.UnreserveIpResponse, error){

	s.Logger.Info("Disassociating Ip from instance")
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	ip.SpotInstanceTemplateId = nil

	s.Logger.Info("Updating Ip in database")
	if err := s.Db.Save(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip updated in database")
	return &pb.UnreserveIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}, nil
}

func (s *Server) UseIp(ctx context.Context, request *pb.UseIpRequest) (*pb.UseIpResponse, error){

	s.Logger.Info("Associating Ip with instance")
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	ip.InstanceId = &request.InstanceId

	s.Logger.Info("Updating Ip in database")
	if err := s.Db.Save(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip updated in database")
	return &pb.UseIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}, nil
}

func (s *Server) UnuseIp(ctx context.Context, request *pb.UnuseIpRequest) (*pb.UnuseIpResponse, error){

	s.Logger.Info("Disassociating Ip from instance")
	ip := &models.Ip{}
	if err := s.Db.Where("id = ?", request.IpId).First(ip).Error; err != nil {
		return nil, err
	}

	ip.InstanceId = nil

	s.Logger.Info("Updating Ip in database")
	if err := s.Db.Save(ip).Error; err != nil {
		return nil, err
	}

	s.Logger.Info("Ip updated in database")
	return &pb.UnuseIpResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}, nil
}