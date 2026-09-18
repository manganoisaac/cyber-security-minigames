package data

import "cyber-security-minigames/types"

// PrivacyInstructions is shown in the modal when the Online Privacy game starts.
const PrivacyInstructions = "Some of these posts are giving away more than they should. Click through the feed and collect every post that leaks real information into your Clues Found list. For the best experience, play in full screen."

func GetPostData() []types.Post {
	copiedData := make([]types.Post, len(PostData))
	for i := range len(PostData) {
		copiedData[i] = PostData[i]
		if PostData[i].Clue != nil {
			copiedClue := *PostData[i].Clue
			copiedData[i].Clue = &copiedClue
		}
	}
	return copiedData
}

var PostData = []types.Post{
	{
		Id:      1,
		Author:  "Angela Bennett",
		Content: "Software Developer at MetaCortex Industries • South Brisbane",
		Clue: &types.Clue{
			Key:   "employer",
			Value: "Software Developer at MetaCortex Industries, South Brisbane",
		},
	},
	{
		Id:      8,
		Author:  "Angela Bennett",
		Content: "Finally tried that new ramen place. 10/10, absolutely going back.",
		Clue:    nil,
	},
	{
		Id:      3,
		Author:  "Angela Bennett",
		Content: "Couldn't ask for a better team. Even David Walker bought us coffee for once!",
		Clue: &types.Clue{
			Key:   "manager",
			Value: "David Walker",
		},
	},
	{
		Id:      6,
		Author:  "Angela Bennett",
		Content: "Quick stop at Paws & Co on Grey Street after work. Max is running dangerously low on Chicken Bites.",
		Clue: &types.Clue{
			Key:   "favourite store",
			Value: "Paws & Co, Grey Street (regular after-work visits)",
		},
	},
	{
		Id:      9,
		Author:  "Angela Bennett",
		Content: "Started this show last night and accidentally watched six episodes. Send help.",
		Clue:    nil,
	},
	{
		Id:      2,
		Author:  "Angela Bennett",
		Content: "Another late night on Project Helios. At least we're finally getting close to the finish line!",
		Clue: &types.Clue{
			Key:   "current project",
			Value: "Project Helios",
		},
	},
	{
		Id:      5,
		Author:  "Angela Bennett",
		Content: "Someone clearly thinks he owns the couch. Meet Max, the actual boss of the house.",
		Clue: &types.Clue{
			Key:   "pet",
			Value: "Max (dog)",
		},
	},
	{
		Id:      10,
		Author:  "Angela Bennett",
		Content: "Sunset over the bay tonight. Needed this after a long week.",
		Clue:    nil,
	},
	{
		Id:      4,
		Author:  "Angela Bennett",
		Content: "Only two weeks until Japan! ✈️ Can't wait to disappear from work from 14–28 October.",
		Clue: &types.Clue{
			Key:   "travel dates",
			Value: "Away from work 14–28 October (Japan trip)",
		},
	},
	{
		Id:      7,
		Author:  "Sophie Miller",
		Content: "Perfect weekend at Mount Tamborine with Angela. Even Max behaved for a few photos!",
		Clue: &types.Clue{
			Key:   "friend",
			Value: "Sophie Miller (recent Mount Tamborine trip)",
		},
	},
}
