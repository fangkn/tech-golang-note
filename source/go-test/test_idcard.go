package main

import (
	"fmt"
	"strconv"
	"time"
)

// ParseChineseID 解析中国的身份证号码，提取出生日期和性别
func ParseChineseID(id string) (string, int, error) {
	if len(id) != 18 {
		return "", 0, fmt.Errorf("身份证号码长度不正确")
	}

	// 提取出生日期
	birthday := id[6:14]
	year, err := strconv.Atoi(birthday[0:4])
	if err != nil {
		return "", 0, fmt.Errorf("出生日期解析错误")
	}
	month, err := strconv.Atoi(birthday[4:6])
	if err != nil {
		return "", 0, fmt.Errorf("出生日期解析错误")
	}
	day, err := strconv.Atoi(birthday[6:8])
	if err != nil {
		return "", 0, fmt.Errorf("出生日期解析错误")
	}

	// 提取性别
	gender := id[16:17]
	genderInt, err := strconv.Atoi(gender)
	if err != nil {
		return "", 0, fmt.Errorf("性别解析错误")
	}
	var genderType int
	if genderInt%2 == 0 {
		genderType = 2
	} else {
		genderType = 1
	}

	// 格式化出生日期
	birthdayStr := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")

	return birthdayStr, genderType, nil
}

func main() {
	id := "11010519491231002X"
	birthday, gender, err := ParseChineseID(id)
	if err != nil {
		fmt.Println("解析错误:", err)
		return
	}
	fmt.Println("出生日期:", birthday)
	fmt.Println("性别:", gender)
}
