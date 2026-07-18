package enums

type SponsoredType string

const (
	Amount              SponsoredType = "amount"
	Equipment           SponsoredType = "equipment"
	Electronics         SponsoredType = "electronics"
	Furnitures          SponsoredType = "furnitures"
	GalleryArrangements SponsoredType = "gallery_arrangements"
	Others              SponsoredType = "others"
)

func (s SponsoredType) IsValid() bool {

	switch s {
	case Amount, Equipment, Electronics, Furnitures, GalleryArrangements, Others:
		return true
	default:
		return false
	}
}
