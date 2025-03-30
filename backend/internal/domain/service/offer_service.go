package service

type OfferService interface {
	SetOffer(id string)
	GetOffer() string
	ClearOffer()
	IsOffer() bool
	IsOfferID(id string) bool
}