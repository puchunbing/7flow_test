package main

import (
	"awesomeProject/log"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"math/rand"
	"net"
	"sort"
	"time"
)

// Horse 结构体表示一匹马
type Horse struct {
	Name       string
	FinishTime time.Duration
}

// 定义一个结构体
type Example struct {
	FieldA string
	FieldB string
	FieldC string
}

// 根据字段的A-Z排序生成MD5哈希值
func generateMD5(e Example) string {
	// 获取字段值
	fields := []string{e.FieldA, e.FieldB, e.FieldC}

	// 按字母顺序排序
	sort.Strings(fields)

	// 拼接字段值
	concatenated := ""
	for _, field := range fields {
		concatenated += field
	}

	// 生成MD5哈希值
	hash := md5.Sum([]byte(concatenated))
	return hex.EncodeToString(hash[:])
}

type RedisConfig struct {
	Host string
	Port int
	Pass string
}

var redisConfigs = map[string]RedisConfig{
	"redis10": {
		Host: "192.168.2.100",
		Port: 6379,
		Pass: "123456",
	},
}

func getRedisClient(key string) (*redis.Client, error) {
	config, exists := redisConfigs[key]
	if !exists {
		return nil, fmt.Errorf("redis config for key %s not found", key)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Pass,
		DB:       0, // use default DB
	})

	return rdb, nil
}

// 自定义排序规则：1-8, a-z, A-Z, 其他符号
func customSort(a, b string) bool {
	// 定义排序优先级
	priority := func(r rune) int {
		switch {
		case r >= '1' && r <= '8':
			return int(r) - '1'
		case r >= 'a' && r <= 'z':
			return int(r) - 'a' + 8
		case r >= 'A' && r <= 'Z':
			return int(r) - 'A' + 34
		default:
			// 符号优先级排在最后
			return int(r) + 100
		}
	}

	// 比较两个字符串
	for i := 0; i < len(a) && i < len(b); i++ {
		if priority(rune(a[i])) != priority(rune(b[i])) {
			return priority(rune(a[i])) < priority(rune(b[i]))
		}
	}

	// 如果当前字符相同，则比较长度
	return len(a) < len(b)
}

var RequestParamSecret = "GpBWe7P5nUXibWd9gq70"

// 生成唯一的雪花 ID 并返回字符串类型
func generateSnowflakeID() string {
	timestamp := time.Now().UnixNano()
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Int63()
	snowflakeID := (timestamp << 22) | (randomNumber & 0x3FFFFF)
	if snowflakeID < 0 {
		snowflakeID = -snowflakeID
	}
	return fmt.Sprintf("%d", snowflakeID)
}

func generateUniqueString() string {
	return uuid.New().String()
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("no local IP found")
}

type BettingReq struct {
	Uid           json.Number `json:"uid" binding:"required"`
	GameId        string      `json:"game_id"`
	RoundId       string      `json:"round_id"`
	TransactionId string      `json:"transaction_id"`
	Amount        int         `json:"amount"`
	AppKey        string      `json:"app_key"`
	Sign          string      `json:"sign"`
	Timestamp     int         `json:"timestamp"`
	Currency      string      `json:"currency"`
}
type ImageURLs struct {
	En string `json:"en"`
	Hi string `json:"hi"`
}

type Content struct {
	TitleImg   ImageURLs `json:"title_img"`
	ContentImg ImageURLs `json:"content_img"`
}

