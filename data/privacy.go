package data

import "cyber-security-minigames/types"

// PrivacyInstructions is shown in the modal when the Online Privacy game starts.
const PrivacyInstructions = `Read all 10 posts on Angela's profile. Select ANALYSE on each post and decide what useful personal information, if any, it reveals. Some posts contain information an attacker could use, while others do not reveal anything significant.

Use the feedback after each answer to understand why the information matters. Analyse all 10 posts to complete the investigation.`

// PrivacyIntro is the pre-instruction overview shown before PrivacyInstructions.
// It frames what the player is about to learn.
const PrivacyIntro = `Online privacy is about controlling how much information about you is available to other people. Every public post, profile or interaction contributes to your digital footprint (the trail of information you leave behind online).

A single detail may seem harmless, but information from different posts and websites can be combined to reveal where you work, who you know, your interests, routines, travel plans and other personal details. Attackers can use this information to make social-engineering attacks more believable.

Eugene has been researching Angela using information she shared publicly. Analyse her profile and see how much a stranger could learn.`

// PrivacyProfile is the account whose feed the player reads through.
var PrivacyProfile = types.SocialProfile{
	Name:      "Angela Bennett",
	Handle:    "@angela.bennett",
	Followers: "1,284",
	Following: "312",
}

func GetPostData() []types.Post {
	copiedData := make([]types.Post, len(PostData))
	for i := range len(PostData) {
		copiedData[i] = PostData[i]
		if PostData[i].Clue != nil {
			copiedClue := *PostData[i].Clue
			copiedData[i].Clue = &copiedClue
		}
		copiedData[i].AnalysisOptions = append([]string(nil), PostData[i].AnalysisOptions...)
	}
	return copiedData
}

