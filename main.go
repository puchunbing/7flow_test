package main

import (
	"7flow_test/config"
	"7flow_test/log"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func init() {
	config.InitTaurus()
}

type Reward struct {
	Rank      int    `json:"rank,omitempty"`
	RankRange string `json:"rank_range,omitempty"`
	Reward    int    `json:"reward"`
}

type RankingRewards struct {
	RankingRewards []Reward `json:"ranking_rewards"`
}

func getRewardByRank(rank int, jsonData string) int {
	var data RankingRewards
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		fmt.Println("解析 JSON 时出错:", err)
		return 0
	}

	for _, reward := range data.RankingRewards {
		if reward.Rank == rank {
			return reward.Reward
		}
		if reward.RankRange != "" {
			ranges := strings.Split(reward.RankRange, "-")
			min, _ := strconv.Atoi(ranges[0])
			max, _ := strconv.Atoi(ranges[1])
			if rank >= min && rank <= max {
				return reward.Reward
			}
		}
	}

	return 0 // 如果未找到排名，返回默认奖励
}

func scheduledTask() {
	// 定时执行的任务
	fmt.Println("执行定时任务:", time.Now())
}

type aaa struct {
	RobotConfig map[int][]Rof `mapstructure:"robot_config" json:"robot_config"` // 机器人配置
}

type Rof struct {
	Fen FnConfig `mapstructure:"fen" json:"fen"`
	No  FnConfig `mapstructure:"no" json:"no"`
}

type FnConfig struct {
	Max int `mapstructure:"max" json:"max"`
	Min int `mapstructure:"min" json:"min"`
}

type Rof1 struct {
	No  Range
	Fen Range
}

type Range struct {
	Min int
	Max int
}

func PiecesWith(piecesConfig []Rof1, no int) float64 {
	amount := 0.0
	for _, val := range piecesConfig {
		if no >= val.No.Min && no <= val.No.Max {
			rand.Seed(time.Now().UnixNano())
			amount = float64(rand.Intn(val.Fen.Max-val.Fen.Min+1) + val.Fen.Min)
			break
		}
	}
	return amount
}

