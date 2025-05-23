package resources

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elasticache"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/ekristen/aws-nuke/v3/pkg/nuke"
)

const ElasticacheUserGroupResource = "ElasticacheUserGroup"

func init() {
	registry.Register(&registry.Registration{
		Name:     ElasticacheUserGroupResource,
		Scope:    nuke.Account,
		Resource: &ElasticacheUserGroup{},
		Lister:   &ElasticacheUserGroupLister{},
	})
}

type ElasticacheUserGroupLister struct{}

func (l *ElasticacheUserGroupLister) List(_ context.Context, o interface{}) ([]resource.Resource, error) {
	opts := o.(*nuke.ListerOpts)

	svc := elasticache.New(opts.Session)
	resources := make([]resource.Resource, 0)
	var nextToken *string

	for {
		params := &elasticache.DescribeUserGroupsInput{
			MaxRecords: aws.Int64(100),
			Marker:     nextToken,
		}
		resp, err := svc.DescribeUserGroups(params)
		if err != nil {
			return nil, err
		}

		for _, userGroup := range resp.UserGroups {
			resources = append(resources, &ElasticacheUserGroup{
				svc:     svc,
				groupID: userGroup.UserGroupId,
			})
		}

		// Check if there are more results
		if resp.Marker == nil {
			break // No more results, exit the loop
		}

		// Set the nextToken for the next iteration
		nextToken = resp.Marker
	}

	return resources, nil
}

type ElasticacheUserGroup struct {
	svc     *elasticache.ElastiCache
	groupID *string
}

func (i *ElasticacheUserGroup) Remove(_ context.Context) error {
	params := &elasticache.DeleteUserGroupInput{
		UserGroupId: i.groupID,
	}

	_, err := i.svc.DeleteUserGroup(params)
	if err != nil {
		return err
	}

	return nil
}

func (i *ElasticacheUserGroup) Properties() types.Properties {
	properties := types.NewProperties()
	properties.Set("ID", i.groupID)
	return properties
}

func (i *ElasticacheUserGroup) String() string {
	return *i.groupID
}
