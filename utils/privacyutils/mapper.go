package privacyutils

import "cyber-security-minigames/types"

// PostToUI maps a post to the data the feed and analysis modal need to render.
func PostToUI(p types.Post) types.PostUI {
	postUI := types.PostUI{
		Id:               p.Id,
		Author:           p.Author,
		Handle:           p.Handle,
		Posted:           p.Posted,
		Context:          p.Context,
		Content:          p.Content,
		Likes:            p.Likes,
		Comments:         p.Comments,
		Shares:           p.Shares,
		Analysed:         p.Analysed,
		HasClue:          p.Clue != nil,
		AnalysisOptions:  append([]string(nil), p.AnalysisOptions...),
		CorrectOption:    p.CorrectOption,
		AnalysisFeedback: p.AnalysisFeedback,
	}
	if p.Clue != nil {
		postUI.AddedToFolder = p.Clue.AddedToFolder
	}
	return postUI
}

func PostDataToUI(posts []types.Post) types.PrivacyUI {
	postUIs := make([]types.PostUI, len(posts))
	clueUIs := []types.ClueUI{}
	analysedCount := 0

	for i, p := range posts {
		postUIs[i] = PostToUI(p)
		if p.Analysed {
			analysedCount++
		}
		if p.Clue != nil && p.Clue.AddedToFolder {
			clueUIs = append(clueUIs, types.ClueUI{
				Key:   p.Clue.Key,
				Value: p.Clue.Value,
			})
		}
	}

	return types.PrivacyUI{
		Posts:         postUIs,
		Clues:         clueUIs,
		AnalysedCount: analysedCount,
		TotalPosts:    len(posts),
	}
}
