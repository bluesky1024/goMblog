## counter_server_scheme(sub of mblog_comment server)
### 功能
* 提供各种计数服务(包括点赞、评论数、转发数)
* 根据计数统计进行排序，返回各类按序列表
* 点赞统计比较特殊，似乎需要保存点赞人的uid，也就是需要快速判断某人是否对某条微博或者评论进行了点赞

### 方案
似乎无非就是各种缓存。。。没什么好思路