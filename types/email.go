package types

type Email struct {
	Title              string
	Subject            string
	Content            string
	IsMalware          bool
	Timestamp          string
	CorrectWordIndexes []int
}
