package privacyutils

import "cyber-security-minigames/types"

func PostDataToUI(posts []types.Post) types.PrivacyUI {
	postUIs := make([]types.PostUI, len(posts))
	clueUIs := []types.ClueUI{}
	for i, p := range posts {
		if p.Clue != nil {
			postUIs[i] = types.PostUI{
				Id:            p.Id,
				Author:        p.Author,
				Content:       p.Content,
				AddedToFolder: p.Clue.AddedToFolder,
			}
			if p.Clue.AddedToFolder == true {
				clueUIs = append(clueUIs, types.ClueUI{
					p.Clue.Key,
					p.Clue.Value,
				})
			}
		} else {
			postUIs[i] = types.PostUI{
				Id:            p.Id,
				Author:        p.Author,
				Content:       p.Content,
				AddedToFolder: false,
			}
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
