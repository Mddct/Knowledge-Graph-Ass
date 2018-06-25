package model

type Profile struct {
	Name     string
	Link     string
	ImageSrc string
	Year     string
	ID       string
	//	Director string
	Abstract string
}

func (p Profile) GetID() string {
	return p.ID
}
