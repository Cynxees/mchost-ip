package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"mchost-spot-instance/server/api"
	"mchost-spot-instance/server/models"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/redis/go-redis/v9"
)

func StartSpotInstanceWorker(server *api.Server) {
	go func() {
		ctx := context.Background()

		for {
			now := time.Now().Unix()

			// Fetch tasks that are due for processing
			fleetRequestIds, err := server.Redis.ZRangeByScore(ctx, "spot_instance_queue", &redis.ZRangeBy{
				Min: "-inf",
				Max: fmt.Sprintf("%d", now),
				Offset: 0,
				Count:  10,
			}).Result()
			if err != nil {
				server.Logger.Error("Failed to fetch from Redis:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, fleetRequestId := range fleetRequestIds {
				// Fetch template by FleetRequestId
				var template models.SpotInstanceTemplate
				if err := server.Db.Where("fleet_request_id = ?", fleetRequestId).First(&template).Error; err != nil {
					server.Logger.Error("Failed to find template:", err)
					continue
				}

				// Check if the Spot Fleet request has been fulfilled
				fleetInstances, err := server.AWSManager.EC2Client.DescribeSpotFleetInstances(ctx, &ec2.DescribeSpotFleetInstancesInput{
					SpotFleetRequestId: &fleetRequestId,
				})
				if err != nil {
					server.Logger.Error("Failed to describe Spot Fleet instances:", err)
					continue
				}

				// If instances are found, update the database
				if len(fleetInstances.ActiveInstances) > 0 {
					instanceID := fleetInstances.ActiveInstances[0].InstanceId
					template.InstanceId = instanceID
					template.Status = "ACTIVE"

					if err := server.Db.Save(&template).Error; err != nil {
						server.Logger.Error("Failed to update template:", err)
						continue
					}

					// Log the instance details
					instanceDetails, _ := json.Marshal(fleetInstances)
					server.Logger.Info("Spot instance provisioned:", string(instanceDetails))

					assignElasticIP(ctx, server, *instanceID);

					// Remove the processed task from the queue
					server.Redis.ZRem(ctx, "spot_instance_queue", fleetRequestId)
				}
			}

			// Sleep for a short interval before checking again
			time.Sleep(10 * time.Second)
		}
	}()
}

func assignElasticIP(ctx context.Context, server *api.Server, instanceID string) error {

	eipAllocationID := "eipalloc-0d3131e17bfd77974"

	// Associate the Elastic IP with the instance
	_, err := server.AWSManager.EC2Client.AssociateAddress(ctx, &ec2.AssociateAddressInput{
		InstanceId:   aws.String(instanceID),
		AllocationId: aws.String(eipAllocationID),
		AllowReassociation: aws.Bool(true),
	})

	if err != nil {
		return fmt.Errorf("failed to associate Elastic IP: %w", err)
	}

	server.Logger.Info("Elastic IP associated with instance:", instanceID)
	return nil
}