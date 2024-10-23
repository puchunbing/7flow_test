package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var TaConf = TConf{}

type TConf struct {
	Taurus []TaurusConf `mapstructure:"taurus" json:"taurus"` // 金牛配置
}

type TaurusConf struct {
	Consume   ConsumeConf `mapstructure:"consume" json:"consume"`       // 消耗配置
	WhiteList []int64     `mapstructure:"white_list" json:"white_list"` // 白名单
}

type ConsumeConf struct {
	ActivityStart string        `mapstructure:"activity_start" json:"activity_start"` // 活动开始时间
	ActivityEnd   string        `mapstructure:"activity_end" json:"activity_end"`     // 活动结束时间
	Betting       float64       `mapstructure:"betting" json:"betting"`               // 下注数
	ConsumeConf   []ConsumeItem `mapstructure:"consume_conf" json:"consume_conf"`     // 消耗配置项
	RewardConfigs []RewardConf  `mapstructure:"reward_configs" json:"reward_configs"` // 奖励配置
}

type ConsumeItem struct {
	BullType   int `mapstructure:"bull_type" json:"bull_type"`     // 牛类型  1-绿牛，2-蓝牛，3-紫牛
	MaxReward  int `mapstructure:"max_reward" json:"max_reward"`   // 最大奖励
	RewardType int `mapstructure:"reward_type" json:"reward_type"` // 奖励类型  1-bmoney, 2-money, 3-bonus
	Debris     int `mapstructure:"debris" json:"debris"`           // 碎片
}

type RewardConf struct {
	BullType   int             `mapstructure:"bull_type" json:"bull_type"`   // 牛类型
	Count      int             `mapstructure:"count" json:"count"`           // 碎片消耗数量
	Reward     []Reward        `mapstructure:"reward" json:"reward"`         // 奖励
	Conditions RewardCondition `mapstructure:"conditions" json:"conditions"` // 返利比例
}

type Reward struct {
	Amount      float64 `mapstructure:"amount" json:"amount"`           // 奖励金额
	Probability float64 `mapstructure:"probability" json:"probability"` // 概率
}

type RewardCondition struct {
	First     float64 `mapstructure:"first" json:"first"`           // 第一次返利比例
	Second    float64 `mapstructure:"second" json:"second"`         // 第二次返利比例
	ThirdPlus float64 `mapstructure:"third_plus" json:"third_plus"` // 第三次及以后返利比例
}

func InitTaurus() {
	viper.SetConfigFile("./config/taurus.json")
	err := viper.ReadInConfig()
	if err != nil {
		return
	}

	if err = viper.Unmarshal(&TaConf); err != nil {
		return
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err = viper.Unmarshal(&TaConf); err != nil {
			return
		}
	})
	return
}

/**
// TaurusLottery 金牛抽奖自动任务
func TaurusLottery() {
 h, m := tools.NowHourMinute()
 fmt.Println("执行金牛抽奖定时任务:", h, m)
 currentHour := time.Now().Hour()
 c := cron_.New()

 // 添加一个在每天0:00执行的任务
 c.AddFunc("0 0 0 * * *", func() {
  if !cron.ActivityEnd() {
   fmt.Println("在0:00执行任务: ", time.Now().Format("2006-01-02 15:04:05"))
   if m == 0 && h == 0 {
    yesterUnix := (tools.TODAY() - 86400)
    issueReward(tools.YesterdayStr(), int(yesterUnix), true) //发放奖励
   }
  }
 })

 // 添加一个在每天23:59执行的任务
 c.AddFunc("59 23 * * *", func() {  已验证
  if !cron.ActivityEnd() {
   fmt.Println("在23:59执行任务: ", time.Now().Format("2006-01-02 15:04:05"))
   if m == 59 && h == 23 {
    Cron(tools.UinxToStrDay(tools.Now()), int(tools.TODAY()), false) //刷新排行统计
   }
  }
 })

 // 添加在指定小时执行的任务
 c.AddFunc("0 2,4,6,8,10,12,14,16,18,20,22 * * *", func() {
  if !cron.ActivityEnd() {
   fmt.Println("在指定小时执行任务: ", time.Now().Format("2006-01-02 15:04:05"))
   if m == 0 && cron.Contains([]int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22}, h) {
    robotRank(h, int(tools.TODAY()))                                 //刷新机器人数据
    Cron(tools.UinxToStrDay(tools.Now()), int(tools.TODAY()), false) //刷新排行统计
   }
  }
 })

 c.Start()
 select {}
}
*/

/**  已验证

   	c.AddFunc("@every 1m", func() {
		fmt.Println("时间到了: ", time.Now().Format("2006-01-02 15:04:05"), "tell")
	})

	c.AddFunc("@every 3m", func() {
		fmt.Println("每3分执行一次", time.Now().Format("2006-01-02 15:04:05"))
	})
	c.Start()
	select {}


*/
