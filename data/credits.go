package data

import "cyber-security-minigames/types"

// Credits is the end-of-game credits roll.
//
// Entries render in the order given, so "THE TEAM" is listed by role rather
// than alphabetically. Adding or removing a few lines needs no other change;
// the roll animates for a fixed duration (--credits-roll-duration in
// features/credits.css), so if this list ever grows a lot, raise that to keep
// the reading pace.
func Credits() types.CreditsData {
	return types.CreditsData{
		Title:    "CAUGHT IN THE NET",
		Subtitle: "Cyber Security Minigames Project",

		CommissionedLabel: "Commissioned on behalf of",
		CommissionedFor:   "Ernest Foo",
		CommissionedOrg:   "Griffith School of Cyber Security",

		Groups: []types.CreditGroup{
			{
				Heading: "THE TEAM",
				Entries: []types.CreditEntry{
					{Name: "Chloe Badrick", Role: "Project Manager & Team Leader"},
					{Name: "Isaac Mangano", Role: "Lead Developer"},
					{Name: "David Pearsall", Role: "Activity & Content Designer / Assistant Developer"},
					{Name: "Jeremy Goodger", Role: "Integration & Deployment / Assistant Developer"},
					{Name: "Cody Perkins", Role: "UI/UX Designer"},
					{Name: "Taylor Campbell", Role: "UI/UX Designer & Testing Lead"},
				},
			},
			{
				Heading: "THE CASE FILE",
				Entries: []types.CreditEntry{
					{Name: "Passwords & Social Engineering", Role: "Chapter 1"},
					{Name: "Online Privacy", Role: "Chapter 2"},
					{Name: "Phishing", Role: "Chapter 3"},
					{Name: "Encryption", Role: "Chapter 4"},
					{Name: "Network Defence", Role: "Chapter 5"},
					{Name: "Digital Investigation", Role: "Chapter 6"},
				},
			},
		},

		ThanksHeading: "THANK YOU FOR PLAYING",

		ThanksParagraphs: []string{
			"Every attack in this game is one that happens for real, and every defence is one that really works.",
			"Strong, unique passwords. Think before you share. Check the sender. Encrypt what matters. Layer your defences. Follow the evidence.",
		},

		SignOff: "STAY SAFE OUT THERE",

		ReturnLabel: "Return to main menu",
		ReturnHref:  "/",
	}
}
