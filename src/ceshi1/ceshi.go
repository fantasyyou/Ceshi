package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

type Card struct {
	Alice  string `json:"alice"`
	Bob    string `json:"bob"`
	Result int    `json:"result"`
}

type Card_List struct {
	Matches []Card `json:"matches"`
}

func main() {

	//读取文件
	file, err := os.Open("C:\\match_1.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	result := string(buffer)
	fmt.Println("bytes read: ", bytesread)

	//反序列化
	var msgs Card_List
	err = json.Unmarshal([]byte(result), &msgs)
	if err != nil {
		fmt.Println("Can't decode json message", err)
	}

	start := time.Now()
	for i := 0; i < len(msgs.Matches); i++ {
		card_a := msgs.Matches[i].Alice
		card_b := msgs.Matches[i].Bob
		msgs.Matches[i].Result = judge_duizi(card_a, card_b)
	}
	tc := time.Since(start) //计算耗时
	fmt.Printf("time cost = %v\n", tc)

	//序列化
	str, err := json.Marshal(msgs) //json实现序列化
	if err != nil {
		fmt.Println("Can't decode json message", err)
	}

	//写入文件
	//data := []byte(strings.Replace(string(strings.Replace(string(str), "{", "{\n", -1)), "}", "}\n", -1))
	data := []byte(string(strings.Replace(string(str), "},", "},\n", -1)))
	err = ioutil.WriteFile("ceshi.json", data, 0666)
	if err != nil {
		fmt.Println("Can't decode json message", err)
	}
}

//判断是否为同花，返回1则是，返回0则不是
func judge_tonghua(card string, card_number [5]int) int {
	cardList := []rune(card)
	for i := 1; i < 5; i++ {
		if string(cardList[1]) != string(cardList[2*i+1]) {
			return judge_shunzi(card_number, 0) //表示0牌型不为同花
		}
	}
	return judge_shunzi(card_number, 1) //1表示牌型为同花
}

//判断完同花之后，判断是否为顺子
//不是则为同花，是则为同花顺，之
//判断是同花顺还是皇家同花顺
func judge_shunzi(card [5]int, result int) int {
	card = popupSort(card)
	for i := 4; i > 0; i-- {
		if card[i] != (card[i-1] - 1) {
			if i == 1 && card[0] == 14 && card[1] == 5 { //判断是否2345A
				break
			}
			if result == 1 { //判断牌型是否为同花
				return 6 //不是顺子是同花
			} else {
				return 1 //既不是同花也不是顺子
			}
		}
	}

	if result == 1 { //判断牌型是否为同花
		if card[4] == 10 {
			return 10 //最大同花顺
		} else {
			return 9 //同花顺
		}
	} else {
		return 5 //是顺子不是同花
	}
}

func judge_duizi(carda string, cardb string) int {
	card_a := get_card_value(carda) //得到a的牌的点数
	card_b := get_card_value(cardb) //得到b的牌的点数

	//统计牌的点数以及张数
	var value_a, value_b, number_a, number_b [5]int
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if card_a[i] == value_a[j] {
				number_a[j]++
				break
			}

			if value_a[j] == 0 {
				value_a[j] = card_a[i]
				number_a[j] = 1
				break
			}
		}
		for j := 0; j < 5; j++ {
			if card_b[i] == value_b[j] {
				number_b[j]++
				break
			}

			if value_b[j] == 0 {
				value_b[j] = card_b[i]
				number_b[j] = 1
				break
			}
		}
	}

	var result_a, result_b int

	//判断是否有对子，无对子判断是否为同花或者顺子，
	if number_a[4] == 1 {
		result_a = judge_tonghua(carda, value_a)
	} else {
		result_a = 4 //无对子，给牌值一个最大值4
	}

	//判断是否有对子，无对子判断是否为同花或者顺子
	if number_b[4] == 1 {
		result_b = judge_tonghua(cardb, value_b)
	} else {
		result_b = 4 //无对子，给牌值一个最大值4
	}

	//判断牌型是否为4+1或者3+2
	if number_a[0]+number_a[1] == 5 {
		result_a = 8 //是则给牌值一个最大值8
	}
	if number_b[0]+number_b[1] == 5 {
		result_b = 8 //是则给牌值一个最大值8
	}

	//判断牌值大小，a大于b返回1，小于返回2
	if result_a > result_b {
		return 1
	}
	if result_a < result_b {
		return 2
	}

	//对牌组先根据牌的点数排序，之后在根据张数排序
	for i := 0; i < 5; i++ {
		for j := 3; j >= i; j-- {
			if value_a[j] < value_a[j+1] {
				number_a[j], number_a[j+1] = number_a[j+1], number_a[j]
				value_a[j], value_a[j+1] = value_a[j+1], value_a[j]
			}
			if number_a[j] < number_a[j+1] {
				number_a[j], number_a[j+1] = number_a[j+1], number_a[j]
				value_a[j], value_a[j+1] = value_a[j+1], value_a[j]
			}

			if value_b[j] < value_b[j+1] {
				number_b[j], number_b[j+1] = number_b[j+1], number_b[j]
				value_b[j], value_b[j+1] = value_b[j+1], value_b[j]
			}
			if number_b[j] < number_b[j+1] {
				number_b[j], number_b[j+1] = number_b[j+1], number_b[j]
				value_b[j], value_b[j+1] = value_b[j+1], value_b[j]
			}
		}
	}

	if result_a == 5 || result_a == 9 {
		if value_a[0] == 14 && value_b[1] == 5 {
			value_a[0] = 5
			value_a[0] = 4
			value_a[0] = 3
			value_a[0] = 2
			value_a[0] = 1
		}
	}

	if result_b == 5 || result_b == 9 {
		if value_b[0] == 14 && value_b[1] == 5 {
			value_b[0] = 5
			value_b[1] = 4
			value_b[2] = 3
			value_b[3] = 2
			value_b[4] = 1
		}
	}

	//通过牌的张数比较大小
	for i := 0; i < 5; i++ {
		if number_a[i] > number_b[i] {
			return 1
		}
		if number_a[i] < number_b[i] {
			return 2
		}
	}

	//通过点数比较大小
	for i := 0; i < 5; i++ {
		if value_a[i] > value_b[i] {
			return 1
		}
		if value_a[i] < value_b[i] {
			return 2
		}
	}

	//点数和张数一样，牌型相同，返回3
	return 0
}

//冒泡排序
func popupSort(card [5]int) [5]int {
	for i := 0; i < 5; i++ {
		for j := 3; j >= i; j-- {
			if card[j] < card[j+1] {
				card[j], card[j+1] = card[j+1], card[j]
			}
		}
	}
	return card
}

//返回牌值，不可优化
func get_card_value(card string) [5]int {
	cardList := []rune(card)
	var result [5]int
	for i := 0; i < 5; i++ {
		if string(cardList[2*i]) == "T" {
			result[i] = 10
			continue
		}
		if string(cardList[2*i]) == "J" {
			result[i] = 11
			continue
		}
		if string(cardList[2*i]) == "Q" {
			result[i] = 12
			continue
		}
		if string(cardList[2*i]) == "K" {
			result[i] = 13
			continue
		}
		if string(cardList[2*i]) == "A" {
			result[i] = 14
			continue
		} else {
			number, err := strconv.Atoi(string(cardList[2*i]))
			result[i] = number
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return result
}
