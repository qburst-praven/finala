package resources

import (
	"errors"
	"fmt"

	"finala/collector"
	"finala/collector/aws/common"
	"finala/collector/aws/pricing"
	"finala/collector/aws/register"
	"finala/collector/config"

	awsClient "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	awspricing "github.com/aws/aws-sdk-go/service/pricing"
	log "github.com/sirupsen/logrus"
)

// ElasticIPClientDescriptor is an interface defining the aws ec2 client
type ElasticIPClientDescriptor interface {
	DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
}

// ElasticIPManager will hold the elastic ip manger strcut
type ElasticIPManager struct {
	client             ElasticIPClientDescriptor
	awsManager         common.AWSManager
	servicePricingCode string
	rateCode           string
	Name               collector.ResourceIdentifier
}

// DetectedElasticIP defines the detected AWS elastic ip
type DetectedElasticIP struct {
	Region        string
	Metric        string
	IP            string
	PricePerHour  float64
	PricePerMonth float64
	Tag           map[string]string
}

func init() {
	register.Registry("elasticip", NewElasticIPManager)
}

// NewElasticIPManager implements AWS GO SDK
func NewElasticIPManager(awsManager common.AWSManager, client interface{}) (common.ResourceDetection, error) {
	if client == nil {
		client = ec2.New(awsManager.GetSession())
	}

	ec2Client, ok := client.(ElasticIPClientDescriptor)
	if !ok {
		return nil, errors.New("invalid ec2 elasticip client")
	}

	return &ElasticIPManager{
		client:             ec2Client,
		awsManager:         awsManager,
		servicePricingCode: "AmazonVPC",
		rateCode:           "6YS6EN2CT7",
		Name:               awsManager.GetResourceIdentifier("elastic_ip"),
	}, nil
}

// Detect checks if elastic ips is under utilized
func (ei *ElasticIPManager) Detect(metrics []config.MetricConfig) (interface{}, error) {
	metric := metrics[0]

	log.WithFields(log.Fields{
		"region":   ei.awsManager.GetRegion(),
		"resource": "elastic ips",
	}).Info("starting to analyze resource")

	ei.awsManager.GetCollector().CollectStart(ei.Name)

	elasticIPs := []DetectedElasticIP{}

	priceFIlters := ei.getPricingFilterInput()
	// Get elastic ip pricing
	price, err := ei.awsManager.GetPricingClient().GetPrice(priceFIlters, ei.rateCode, ei.awsManager.GetRegion())
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"rate_code":     ei.rateCode,
			"region":        ei.awsManager.GetRegion(),
			"price_filters": priceFIlters,
		}).Error("could not get elastic ip price")
		ei.awsManager.GetCollector().CollectError(ei.Name, err)

		return elasticIPs, err
	}

	// Getting all elastic ip addressess
	ips, err := ei.describeAddressess()
	if err != nil {
		ei.awsManager.GetCollector().CollectError(ei.Name, err)
		return elasticIPs, err
	}

	for _, ip := range ips {
		if ip.PrivateIpAddress == nil && ip.AssociationId == nil && ip.InstanceId == nil && ip.NetworkInterfaceId == nil {
			tagsData := map[string]string{}
			if err == nil {
				for _, tag := range ip.Tags {
					tagsData[*tag.Key] = *tag.Value
				}
			}

			eIP := DetectedElasticIP{
				Region:        ei.awsManager.GetRegion(),
				Metric:        metric.Description,
				IP:            *ip.PublicIp,
				PricePerHour:  price,
				PricePerMonth: price * collector.TotalMonthHours,
				Tag:           tagsData,
			}

			ei.awsManager.GetCollector().AddResource(collector.EventCollector{
				ResourceName: ei.Name,
				Data:         eIP,
			})

			elasticIPs = append(elasticIPs, eIP)
		}
	}

	ei.awsManager.GetCollector().CollectFinish(ei.Name)

	return elasticIPs, nil
}

// getPricingFilterInput returns the elastic ip price filters.
func (ei *ElasticIPManager) getPricingFilterInput() awspricing.GetProductsInput {
	region := ei.awsManager.GetRegion()
	pricingClient := ei.awsManager.GetPricingClient()

	regionPrefixForUsageType, err := pricingClient.GetRegionPrefix(region)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"region": region,
		}).Error("Could not get region prefix for usagetype")
		return awspricing.GetProductsInput{ServiceCode: &ei.servicePricingCode}
	}

	// Get region info for location filter
	regionInfo, found := pricing.RegionsInfo[region]
	if !found {
		log.WithField("region", region).Error("Could not get region info for location filter")
		return awspricing.GetProductsInput{ServiceCode: &ei.servicePricingCode}
	}

	// The format is regionPrefix + "PublicIPv4:IdleAddress"
	regionSpecificUsageType := fmt.Sprintf("%sPublicIPv4:IdleAddress", regionPrefixForUsageType)
	log.WithFields(log.Fields{
		"region":     region,
		"prefix":     regionPrefixForUsageType,
		"usage_type": regionSpecificUsageType,
	}).Debug("Constructed usage type for pricing filter")

	return awspricing.GetProductsInput{
		ServiceCode: &ei.servicePricingCode,
		Filters: []*awspricing.Filter{
			{
				Type:  awsClient.String("TERM_MATCH"),
				Field: awsClient.String("location"),
				Value: awsClient.String(regionInfo.FullName),
			},
			{
				Type:  awsClient.String("TERM_MATCH"),
				Field: awsClient.String("group"),
				Value: awsClient.String("VPCPublicIPv4Address"),
			},
			{
				Type:  awsClient.String("TERM_MATCH"),
				Field: awsClient.String("usagetype"),
				Value: awsClient.String(regionSpecificUsageType),
			},
			{
				Type:  awsClient.String("TERM_MATCH"),
				Field: awsClient.String("termType"),
				Value: awsClient.String("OnDemand"),
			},
		},
	}
}

// describeAddressess returns list of elastic ips addresses
func (ei *ElasticIPManager) describeAddressess() ([]*ec2.Address, error) {
	input := &ec2.DescribeAddressesInput{}

	resp, err := ei.client.DescribeAddresses(input)
	if err != nil {
		log.WithField("error", err).Error("could not describe elastic ips addresses")
		return nil, err
	}

	return resp.Addresses, nil
}
