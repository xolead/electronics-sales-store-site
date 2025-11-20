package models

type ProductBuilder struct {
	product ProductResponse
}

func NewProductBuilder() *ProductBuilder {
	return &ProductBuilder{
		product: ProductResponse{
			Images: make([]string, 0),
		},
	}
}

func (b *ProductBuilder) WithID(id int) *ProductBuilder {
	b.product.ID = id
	return b
}

func (b *ProductBuilder) WithName(name string) *ProductBuilder {
	b.product.Name = name
	return b
}

func (b *ProductBuilder) WithDescription(description string) *ProductBuilder {
	b.product.Description = description
	return b
}

func (b *ProductBuilder) WithParameters(parameters string) *ProductBuilder {
	b.product.Parameters = parameters
	return b
}

func (b *ProductBuilder) WithCount(count int) *ProductBuilder {
	b.product.Count = count
	return b
}

func (b *ProductBuilder) WithPrice(price int) *ProductBuilder {
	b.product.Price = price
	return b
}

func (b *ProductBuilder) AddImage(url string) *ProductBuilder {
	b.product.Images = append(b.product.Images, url)
	return b
}

func (b *ProductBuilder) Build() ProductResponse {
	return b.product
}
