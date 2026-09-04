package data

import "cyber-security-minigames/types"

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
		Author:  "Alice Wonder",
		Content: "Just finished setting up my new smart home system. Everything is connected and voice-controlled now!",
		Clue: &types.Clue{
			Key:   "device name",
			Value: "EchoDot LivingRoom",
		},
	},
	{
		Id:      2,
		Author:  "Bob Builder",
		Content: "Hot take: pineapple on pizza is actually good. Fight me.",
		Clue:    nil,
	},
	{
		Id:      3,
		Author:  "Charlie Dev",
		Content: "Shipped the new feature to production today. No bugs... yet.",
		Clue: &types.Clue{
			Key:   "api endpoint",
			Value: "/internal/debug/logs",
		},
	},
	{
		Id:      4,
		Author:  "Diana Art",
		Content: "Check out my latest digital painting! Took me 3 hours.",
		Clue:    nil,
	},
	{
		Id:      5,
		Author:  "Eve Sec",
		Content: "Reminder: never share your passwords or two-factor codes with anyone.",
		Clue:    nil,
	},
	{
		Id:      6,
		Author:  "Frank Music",
		Content: "New playlist dropping Friday. Lo-fi beats to study/relax to.",
		Clue: &types.Clue{
			Key:   "playlist id",
			Value: "spotify:playlist:3cI8rXJw7X8rK3mNpQ",
		},
	},
	{
		Id:      7,
		Author:  "Grace Travels",
		Content: "Japan trip was absolutely incredible. The food alone was worth it.",
		Clue:    nil,
	},
	{
		Id:      8,
		Author:  "Hank Fit",
		Content: "Morning run done. 5K in under 25 minutes. Personal best!",
		Clue:    nil,
	},
	{
		Id:      9,
		Author:  "Ivy Reads",
		Content: "Just started reading Dune. Already 200 pages in and hooked.",
		Clue:    nil,
	},
	{
		Id:      10,
		Author:  "Jack Cook",
		Content: "Tried making sourdough from scratch. Day 3 of the starter and it smells terrible.",
		Clue:    nil,
	},
}
