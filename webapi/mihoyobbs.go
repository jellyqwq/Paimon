package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jellyqwq/Paimon/requests"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongodb 数据库权限管理
// https://www.cnblogs.com/harrychinese/p/mongodb.html

type BBSGoodsType struct {
	GoodsID   string `json:"goods_id" bson:"goods_id"`
	GoodsName string `json:"goods_name" bson:"goods_name"`
	Type      int    `json:"type" bson:"type"`
	Price     int    `json:"price" bson:"price"`
	Icon      string `json:"icon" bson:"icon"`
	NextTime  int    `json:"next_time" bson:"next_time"`
	NowTime   int    `json:"now_time" bson:"now_time"`
}

func MihoyoBBSGoodsUpdate() error {
	clientOptions := options.Client().ApplyURI("mongodb://Nahida:Nahida1027@localhost:27017/Paimon")
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	MongoDBClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error(err)
		return err
	}
	err = MongoDBClient.Ping(context.TODO(), nil)
	// 测试连接
	if err != nil {
		log.Error(err)
		return err
	}
	collection := MongoDBClient.Database("Paimon").Collection("mihoyobbsGoods")
	
	// 检查今天是否更新, 更新过了就直接返回
	// type=-1 数据的next_time
	var queryUpdate struct {
		Type int `bson:"type"`
		Nextime int `bson:"next_time"`
	}

	err = collection.FindOne(context.TODO(), bson.M{"type": -1}).Decode(&queryUpdate)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 找不到文档记录
			log.Error("文档不存在")
			res, err := collection.InsertOne(context.TODO(), bson.M{"type": -1, "next_time": time.Now().Local().Unix()})
			if err != nil {
				log.Error(err)
			}
			log.Info("更新记录: ", res.InsertedID)
		} else {
			return err
		}
	}

	// 时间戳转换判断是否在同一天
	nowTime := time.Now().Local().Unix()
	latestTime := int64(queryUpdate.Nextime)
	
	// 如果不在同一天
	if !IsSameDay(nowTime, latestTime, 4*3600) {
		// 先尝试删除所有的文档
		deleteResult, err := collection.DeleteMany(context.TODO(), bson.M{"type": bson.M{"$ne": -1}})
		if err != nil {
			log.Error(err)
			return err
		}
		log.Infof("Deleted %v documents in the trainers collection\n", deleteResult.DeletedCount)
	} else {
		log.Info("数据暂时不需要更新")
		return nil
	}

	// 开始写入
	ch := make(chan int)
	for i := 1; i < 6; i++ {
		go func(i int) {
			god := getGoods(i)
			if god == nil {
				return
			}
			_, err = collection.InsertMany(context.TODO(), god)
			if err != nil {
				log.Error(err)
				return
			}
			time.Sleep(3 * time.Second)
			ch <- i
		}(i)
	}
	<-ch

	// 更新type数据, type为-1的nexttime的值为该表最近更新的值
	res, err := collection.UpdateOne(context.TODO(), bson.M{"type": -1}, bson.M{"$set": bson.M{"next_time": time.Now().Local().Unix()}})
	if err != nil {
		log.Error(err)
	}
	log.Infof("更新了 %d 条记录", res.UpsertedCount)

	// 断开连接
	err = MongoDBClient.Disconnect(context.TODO())
	if err != nil {
		log.Error(err)
		return err
	}
	fmt.Println("Connection to MongoDB closed.")
	return nil
}

// 对商品列表进行遍历
func getGoods(pageNumber int) []interface{} {
	resp, err := requests.Bronya("GET", fmt.Sprintf("https://api-takumi.miyoushe.com/mall/v1/web/goods/list?app_id=1&point_sn=myb&page_size=20&page=%d&game=", pageNumber), map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
	}, nil, nil, false)
	if err != nil {
		log.Error(err)
		return nil
	}
	body := resp.Body
	jsonRes := map[string]interface{}{}
	json.Unmarshal(body, &jsonRes)
	if jsonRes["retcode"].(float64) != 0 {
		log.Info(jsonRes["message"])
		return nil
	}
	tmplist := jsonRes["data"].(map[string]interface{})["list"]
	tmp, _ := json.Marshal(tmplist)
	var goodlist []interface{}
	json.Unmarshal(tmp, &goodlist)
	length := len(goodlist)
	if length == 0 {
		log.Warn("goods is empty")
		return nil
	}
	// log.Println(goodlist)
	return goodlist
}