var PostData = []types.Post{
	{
		Id:       1,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "2h",
		Likes:    42,
		Comments: 8,
		Shares:   1,
		Content:  "Software Developer at MetaCortex Industries • South Brisbane",
		Clue: &types.Clue{
			Key:   "employer and location",
			Value: "Software Developer at MetaCortex Industries, South Brisbane",
		},
		AnalysisOptions: []string{
			"Her upcoming travel dates",
			"Her employer, job role and workplace area",
			"Her pet's name",
			"No significant personal information",
		},
		CorrectOption:    1,
		AnalysisFeedback: "Angela has revealed where she works, what she does and the area she works in. Those details can make a targeted workplace message much more convincing.",
	},
	{
		Id:       8,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "5h",
		Likes:    31,
		Comments: 6,
		Shares:   0,
		Content:  "Finally tried that new ramen place. 10/10, absolutely going back.",
		Clue:     nil,
		AnalysisOptions: []string{
			"Her exact home address",
			"Her employer and job title",
			"A reliable daily routine and location",
			"No significant personal information",
		},
		CorrectOption:    3,
		AnalysisFeedback: "This post shares a personal preference, but it does not name the restaurant, reveal a location, or establish a useful routine an attacker could reliably exploit.",
	},
	{
		Id:       3,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "1d",
		Likes:    57,
		Comments: 12,
		Shares:   2,
		Content:  "Couldn't ask for a better team. Even David Walker bought us coffee for once!",
		Clue: &types.Clue{
			Key:   "work contact",
			Value: "David Walker",
		},
		AnalysisOptions: []string{
			"Her travel dates",
			"Her pet's name",
			"The name of a real workplace contact",
			"No significant personal information",
		},
		CorrectOption:    2,
		AnalysisFeedback: "Naming David Walker gives an attacker a real workplace contact they could impersonate in a believable message to Angela.",
	},
	{
		Id:       6,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "2d",
		Likes:    23,
		Comments: 4,
		Shares:   0,
		Content:  "Quick stop at Paws & Co on Grey Street after work. Max is running dangerously low on Chicken Bites.",
		Clue: &types.Clue{
			Key:   "routine and location",
			Value: "Paws & Co, Grey Street (regular after-work visits)",
		},
		AnalysisOptions: []string{
			"A regular after-work location and routine",
			"Her current work project",
			"Her manager's email address",
			"No significant personal information",
		},
		CorrectOption:    0,
		AnalysisFeedback: "Angela has named a specific place she visits after work and implied that she goes there regularly. That can reveal part of her routine and likely movements.",
	},
	{
		Id:       9,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "3d",
		Likes:    19,
		Comments: 7,
		Shares:   0,
		Content:  "Started this show last night and accidentally watched six episodes. Send help.",
		Clue:     nil,
		AnalysisOptions: []string{
			"Her work schedule",
			"No significant personal information",
			"Her home address",
			"Her account password",
		},
		CorrectOption:    1,
		AnalysisFeedback: "This post reveals an entertainment preference, but nothing specific enough to expose Angela's location, work, relationships, schedule or other useful personal details.",
	},
	{
		Id:       2,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "4d",
		Likes:    36,
		Comments: 5,
		Shares:   1,
		Content:  "Another late night on Project Helios. At least we're finally getting close to the finish line!",
		Clue: &types.Clue{
			Key:   "current project",
			Value: "Project HELIOS",
		},
		AnalysisOptions: []string{
			"Her pet's name",
			"Her holiday destination",
			"No significant personal information",
			"The name of a current work project",
		},
		CorrectOption:    3,
		AnalysisFeedback: "Angela has publicly named Project HELIOS and linked herself to it. An attacker can use that project name to make a work-themed message feel legitimate.",
	},
	{
		Id:       5,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "5d",
		Likes:    128,
		Comments: 24,
		Shares:   6,
		Content:  "Someone clearly thinks he owns the couch. Meet Max, the actual boss of the house.",
		Clue: &types.Clue{
			Key:   "pet",
			Value: "Max (dog)",
		},
		AnalysisOptions: []string{
			"Her workplace location",
			"Her pet's name",
			"Her current project",
			"No significant personal information",
		},
		CorrectOption:    1,
		AnalysisFeedback: "Angela has revealed her pet's name. Details like pet names are useful personal context and are sometimes reused in passwords, password hints or security questions.",
	},
	{
		Id:       10,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "6d",
		Likes:    64,
		Comments: 9,
		Shares:   2,
		Content:  "Sunset over the bay tonight. Needed this after a long week.",
		Clue:     nil,
		AnalysisOptions: []string{
			"No significant personal information",
			"Her exact current location",
			"Her employer and job role",
			"Her Japan travel dates",
		},
		CorrectOption:    0,
		AnalysisFeedback: "The post is vague and does not identify a specific place, routine or personal detail that would meaningfully help an attacker.",
	},
	{
		Id:       4,
		Author:   "Angela Bennett",
		Handle:   "@angela.bennett",
		Posted:   "1w",
		Likes:    88,
		Comments: 21,
		Shares:   3,
		Content:  "Only two weeks until Japan! ✈️ Can't wait to disappear from work from 14–28 October.",
		Clue: &types.Clue{
			Key:   "travel dates",
			Value: "Away from work 14–28 October (Japan trip)",
		},
		AnalysisOptions: []string{
			"Her current work project",
			"The name of a workplace contact",
			"The exact period she will be away from work",
			"No significant personal information",
		},
		CorrectOption:    2,
		AnalysisFeedback: "Angela has announced both her destination and the exact dates she will be away from work. An attacker now knows when she is unavailable and can reference the trip in a targeted message.",
	},
	{
		Id:       7,
		Author:   "Sophie Miller",
		Handle:   "@sophie.miller",
		Posted:   "2w",
		Context:  "Tagged Angela Bennett",
		Likes:    47,
		Comments: 11,
		Shares:   2,
		Content:  "Perfect weekend at Mount Tamborine with Angela. Even Max behaved for a few photos!",
		Clue: &types.Clue{
			Key:   "friend",
			Value: "Sophie Miller (recent Mount Tamborine trip)",
		},
		AnalysisOptions: []string{
			"A friend or close social connection",
			"Her workplace password",
			"Her Japan travel dates",
			"No significant personal information",
		},
		CorrectOption:    0,
		AnalysisFeedback: "The tag links Angela to Sophie Miller and shows a recent shared trip. Public relationships can give an attacker names and context for impersonation or social engineering.",
	},
}