func main() {

	log.Logger.Info("Hello, World!")
	log.Logger.Debug("Hello, World!")
	five := float64(5)
	u := int64(1000)

	y := float64(five) / 100
	v := u - (u * int64(five) / 100)
	fmt.Println(y, v)

	//redesTest()
	//piecesConfig := []Rof1{
	//	{No: Range{Min: 1, Max: 5}, Fen: Range{Min: 10, Max: 20}},
	//	{No: Range{Min: 6, Max: 10}, Fen: Range{Min: 15, Max: 25}},
	//}
	//no := 7
	//amount := PiecesWith(piecesConfig, no)
	//println("Amount:", amount)
	//
	//t := time.Now()
	//fmt.Println(t)
	//ce := 1
	//ce2 := []int{1, 2, 3, 4, 5}
	//
	//if !slices.Contains(ce2, ce) {
	//	fmt.Println("no")
	//}
	lucky := `

{
  "lucky_list": [
    {
      "rank": 1,
      "name": "player75957",
      "total": 384,
      "cash": 10000,
      "uid": 4682,
      "avatar": "head_9.png",
      "avatar_url": ""
    },
    {
      "rank": 2,
      "name": "Arman Saif",
      "total": 351,
      "cash": 6000,
      "uid": 2872,
      "avatar": "image_2872.png",
      "avatar_url": "http://load.ala11.in/india_avatar/image_2872.png"
    },
    {
      "rank": 49,
      "name": "player28085",
      "total": 190,
      "cash": 500,
      "uid": 9964,
      "avatar": "head_2.png",
      "avatar_url": ""
    },
    {
      "rank": 50,
      "name": "Bhim Kumar",
      "total": 180,
      "cash": 500,
      "uid": 7888,
      "avatar": "image_6738.png",
      "avatar_url": "http://load.ala11.in/india_avatar/image_6738.png"
    }
  ]
}
`
	me := `

  {
      "rank": 0,
      "name": "shit",
      "total": 500,
      "cash": 0,
      "uid": 7888,
      "avatar": "8888.jpg",
      "avatar_url": "http://hao123.com"
    }



`
	type UserList struct {
		Rank      any     `json:"rank"`  // 排名
		Name      string  `json:"name"`  // 用户名
		Total     float64 `json:"total"` // 总碎片
		Cash      float64 `json:"cash"`  // 奖励金额
		Uid       int64   `json:"uid"`   // 用户id
		Avatar    string  `json:"avatar"`
		AvatarUrl string  `json:"avatar_url"`
	}

	type luck struct {
		LuckyList []UserList `json:"lucky_list"` // 最近24小时前50名的获奖玩家
	}
	var data luck
	err1 := json.Unmarshal([]byte(lucky), &data)
	if err1 != nil {
		fmt.Println("解析 JSON 时出错:", err1)
	}
	fmt.Println(data)

	var meData UserList

	err2 := json.Unmarshal([]byte(me), &meData)
	if err2 != nil {
		fmt.Println("解析 me JSON 时出错:", err2)
		return
	}
	fmt.Println(data.LuckyList[len(data.LuckyList)-1].Total)
	if data.LuckyList[len(data.LuckyList)-1].Total > meData.Total {
		fmt.Println("干不过")
	} else {

		for i, luckyUser := range data.LuckyList {
			if luckyUser.Uid == meData.Uid && luckyUser.Total != meData.Total {
				if luckyUser.Total <= meData.Total { //更新
					meData.Rank = luckyUser.Rank
					meData.Cash = luckyUser.Cash
					data.LuckyList[i] = meData
				}

			} else {
				meData.Rank = luckyUser.Rank
				meData.Cash = luckyUser.Cash
				data.LuckyList[i] = meData
				break
			}
		}

	}

	fmt.Println(data)

	a := time.Now()
	fmt.Println(a, meData)

	//redesTest()

	var dd aaa
	ddd := `{  "robot_config": {
        "2": [
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "4": [
          {
            "fen": {
              "max": 170,
              "min": 150
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 0,
              "min": 0
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "6": [
          {
            "fen": {
              "max": 220,
              "min": 190
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 190,
              "min": 170
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 170,
              "min": 150
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 150,
              "min": 100
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 120,
              "min": 70
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "8": [
          {
            "fen": {
              "max": 280,
              "min": 250
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 250,
              "min": 220
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 220,
              "min": 190
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 190,
              "min": 160
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 160,
              "min": 140
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "10": [
          {
            "fen": {
              "max": 340,
              "min": 300
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 300,
              "min": 270
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 270,
              "min": 220
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 220,
              "min": 170
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 160,
              "min": 150
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "12": [
          {
            "fen": {
              "max": 420,
              "min": 370
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 370,
              "min": 320
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 320,
              "min": 270
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 270,
              "min": 220
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 220,
              "min": 190
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "14": [
          {
            "fen": {
              "max": 440,
              "min": 400
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 400,
              "min": 370
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 370,
              "min": 300
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 300,
              "min": 240
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 240,
              "min": 200
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "16": [
          {
            "fen": {
              "max": 500,
              "min": 470
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 470,
              "min": 420
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 420,
              "min": 320
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 320,
              "min": 280
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 270,
              "min": 210
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "18": [
          {
            "fen": {
              "max": 540,
              "min": 500
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 500,
              "min": 470
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 470,
              "min": 340
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 340,
              "min": 300
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 280,
              "min": 240
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "20": [
          {
            "fen": {
              "max": 570,
              "min": 520
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 520,
              "min": 500
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 500,
              "min": 420
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 420,
              "min": 340
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 320,
              "min": 270
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ],
        "22": [
          {
            "fen": {
              "max": 620,
              "min": 580
            },
            "no": {
              "max": 1,
              "min": 1
            }
          },
          {
            "fen": {
              "max": 580,
              "min": 520
            },
            "no": {
              "max": 2,
              "min": 2
            }
          },
          {
            "fen": {
              "max": 520,
              "min": 470
            },
            "no": {
              "max": 3,
              "min": 3
            }
          },
          {
            "fen": {
              "max": 470,
              "min": 400
            },
            "no": {
              "max": 10,
              "min": 4
            }
          },
          {
            "fen": {
              "max": 370,
              "min": 320
            },
            "no": {
              "max": 50,
              "min": 11
            }
          }
        ]
      }
}

`

	err := json.Unmarshal([]byte(ddd), &dd)
	fmt.Println(dd)

	robotConfig := map[int][]map[string]map[string]int{
		2: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 0, "max": 0}},
		},
		4: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 150, "max": 170}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 0, "max": 0}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 0, "max": 0}},
		},
		6: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 190, "max": 220}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 170, "max": 190}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 150, "max": 170}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 100, "max": 150}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 70, "max": 120}},
		},
		8: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 250, "max": 280}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 220, "max": 250}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 190, "max": 220}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 160, "max": 190}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 140, "max": 160}},
		},
		10: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 300, "max": 340}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 270, "max": 300}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 220, "max": 270}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 170, "max": 220}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 150, "max": 160}},
		},
		12: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 370, "max": 420}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 320, "max": 370}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 270, "max": 320}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 220, "max": 270}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 190, "max": 220}},
		},
		14: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 400, "max": 440}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 370, "max": 400}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 300, "max": 370}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 240, "max": 300}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 200, "max": 240}},
		},
		16: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 470, "max": 500}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 420, "max": 470}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 320, "max": 420}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 280, "max": 320}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 210, "max": 270}},
		},
		18: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 500, "max": 540}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 470, "max": 500}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 340, "max": 470}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 300, "max": 340}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 240, "max": 280}},
		},
		20: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 520, "max": 570}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 500, "max": 520}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 420, "max": 500}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 340, "max": 420}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 270, "max": 320}},
		},
		22: {
			{"no": {"min": 1, "max": 1}, "fen": {"min": 580, "max": 620}},
			{"no": {"min": 2, "max": 2}, "fen": {"min": 520, "max": 580}},
			{"no": {"min": 3, "max": 3}, "fen": {"min": 470, "max": 520}},
			{"no": {"min": 4, "max": 10}, "fen": {"min": 400, "max": 470}},
			{"no": {"min": 11, "max": 50}, "fen": {"min": 320, "max": 370}},
		},
	}

	jsonData1, err := json.MarshalIndent(robotConfig, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	fmt.Println(string(jsonData1))

	//// 创建一个新的cron调度器
	//c := cron.New()
	//
	//// 添加定时任务，使用自定义时间表达式
	//// 添加定时任务，每10分钟执行一次
	//_, err := c.AddFunc("@every 10m", scheduledTask)
	//if err != nil {
	//	fmt.Println("添加定时任务时出错:", err)
	//	return
	//}
	//
	//// 启动调度器
	//c.Start()
	//select {}
	//redesTest()

	jsonData := `{
		"ranking_rewards": [
			{
				"rank": 1,
				"reward": 10000
			},
			{
				"rank": 2,
				"reward": 6000
			},
			{
				"rank": 3,
				"reward": 3000
			},
			{
				"rank_range": "4-10",
				"reward": 1000
			},
			{
				"rank_range": "11-50",
				"reward": 500
			}
		]
	}`

	rank := 2
	reward := getRewardByRank(rank, jsonData)
	fmt.Printf("排名: %d, 奖励: %d\n", rank, reward) // 输出: 排名: 5, 奖励: 1000

	Reard()

	CONF := int((config.TaConf.Taurus[0].Consume.RewardConfigs[2].Reward[5].Probability) * 10000)
	fmt.Println(CONF)

	//redesTest()

	fmt.Println("Hello, World!")
	startTimeStr := "2024-08-14 00:00:00"
	endTimeStr := "2024-08-15 23:59:59"
	remainingSeconds := RemainingTime(startTimeStr, endTimeStr)
	A := time.Now()
	fmt.Println(A)
	remainingHours := remainingSeconds / 3600
	fmt.Printf("Remaining time: %d hours\n", remainingHours)

}

