package chatService

func (s *chatService) initRoomResource(uid int64) (err error) {

	////开启指定room的kafka队列（这一步理论上不需要）
	//
	////获取room中的配置信息
	//roomConfig, found := s.dbRepo.GetRoomConfigByUid(uid)
	//if !found {
	//	err = errors.New("not found the room about" + strconv.FormatInt(uid, 64))
	//	logger.Err(logType, err.Error())
	//	return err
	//}
	//
	////根据room的配置，启动redisSetCnt个consumer group消费队列中的数据
	//redisSetCnt := roomConfig.RedisSetCnt
	//for i := 0; i < redisSetCnt; i++ {
	//	go func() {
	//
	//	}()
	//}

	return nil
}

func (s *chatService) releaseRoomResource(uid int64) (err error) {
	return nil
}
