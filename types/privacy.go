package types

type Post struct {
	Id      int
	Author  string
	Content string
	Clue    *Clue
}

type Clue struct {
	AddedToFolder bool
	Key           string
	Value         string
}

type PostUI struct {
	Id            int
	Author        string
	Content       string
	AddedToFolder bool
}

type ClueUI struct {
	Key   string
	Value string
}

type PrivacyUI struct {
	Posts []PostUI
	Clues []ClueUI
}
