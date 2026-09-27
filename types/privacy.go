package types

type Post struct {
	Id     int
	Author string
	Handle string
	// Posted is the relative time shown on the post, e.g. "2h" or "1w".
	Posted string
	// Context is an optional line above the post explaining why it is on
	// this profile, e.g. a friend's post that tags the profile owner.
	Context  string
	Content  string
	Likes    int
	Comments int
	Shares   int
	Clue     *Clue

	// Analysed tracks whether the player has correctly classified this post.
	Analysed bool
	// AnalysisOptions are the four multiple-choice answers shown when the
	// player analyses the post. CorrectOption is the zero-based index of the
	// right answer.
	AnalysisOptions  []string
	CorrectOption    int
	AnalysisFeedback string
}

type Clue struct {
	AddedToFolder bool
	Key           string
	Value         string
}

type PostUI struct {
	Id               int
	Author           string
	Handle           string
	Posted           string
	Context          string
	Content          string
	Likes            int
	Comments         int
	Shares           int
	AddedToFolder    bool
	Analysed         bool
	HasClue          bool
	AnalysisOptions  []string
	CorrectOption    int
	AnalysisFeedback string
}

type ClueUI struct {
	Key   string
	Value string
}

type PrivacyUI struct {
	Posts         []PostUI
	Clues         []ClueUI
	AnalysedCount int
	TotalPosts    int
}

// SocialProfile is the owner of the feed in the Online Privacy game, shown in
// the profile header above the posts. Counts are display strings ("1,284")
// because they are only ever shown, never calculated with.
type SocialProfile struct {
	Name      string
	Handle    string
	Followers string
	Following string
}
