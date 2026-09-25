package data

import "cyber-security-minigames/types"

// PrivacyInstructions is shown in the modal when the Online Privacy game starts.
const PrivacyInstructions = "Some of these posts are giving away more than they should. Click through the feed and collect every post that leaks real information into your Clues Found list. For the best experience, play in full screen."

// PrivacyIntro is the pre-instruction overview shown before PrivacyInstructions.
// It frames what the player is about to learn.
const PrivacyIntro = "Online Privacy - Every post adds to your digital footprint. On their own the details people share seem harmless, but combined they can reveal where you work, who you know and when you'll be away. Eugene built a picture of Angela from what she posted publicly. Read through her feed and collect the posts that give away real personal information - and see how much a stranger could learn."

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
			Key:   "employer",
			Value: "Software Developer at MetaCortex Industries, South Brisbane",
		},
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
			Key:   "manager",
			Value: "David Walker",
		},
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
			Key:   "favourite store",
			Value: "Paws & Co, Grey Street (regular after-work visits)",
		},
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
			Value: "Project Helios",
		},
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
	},
}
