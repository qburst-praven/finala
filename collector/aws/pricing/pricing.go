package pricing

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	awsClient "github.com/aws/aws-sdk-go/aws"
	awsPricing "github.com/aws/aws-sdk-go/service/pricing"
	log "github.com/sirupsen/logrus"
)

// ErrRegionNotFound when a region is not found
var ErrRegionNotFound = errors.New("region was not found as part of the regionsInfo map")

// RegionInfo will hold data about a region pricing options
type RegionInfo struct {
	FullName string
	Prefix   string
}

var RegionsInfo = map[string]RegionInfo{
	"us-east-2":      {FullName: "US East (Ohio)", Prefix: "USE2"},
	"us-east-1":      {FullName: "US East (N. Virginia)", Prefix: ""},
	"us-west-1":      {FullName: "US West (N. California)", Prefix: "USW1"},
	"us-west-2":      {FullName: "US West (Oregon)", Prefix: "USW2"},
	"ap-east-1":      {FullName: "Asia Pacific (Hong Kong)", Prefix: "APE1"},
	"ap-south-1":     {FullName: "Asia Pacific (Mumbai)", Prefix: ""},
	"ap-northeast-3": {FullName: "Asia Pacific (Osaka-Local)", Prefix: "APN3"},
	"ap-northeast-2": {FullName: "Asia Pacific (Seoul)", Prefix: "APN2"},
	"ap-southeast-1": {FullName: "Asia Pacific (Singapore)", Prefix: "APS1"},
	"ap-southeast-2": {FullName: "Asia Pacific (Sydney)", Prefix: "APS2"},
	"ap-northeast-1": {FullName: "Asia Pacific (Tokyo)", Prefix: "APN1"},
	"ca-central-1":   {FullName: "Canada (Central)", Prefix: "CAN1"},
	"cn-north-1":     {FullName: "China (Beijing)", Prefix: ""},
	"cn-northwest-1": {FullName: "China (Ningxia)", Prefix: ""},
	"eu-central-1":   {FullName: "EU (Frankfurt)", Prefix: "EUC1"},
	"eu-west-1":      {FullName: "EU (Ireland)", Prefix: "EUW1"},
	"eu-west-2":      {FullName: "EU (London)", Prefix: "EUW2"},
	"eu-west-3":      {FullName: "EU (Paris)", Prefix: "EUW3"},
	"eu-south-1":     {FullName: "EU (Milan)", Prefix: "EUS1"},
	"eu-north-1":     {FullName: "EU (Stockholm)", Prefix: "EUN1"},
	"sa-east-1":      {FullName: "South America (Sao Paulo)", Prefix: "SAE1"},
	"us-gov-east-1":  {FullName: "AWS GovCloud (US-East)", Prefix: "UGE1"},
	"us-gov-west-1":  {FullName: "AWS GovCloud (US)", Prefix: "UGW1"},
	"af-south-1":     {FullName: "Africa (Cape Town)", Prefix: "AFS1"},
	"me-south-1":     {FullName: "Middle East (Bahrain)", Prefix: "MES1"},
}

// PricingClientDescreptor is an interface defining the aws pricing client
type PricingClientDescreptor interface {
	GetProducts(*awsPricing.GetProductsInput) (*awsPricing.GetProductsOutput, error)
}

// PricingManager Pricing
type PricingManager struct {
	client         PricingClientDescreptor
	region         string
	priceResponses map[uint64]float64
}

// PricingResponse describ the response of AWS pricing
type PricingResponse struct {
	Products PricingProduct `json:"product"`
	Terms    PricingTerms   `json:"terms"`
}

// PricingProduct describe the product details
type PricingProduct struct {
	SKU string `json:"sku"`
}

// PricingTerms describe the product terms
type PricingTerms struct {
	OnDemand map[string]*PricingOfferTerm `json:"OnDemand"`
}

// PricingOfferTerm describe the product offer terms
type PricingOfferTerm struct {
	SKU             string                    `json:"sku"`
	PriceDimensions map[string]*PriceRateCode `json:"priceDimensions"`
}

// PriceRateCode describe the product price
type PriceRateCode struct {
	Unit         string            `json:"unit"`
	PricePerUnit PriceCurrencyCode `json:"pricePerUnit"`
}

// PriceCurrencyCode Descrive the pricing currency
type PriceCurrencyCode struct {
	USD string `json:"USD"`
}

// NewPricingManager implements AWS GO SDK
func NewPricingManager(client PricingClientDescreptor, region string) *PricingManager {
	log.Debug("Initializing aws pricing SDK client")
	return &PricingManager{
		client:         client,
		region:         region,
		priceResponses: make(map[uint64]float64),
	}
}

// GetPrice returns the price for the given filters and rate code.
func (p *PricingManager) GetPrice(filters awsPricing.GetProductsInput, rateCode string, region string) (float64, error) {
	// Add location filter
	regionInfo, found := RegionsInfo[region]
	if !found {
		return 0, fmt.Errorf("region info not found for %s", region)
	}

	filters.Filters = append(filters.Filters, &awsPricing.Filter{
		Type:  awsClient.String("TERM_MATCH"),
		Field: awsClient.String("location"),
		Value: awsClient.String(regionInfo.FullName),
	})

	// Get products
	products, err := p.client.GetProducts(&filters)
	if err != nil {
		return 0, err
	}

	if len(products.PriceList) == 0 {
		return 0, fmt.Errorf("no products found for the given filters")
	}

	// Get the first product
	product := products.PriceList[0]

	// Unmarshal the product into our PricingResponse struct
	var pricingResponse PricingResponse
	productJSON, err := json.Marshal(product)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal product: %v", err)
	}

	if err := json.Unmarshal(productJSON, &pricingResponse); err != nil {
		return 0, fmt.Errorf("failed to unmarshal product: %v", err)
	}

	// Get the price
	var price float64
	var priceFound bool

	// Get the price from the terms
	for _, term := range pricingResponse.Terms.OnDemand {
		// Get the price from the price dimensions
		for _, priceDimension := range term.PriceDimensions {
			// Get the price from the price per unit
			if priceDimension.PricePerUnit.USD != "" {
				price, err = strconv.ParseFloat(priceDimension.PricePerUnit.USD, 64)
				if err != nil {
					return 0, fmt.Errorf("failed to parse price: %v", err)
				}
				priceFound = true
				break
			}
		}
		if priceFound {
			break
		}
	}

	if !priceFound {
		return 0, fmt.Errorf("no price found for the given filters")
	}

	return price, nil
}

// GetRegionPrefix will return the prefix for a
// pricing filter value according to a given region.
// For example:
// Region: "us-east-2" prefix will be: "USE2-"
func (p *PricingManager) GetRegionPrefix(region string) (string, error) {
	var prefix string
	regionInfo, found := RegionsInfo[region]
	if !found {
		return prefix, ErrRegionNotFound
	}

	switch regionInfo.Prefix {
	case "":
		prefix = ""
	default:
		prefix = fmt.Sprintf("%s-", RegionsInfo[region].Prefix)
	}
	return prefix, nil
}

// Add this method to expose the raw client
func (p *PricingManager) RawClient() PricingClientDescreptor {
	return p.client
}
