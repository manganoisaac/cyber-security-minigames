package data

import "cyber-security-minigames/types"

// Credits is the end-of-game credits roll.
//
// TODO(team): the entries in "THE TEAM" below are PLACEHOLDERS. Replace each
// Name and Role with the real ones before this ships. Nothing else on the page
// needs editing to do that - add or remove CreditEntry lines freely, the roll
// is driven by its own content height so it re-times itself automatically.
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
					{Name: "[ Team member 1 ]", Role: "[ Role ]"},
					{Name: "[ Team member 2 ]", Role: "[ Role ]"},
					{Name: "[ Team member 3 ]", Role: "[ Role ]"},
					{Name: "[ Team member 4 ]", Role: "[ Role ]"},
					{Name: "[ Team member 5 ]", Role: "[ Role ]"},
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
