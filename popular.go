package main

type postByPopular []postAPI

func (p postByPopular) Len() int      { return len(p) }
func (p postByPopular) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p postByPopular) Less(i, j int) bool {
	return p[i].LikeCount+p[i].DislikeCount+p[i].CommentCount > p[j].LikeCount+p[j].DislikeCount+p[j].CommentCount
}
