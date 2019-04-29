package feedService

func (f *feedService) appendNewMid(uid int64, groupId int64, mid int64) (err error) {
	return f.feedRdRepo.AppendNewMid(uid, groupId, mid)
}

func (f *feedService) removeMultiMids(uid int64, groupId int64, mids []int64) (err error) {
	return f.feedRdRepo.RemoveMids(uid, groupId, mids)
}