func Reard() {
	// 定义奖品和它们的权重
	weights := []float64{0.004, 0.006, 0.1, 0.2, 0.4, 0.2936}
	prizes := []float64{30000, 10000, 5000, 2000, 1000, 300} // 对应的中奖金额

	// 计算总权重
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight
	}

	// 计算每个奖品的随机数区间
	var sum float64
	for i, weight := range weights {
		startRange := sum
		sum += weight
		endRange := sum
		fmt.Printf("奖品 %d 对应的奖金是 %.2f，随机数区间是 [%.4f, %.4f)\n", i+1, prizes[i], startRange, endRange)
	}

	// 生成0到总权重之间的随机数
	randNum := rand.Float64() * totalWeight
	// 累加权重，直到累加值大于随机数
	sum = 0.0
	for i, weight := range weights {
		sum += weight
		if sum >= randNum {
			fmt.Printf("恭喜你，抽中了奖品 %d，摇出的数字是 %.4f，对应的奖金是 %.2f\n", i+1, randNum, prizes[i])
			break
		}
	}
}

func RemainingTime(startTimeStr, endTimeStr string) int64 {
	const layout = "2006-01-02 15:04:05"
	startTime, err := time.Parse(layout, startTimeStr)
	if err != nil {
		return 0
	}
	endTime, err := time.Parse(layout, endTimeStr)
	if err != nil {
		return 0
	}
	currentTime := time.Now()
	if currentTime.Before(startTime) || currentTime.After(endTime) {
		return 0
	}
	remainingTime := endTime.Sub(currentTime).Seconds()
	return int64(remainingTime)
}

