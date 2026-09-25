package privacyutils

import "cyber-security-minigames/types"

// PostToUI maps a post to what the feed renders. A post only counts as added
// to the folder if it has a clue to add.
func PostToUI(p types.Post) types.PostUI {
	postUI := types.PostUI{
		Id:       p.Id,
		Author:   p.Author,
		Handle:   p.Handle,
		Posted:   p.Posted,
		Context:  p.Context,
		Content:  p.Content,
		Likes:    p.Likes,
		Comments: p.Comments,
		Shares:   p.Shares,
	}
	if p.Clue != nil {
		postUI.AddedToFolder = p.Clue.AddedToFolder
	}
	return postUI
}

func PostDataToUI(posts []types.Post) types.PrivacyUI {
	postUIs := make([]types.PostUI, len(posts))
	clueUIs := []types.ClueUI{}
	for i, p := range posts {
		postUIs[i] = PostToUI(p)
		if p.Clue != nil && p.Clue.AddedToFolder {
			clueUIs = append(clueUIs, types.ClueUI{
				Key:   p.Clue.Key,
				Value: p.Clue.Value,
			})
		}
	}
	return types.PrivacyUI{
		Posts: postUIs,
		Clues: clueUIs,
	}
}

func NumberOfClues(posts []types.Post) int {
	totalClues := 0
	for _, post := range posts {
		if post.Clue != nil {
			totalClues = totalClues + 1
		}
	}
	return totalClues

}
