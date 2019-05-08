## some_keng
### 环境部署


### 前端
#### 1.js精度
#### throw
* js处理int64位数据存在精度损失，以mid=4621966154320056320为例，通过ajax请求，发至前端浏览器，处理时mid变成了4621966154320056000。

#### catch
* 目前解决的方案，将这些参数转成string型
* 由此考虑，所有需要在前端展示的数据结构，如MblogInfo，对应生成MblogInfoView结构，一方面将int64数据转换成string型，防止精度损失，另一方面屏蔽一些敏感信息，如User.PassWord

### 后端
#### 1.is_friend字段的更新问题
##### throw
* 每次接收到关注、取关消息的时候，relation服务需要处理该消息，修改follow、fan表，其中有个字段是is_friend，表征两者是否互关
* 由于消息的生产分区依据是根据发起关注人的uid进行取模hash,所以a关注b,b取关a产生的消息将由同一个消费组中的两个不同的消费者进行消费。因而可能产生时间先后与处理时间先后顺序相反，或者并发。
* 正确顺序（首先a先关注b的消息处理,然后b取关a的消息处理）：11.新增fan表记录；12.检测b是否也关注a，发现b也关注a;3.设置a和b的follow_info.is_friend=1。若在2-3之间b取关a,此时未进行额外检测，则is_friend=1此时即是错误数据。然后处理b取关a的消息：21.修改fan表，22.修改follow表的两条数据的is_friend=0。正确顺序似乎没问题。
* 交叉顺序（处理b取关a的消息快于处理a关注b的消息）：11.a关注b；12.开始消费a关注b的消息；13新增a关注b的fan_info表记录；14.检测b是否也关注a,结果为是，判定两者is_friend=1;21.b取关a;22.开始消费b取关a的消息；23.将两者的follow_info.is_friend=0；15.设置两者的follow_info.is_friend=1。所以出现问题了。。。
* 会不会出现死锁？？？
    ```sql
        update follow_info set is_friend= 1 where uid=a and follow_uid=b;
        update follow_info set is_friend= 1 where uid=b and follow_uid=a;
    ```
    ```sql
        update follow_info set is_friend= 0 where uid=b and follow_uid=a;
        update follow_info set is_friend= 0 where uid=a and follow_uid=b;
    ```
两个事务均执行了第一句，将该行数据锁定。进而导致死锁？？？

##### catch
* 一把梭，前端发起关注请求时，将所有数据的更新同步处理

### go
#### 1.go test 缓存
##### throw
* go test 指令进行单元测试存在缓存，若功能代码和测试代码没变动，则两次执行单元测试会采用缓存结果
##### catch
* 在执行测试时增加参数 -count=1,(etc: go test -v xx_test.go -count=1 -test.run TestXxFunc)
