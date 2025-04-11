package offerservice

import "example.com/webrtc-practice/internal/domain/service"

type OfferServiceImpl struct {
	id string
}

func NewOfferService() service.OfferService {
	return &OfferServiceImpl{
		id: "",
	}
}

func (o *OfferServiceImpl) SetOffer(id string) {
	o.id = id
}

func (o *OfferServiceImpl) GetOffer() string {
	return o.id
}

func (o *OfferServiceImpl) ClearOffer() {
	o.id = ""
}

func (o *OfferServiceImpl) IsOffer() bool {
	return o.id != ""
}

func (o *OfferServiceImpl) IsOfferID(id string) bool {
	return o.id == id
}