func GetMihoyoGoods() ([]BBSGoodsType, error) {
	clientOptions := options.Client().ApplyURI("mongodb://Nahida:Nahida1027@localhost:27017/Paimon")
	// clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	MongoDBClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	err = MongoDBClient.Ping(context.TODO(), nil)
	// 测试连接
	if err != nil {
		log.Error(err)
		return nil, err
	}
	collection := MongoDBClient.Database("Paimon").Collection("mihoyobbsGoods")

	findOptions := options.Find()
	// findOptions.SetLimit(2)

	var results []BBSGoodsType
	cur, err := collection.Find(context.TODO(), bson.M{
		"$and": []bson.M{
			{"type": bson.M{"$ne": -1}},
			{"type": bson.M{"$ne": 0}},
			{"next_time": bson.M{"$ne": 0}},
		},
	}, findOptions)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if err = cur.All(context.TODO(), &results); err != nil {
		log.Error(err)
		return nil, err
	}
	// log.Println(results)
	return results, nil
}

func MihoyoBBSGoodsForQuery(uid string) ([]interface{}, error) {
	goods, err := GetMihoyoGoods()
	if err != nil {
		return nil, err
	}

	// 最大限制50个查询
	if len(goods) > 50 {
		goods = goods[1:51]
	}

	timeStamp16 := time.Now().UnixNano() / 1e6
	timeStamp16String := strconv.FormatInt(timeStamp16, 10)
	var results []interface{}

	for n, good := range goods {
		article := tgbotapi.NewInlineQueryResultArticleMarkdown(
			timeStamp16String+fmt.Sprint(n),
			good.GoodsName,
			"",
		)
		article.Title = good.GoodsName
		fmt.Println("title", good.GoodsName)
		article.ThumbURL=good.Icon
		
		tm := int64(good.NextTime)
		ts := TimeStampToString(tm)

		// shell block
		var shell, goodType string
		if good.Type == 1 {
			goodType = "real"
			shell = fmt.Sprintf("./ShotGoods -type \"real\" -id \"%s\" -target \"%d\" -config \"config.json\"", good.GoodsID, good.NextTime)
		} else {
			goodType = "virtual"
			shell = fmt.Sprintf("./ShotGoods -type \"virtual\" -id \"%s\" -uid \"%s\" -target \"%d\" -config \"config.json\"", good.GoodsID, uid, good.NextTime)
		}
		article.InputMessageContent = tgbotapi.InputTextMessageContent{
			Text: fmt.Sprintf("%s\nType: %s\nPrice: %d\nExchangeTime: %s\nGoodsId: `%s`\nShell: `%s`", strings.Replace(good.GoodsName, "*", "\\*", -1), goodType, good.Price, ts, good.GoodsID, shell),
			ParseMode: "Markdown",
		}
		results = append(results, article)
	}
	return results, nil
}

// 判断两个时间戳是否为同一天, 请确保输入时间的时区是一样的
// 
// offset: 每天的分割点, 默认为 00:00:00 即 0s, 00:01:01 -> 61s
func IsSameDay(latestTime, nowTime, offset int64) bool {
	time.LoadLocation("Asia/Shanghai")
	lt := time.Unix(latestTime-offset, 0).Local()
	nt := time.Unix(nowTime-offset, 0).Local()
	result := true
	if nt.Day() != lt.Day() {
		result = false
	}
	return result
}

func TimeStampToString(timeStamp int64) string {
	time.LoadLocation("Asia/Shanghai")
	return time.Unix(timeStamp, 0).Local().Format("2006-01-02 15:04:05")
}

func StringToTimeStamp(dateStr string) int64 {
	loc, _  := time.LoadLocation("Asia/Shanghai")
	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)
	return timestamp.Unix()
}