func main() {

	ip, _ := getLocalIP()
	fmt.Println(ip)
	jsonString := `{
        "title_img": {
            "en": "https://cocos-teenpatti.s3.ap-south-1.amazonaws.com/service/assets/images/invite_banner.png",
            "hi": "https://cocos-teenpatti.s3.ap-south-1.amazonaws.com/service/assets/images/invite_banner.png"
        },
        "content_img": {
            "en": "https://cocos-teenpatti.s3.ap-south-1.amazonaws.com/service/assets/images/invite_content.png",
            "hi": "https://cocos-teenpatti.s3.ap-south-1.amazonaws.com/service/assets/images/invite_content.png"
        }
    }`

	var content Content
	err := json.Unmarshal([]byte(jsonString), &content)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	va := `
{
    "game_id": "100",
   "currency": "INR",
    "app_key": "FD9B43A96B93B529D9E2DC68BBE1B69B",
    "timestamp": 100,
    "sign": "123456789f9191460984a38486be44696ef7b5d512345678"
}
`

	var req map[string]any
	// 解析请求数据
	err = json.Unmarshal([]byte(va), &req)
	if err != nil {
		return
	}

	req["app_secret"] = "D24174096B5CDC4350AB2844CBDA53BD"
	// 对请求参数的键进行排序 (1-8, a-z, A-Z, 符号)，排除 sign 字段
	keys := make([]string, 0, len(req))
	for key := range req {
		if key == "app_key" || key == "timestamp" {
			keys = append(keys, key)
		}
	}

	hasAppKey := false
	hasTimestamp := false

	for _, k := range keys {
		if k == "app_key" {
			hasAppKey = true
		}
		if k == "timestamp" {
			hasTimestamp = true
		}
	}

	if !hasAppKey || !hasTimestamp {

		return
	}

	if _, ok := req["timestamp"].(float64); !ok {

		return
	}

	signatureString := fmt.Sprintf("%v%v%v", req["app_key"], req["app_secret"], int64(req["timestamp"].(float64)))

	// 移除最后一个多余的 "&"
	//signatureString := strings.TrimSuffix(builder.String()+global.RequestParamSecret, "&")

	// 生成 MD5 签名
	md5Hash := md5.New()
	md5Hash.Write([]byte(signatureString))
	signature := hex.EncodeToString(md5Hash.Sum(nil))
	newSign := "12345678" + signature + "12345678"
	fmt.Println(newSign)
	log.Logger.Infof("签名:%v", newSign)

	//ctx := context.Background()
	//rdb, err := getRedisClient("redis10")
	//if err != nil {
	//	fmt.Println("Failed to get Redis client:", err)
	//	return
	//}
	//
	//// 使用 Redis 客户端进行操作
	//ukey := "hu:3166450"
	//val, err := rdb.Get(ctx, ukey).Result()
	//if err != nil {
	//	fmt.Println("Failed to get value from Redis:", err)
	//	return
	//}
	//
	//fmt.Println("Value:", val)
	//
	//e := Example{
	//	FieldA: "valueA",
	//	FieldB: "valueB",
	//	FieldC: "valueC",
	//}
	//
	//md5Hash := generateMD5(e)
	//fmt.Println("MD5哈希值:", md5Hash)
	//
	//// 存储原始字段值和哈希值
	//storedData := map[string]Example{
	//	md5Hash: e,
	//}
	//
	//// 解析出原始字段值
	//if original, exists := storedData[md5Hash]; exists {
	//	fmt.Println("原始字段值:", original)
	//} else {
	//	fmt.Println("未找到原始字段值")
	//}
	//
	//rand.Seed(time.Now().UnixNano())
	//// 定义五匹马
	//ch := make(chan []Horse)
	//horses := []Horse{
	//	{Name: "马1"},
	//	{Name: "马2"},
	//	{Name: "马3"},
	//	{Name: "马4"},
	//	{Name: "马5"},
	//}
	//
	//go func() {
	//	ch <- horses
	//}()
	//read := <-ch
	//close(ch)
	//fmt.Println(read)
	//
	//var wg sync.WaitGroup
	//results := make(chan Horse, len(horses))
	//
	//// 启动五个goroutine来模拟比赛
	//for i := range horses {
	//	wg.Add(1)
	//	go func(h *Horse) {
	//		defer wg.Done()
	//		// 随机生成到达终点的时间
	//		h.FinishTime = time.Duration(rand.Intn(10)+1) * time.Second
	//		time.Sleep(h.FinishTime)
	//		results <- *h
	//	}(&horses[i])
	//}
	//// 等待所有马到达终点
	//go func() {
	//	wg.Wait()
	//	close(results)
	//}()
	//
	//// 收集比赛结果
	//var raceResults []Horse
	//for result := range results {
	//	raceResults = append(raceResults, result)
	//}
	//
	//// 按到达时间排序
	//sort.Slice(raceResults, func(i, j int) bool {
	//	return raceResults[i].FinishTime < raceResults[j].FinishTime
	//})
	//
	//// 打印比赛结果
	//fmt.Println("比赛结果：")
	//for i, horse := range raceResults {
	//	fmt.Printf("第%d名: %s, 时间: %v\n", i+1, horse.Name, horse.FinishTime)
	//}
}
