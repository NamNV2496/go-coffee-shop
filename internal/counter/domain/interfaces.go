package domain

type ProductDomainServer interface {
	GetProductByIdOrName(id int32, name string) ([]Item, error)
}
