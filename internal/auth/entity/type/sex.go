package custom_type

type Sex uint

const (
	SexOther Sex = iota
	SexFemale
	SexMale
)

func (s Sex) String() string {
	return [...]string{"Other", "Male", "Female"}[s]
}

func (s Sex) EnumIndex() uint {
	return uint(s)
}

func (s Sex) IsSexOther() bool {
	return s == SexOther
}

func (s Sex) IsSexMale() bool {
	return s == SexMale
}

func (s Sex) IsSexFemale() bool {
	return s == SexFemale
}
