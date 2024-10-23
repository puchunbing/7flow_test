# 神牛重构

## 功能需求
    - 展示活动倒计时
    - 展示当前累计投注：XX/1000
    - 展示当前获得祈福券数量：拥有数量/获得数量
    - 展示排行榜
    - 展示个人获奖记录
    - 埋点统计
bet转fragments，中奖金额bonus
## 要点
- 数据均从redis：rank中获取
   --除排行榜从tmp中获取
   --bonus从data中获取
- 押注金额(Betting)：
    - 1000为配置金额 
    - 每押注1000金额，则获得+1张祈福券(碎片)
    - 超过1000金额，清空剩余继续展示(如：已有1100，则展示100/1000)  重点解决！！！
    - redis中存储的数据为：act_10:user_all_bet20240814888
- 用户祈福券，碎片(Fragments)：
    - 展示为：剩余数量/当天获取数量
    - 每天清0
    - 剩余数量redis中存储的数据为：act_10:user_user_fragments20240814888
    - 当天获取数量redis中存储的数据为：act_10:user_fragments_all20240814888
- 抽奖  单位：bonus
    - 紫牛加倍逻辑：类似于回扣，第一次返利100%，第二次返利50%，第三次及三次以后返利10%，每天清0
    - 
    - 若用户抽中5000以上的奖励，则全服滚动广播
    - 抽奖奖励直接写入账户 无数据得写入
      - 奖励发放：
          下方获奖处理
             flow_id = 27
- 
- 排行榜，机器人均为自动化处理
   -奖励发放
    - 当时间为印度当地0或者24小时时，   发放前一天的奖励，  展示最近24小时内积分的前50名玩家
      1.初始化
        - cron 方法的主要功能是刷新活动排行榜数据。具体步骤如下：  
      设置日期参数：  
      如果传入了 ymd 和 today 参数，则设置相应的类属性，并标记 isShow 为 true。
      检查活动是否开放：  
      调用 show 方法检查活动是否开放，如果未开放且 isShow 为 false，则返回 true。
      从 Redis 获取用户积分数据：  
      使用 scan 方法遍历 Redis 中的键，获取用户积分数据，并存储在 $aTmp 和 $aSort 数组中。
      排序用户积分数据：  
      使用 array_multisort 对用户积分数据进行排序，并取前50名用户。
      从数据库获取用户和机器人排行榜数据：  
      查询 act_10_rank 表中当天的用户和机器人排行榜数据。
      合并排序排行榜数据：  
      将用户排行榜数据和机器人排行榜数据进行合并排序。
      更新排行榜数据到数据库：
      根据奖励配置更新排行榜数据到 act_10_rank 表中。
        2.排行榜数据更新
          获取邮件配置和排行榜数据：  
          从 aConfig 中获取邮件配置。
          从 act_10_rank 表中获取当天的排行榜数据。
          处理排行榜数据：  
          遍历排行榜数据，向符合条件的用户发送奖励邮件，并更新用户的奖励状态。
          统计用户数据：
          如果有用户数据，统计用户的充值、提现和玩牌消耗数据，并更新到 act_10_rank 表中
- 排名奖励写入act_10_rank表中reward字段


- 机器人处理逻辑   robotRank
刷新机器人数据
获取当前小时的机器人积分升级配置和奖励配置。
如果没有找到机器人积分升级配置，记录日志并返回。
检查是否显示数据，如果不显示，记录日志并返回。
从数据库中获取当天的机器人排名数据。
如果机器人数据不足50条，重新初始化机器人数据并返回。
遍历机器人数据，更新每个机器人的积分。
对更新后的机器人数据进行排序。
更新数据库中的机器人排名数据。


- 初始化机器人 initRobot
检查当天是否已有机器人数据，如果有则删除。
初始化机器人数据列表和排序数组。
根据配置生成50个机器人数据。
对生成的机器人数据进行排序。
将排序后的机器人数据插入数据库。


- 积分数据处理 piecesWith
遍历配置数据，找到符合条件的积分范围。
随机生成积分值并返回。



global.UserIntegral  用户积分每一千一积分

act10_rank 

reward  奖励
amount  积分


- 活动埋点
- `buriPoint` 函数的作用是记录用户在活动中的参与情况和金币发放情况。具体步骤如下：

1. 获取用户的登录渠道、渠道ID和创建时间。
2. 根据用户的登录渠道或渠道ID，确定 `channel_id`。
3. 获取当前日期，格式为 `Ymd`。
4. 如果 `type` 为 2 且 `amt` 大于 0，表示金币发放：
    - 生成一个键 `key`，格式为 `act_10:stat{today}:{channel_id}:3`。
    - 调用 `Util::BuriPointStat()->limit($key, $amt)` 记录金币发放。
5. 记录参与用户：
    - 生成一个键 `key`，格式为 `act_10:user{today}:{channel_id}:1:{uid}`。
    - 生成一个子键 `subKey`，格式为 `act_10:stat{today}:{channel_id}:1`。
    - 调用 `BuriPointStatModel::class->deduplication($key, $subKey)` 记录去重的参与用户。
6. 如果用户的创建时间是今天，记录新增参与用户：
    - 生成一个键 `key`，格式为 `act_10:user{today}:{channel_id}:2:{uid}`。
    - 生成一个子键 `subKey`，格式为 `act_10:stat{today}:{channel_id}:2`。
    - 调用 `Util::BuriPointStat()->deduplication($key, $subKey)` 记录去重的新增参与用户。
    - 如果 `type` 为 2 且 `amt` 大于 0，表示金币发放：
        - 生成一个键 `key`，格式为 `act_10:stat{today}:{channel_id}:4`。
        - 调用 `Util::BuriPointStat()->limit($key, $amt)` 记录金币发放。
7. 返回 `true` 表示操作成功。


# 更改点
 - sql拼接有问题








```sql
定义一个 cron 函数，处理定时刷新排行统计。
初始化两个切片 aTmp 和 aSort。
使用 Redis 的 Scan 方法遍历 Redis 键，获取用户积分数据并存储在 aTmp 和 aSort 中。
如果 aSort 和 aTmp 的长度不相等，记录日志并返回 false。
使用 sort.SliceStable 对 aTmp 进行排序，并将前 50 名用户存储在 aUserRank 中。
从数据库中获取用户和机器人排行榜数据，存储在 aUserLists 和 aRobotLists 中。
根据不同条件合并并排序 aUserRank、aUserLists 和 aRobotLists，存储在 aRank 中。
遍历 aRank，根据奖励配置更新数据库中的排行榜数据。
返回 true。
```







 
