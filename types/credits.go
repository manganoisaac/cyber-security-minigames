package types

// CreditEntry is one person in the credits roll. Role may be empty, for an
// entry that is only a name.
type CreditEntry struct {
	Name string
	Role string
}

// CreditGroup is a headed block of the roll, e.g. "DEVELOPMENT" followed by
// the people who did it. Entries render in the order given.
type CreditGroup struct {
	Heading string
	Entries []CreditEntry
}

// CreditsData is the whole roll, top to bottom.
type CreditsData struct {
	Title    string
	Subtitle string

	// CommissionedLabel and CommissionedFor are split so the label can be
	// styled as a quiet lead-in to the name it introduces.
	CommissionedLabel string
	CommissionedFor   string
	CommissionedOrg   string

	Groups []CreditGroup

	ThanksHeading    string
	ThanksParagraphs []string

	SignOff string

	ReturnLabel string
	ReturnHref  string
}
