package commentService

type CommentServicer interface {
	//点赞微博
	LikeMblog(uid int64, uidLiked int64, mid int64) error
	UnLikeMblog(uid int64, uidLiked int64, mid int64) error
	GetLikeCnt(mids []int64) map[int64]int32
	CheckLikesByUid(mids []int64, uid int64) map[int64]bool

	//转发
	TransMblog(uid int64, uidOri int64, midOri int64, content string) error
	GetTransCnt(mids []int64) map[int64]int32

	//评论微博
	AddComment(uid int64, uidMblog int64, midMblog int64, content string) error
	GetCommentCnt(mids []int64) map[int64]int32

	//点赞评论
	LikeComment(uid int64, uidLiked int64, mid int64, commentId int64) error
	UnLikeComment(uid int64, uidLiked int64, mid int64, commentId int64) error
	GetCommentLikeCnt(mid int64, commentIds []int64) map[int64]int32
	CheckCommentLikesByUid(mid int64, commentIds []int64, uid int64) map[int64]bool

	//评论评论
}