func redesTest() {

	// Create a new Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.2.100:6379", // Redis server address
		Password: "123456",             // Redis Password
	})

	// Define user ID and fragment amount
	uid := 333
	amount := 888
	ymd := time.Now().Format("20060102") // Current date in Ymd format

	// Create the Redis key
	fraKey := fmt.Sprintf("act_10:user_fragments:%s:%d", ymd, uid)
	fraKey = "act_10:user_fragments:20240910:3166450"
	//// 查询键的过期时间
	//ttl := rdb.TTL(context.Background(), fraKey).Val()
	//fmt.Printf("Key expiration time: %v\n", ttl)

	value, err := rdb.Get(context.Background(), fraKey).Float64()
	if err != nil {
		fmt.Println("Error getting value from Redis:", err)
		return
	}
	fmt.Println(value)

	// Set the user fragments data
	ctx := context.Background()
	//err := rdb.Set(ctx, fraKey, amount, 0).Err()

	err = rdb.IncrByFloat(ctx, fraKey, float64(amount)).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Expire(ctx, fraKey, 48*time.Hour).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("User fragments data created successfully.")
}

func mysqlTest() {
	// 连接到 MySQL 数据库
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.2.100:3306)/yd_game")
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer db.Close()

	// 定义插入数据的 SQL 语句
	sql := "INSERT INTO act_10_rank (uid, type, ctime, no, amount, reward) VALUES (12345, 1, UNIX_TIMESTAMP(CURDATE()), 1, 100, 50)"

	// 执行 SQL 语句插入一条测试数据
	_, err = db.Exec(sql)
	if err != nil {
		fmt.Println("插入数据时出错:", err)
		return
	}

	fmt.Println("测试数据插入成功")
}

//docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3307:3306 -d mysql
//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
