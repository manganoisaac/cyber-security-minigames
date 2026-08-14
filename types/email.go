package types

type Email struct {
	ID                   string
	Title                string
	Subject              string
	To                   string
	From                 string
	Content              string
	Reason               string
	IsMalware            bool
	Timestamp            string
	CorrectPhraseIndexes []int
}
