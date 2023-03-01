package modules

// func SearchInRdsAndMgo(key string) (interface{}, error) {
// 	val, err := rds.Client.Get(key).Result()
// 	if err == redis.Nil { //继续到mgo查询
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		return nil, err
// 	} else { //在redis缓存中找到
// 		fmt.Println("val", val)
// 	}

// }
