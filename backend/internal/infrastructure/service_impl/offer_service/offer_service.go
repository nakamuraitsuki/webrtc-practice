package offerservice

type OfferService struct {
	id string
}

func NewOfferService() OfferService {
	return OfferService{
		id : "",
	}
}

func (o *OfferService) SetOffer(id string) {
	o.id = id
}

func (o *OfferService) GetOffer() string {
	return o.id
}

func (o *OfferService) ClearOffer() {
	o.id = ""
}

func (o *OfferService) IsOffer() bool {
	return o.id != ""
}

func (o *OfferService) IsOfferID(id string) bool {
	return o.id == id
}