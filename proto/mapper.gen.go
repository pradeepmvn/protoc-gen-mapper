//////////////////////////////////////////////////////////////////////////////////////////////
// Code generated by protoc-gen-mapper. .
// DO NOT EDIT.
//////////////////////////////////////////////////////////////////////////////////////////////

package product

import "strconv"
import "fmt"

const ProductName = "product.name"
const ProductDescription = "product.description"
const ProductPriceDetails = "product.priceDetails"
const ProductStarRatingStars = "product.starrating.stars"
const ProductStarRatingCount = "product.starrating.count"
const ProductStarRatingDetailsSomething = "product.starrating.details.something"
const ProductStarRatingDetailsNothing = "product.starrating.details.nothing"
const ProductQuery = "product.query"
const ProductPageNumber = "product.pageNumber"
const ProductResultPerPage = "product.resultPerPage"
const ProductIndicator = "product.indicator"

// ToMap Convert a struct into a Map
func (p *Product) ToMap() map[string]string {
	m := make(map[string]string)
	m[ProductName] = p.Name
	m[ProductDescription] = p.Description
	m[ProductPriceDetails] = p.PriceDetails
	if p.StarRating != nil {
		m[ProductStarRatingStars] = p.StarRating.Stars
		m[ProductStarRatingCount] = strconv.Itoa(int(p.StarRating.Count))
		if p.StarRating.Details != nil {
			m[ProductStarRatingDetailsSomething] = p.StarRating.Details.Something
			m[ProductStarRatingDetailsNothing] = p.StarRating.Details.Nothing
		}
	}
	m[ProductQuery] = p.Query
	m[ProductPageNumber] = fmt.Sprintf("%f", p.PageNumber)
	m[ProductResultPerPage] = strconv.Itoa(int(p.ResultPerPage))
	m[ProductIndicator] = strconv.FormatBool(p.Indicator)
	return m
}

// FromMap Convert a Map into a Struct
func FromMap(m map[string]string) *Product {
	var p = new(Product)
	p.Name = m[ProductName]
	p.Description = m[ProductDescription]
	p.PriceDetails = m[ProductPriceDetails]
	if p.StarRating != nil {
		p.StarRating.Stars = m[ProductStarRatingStars]
		iG, _ := strconv.Atoi(m[ProductStarRatingCount])
		p.StarRating.Count = int32(iG)
		if p.StarRating.Details != nil {
			p.StarRating.Details.Something = m[ProductStarRatingDetailsSomething]
			p.StarRating.Details.Nothing = m[ProductStarRatingDetailsNothing]
		}
	}
	p.Query = m[ProductQuery]
	fK, _ := strconv.ParseFloat(m[ProductPageNumber], 64)
	p.PageNumber = fK
	iF, _ := strconv.Atoi(m[ProductResultPerPage])
	p.ResultPerPage = int32(iF)
	bp, _ := strconv.ParseBool(m[ProductIndicator])
	p.Indicator = bp
	return p
}
