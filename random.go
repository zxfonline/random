// Copyright 2016 zxfonline@sina.com. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package random

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func RandInt32(min int32, max int32) int32 {
	if max <= min {
		max = min + 1
	}
	var base int32 = 0
	if min < 0 {
		base = -min
		min += base
		max += base
	}
	return -base + min + rand.Int31n(max-min)
}
func RandInt64(min int64, max int64) int64 {
	if max <= min {
		max = min + 1
	}
	var base int64 = 0
	if min < 0 {
		base = -min
		min += base
		max += base
	}
	return -base + min + rand.Int63n(max-min)
}
func RandInt(min int, max int) int {
	if max <= min {
		max = min + 1
	}
	var base int = 0
	if min < 0 {
		base = -min
		min += base
		max += base
	}
	return -base + min + rand.Intn(max-min)
}

func main1() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Println(GetRandomNumber("1", random))
	fmt.Println(GetRandomNumber("1~10", random))
	fmt.Println(GetRandomNumber("1~10,2,2~5", random))
	fmt.Println(GetRandomNumber("1,2,4", random))
	fmt.Println(GetRandomNumber("2~10:40", random))
	fmt.Println(GetRandomNumber("1:40", random))
	fmt.Println(GetRandomNumber("1:20,1~4:30,4:500", random))
	fmt.Println(GetRandomNumbers("2~10:40#10:20,10~45:30,40~80:500", random))
}

/**
 * 获得随机数
 *
 * @param numbers
 * @param n 要随机多n个数
 * @return
 */
func GetRandomValues(numbers []int, n int) []int {
	size := len(numbers)
	filter := make([]int, size)
	copy(filter, numbers)
	if size == 0 || n >= size {
		return filter
	}
	list := make([]int, 0, n)
	for i := 0; i < n; i++ {
		index := rand.Intn(len(filter))
		list = append(list, filter[index])
		filter = append(filter[:index], filter[index+1:]...)
	}
	return list
}

/**
 * 加权随机数 数值抽取器
 *
 * @param args 以“#”劈分数组然后再在数组元素中每一位获取一个随机数
 * @param random
 * @see #GetRandomNumber(String, Random)
 * @return
 */
func GetRandomNumbers(args string, random *rand.Rand) []int {
	strs := strings.Split(args, "#")
	size := len(strs)
	ints := make([]int, 0, size)
	for i := 0; i < size; i++ {
		ints = append(ints, GetRandomNumber(strs[i], random))
	}
	return ints
}

/**
 * 数值抽取器(从枚举值,范围随机值,定值 中随机抽出一个值)
 * 枚举值(支持单个出现概率)：1,2,4 或 1:10,2:30,4:30
 * 范围值： 1~10
 * 定值：12
 * 支持混合使用 如：1~10,44~89,2~5 又如 2~10:40
 * 支持概率后缀 代表该值被抽取出来的几率 值越高被抽出的概率越大 并不限定后缀值的范围 如：1:20,1~4:30,4:500
 *
 * @param args
 * @param random
 * @return
 */
func GetRandomNumber(args string, random *rand.Rand) int {
	var err error
	var value int
	//	if strings.Index(args, ",") > 0 {
	values := strings.Split(args, ",")
	size := len(values)
	numbers := make([]int, size)
	if strings.Index(args, ":") > 0 { // 1:20,1~4:30,4:500
		weights := make([]int, size)
		var valuesStr, weightStr string
		var weightSum int
		for i := 0; i < size; i++ {
			if endIndex := strings.Index(values[i], ":"); endIndex <= 0 {
				panic(fmt.Errorf("invalid args:%v,err:%v", args, values[i]))
			} else {
				valuesStr = string(values[i][:endIndex])
			}
			weightStr = string(values[i][strings.Index(values[i], ":")+1:])
			if value, err = strconv.Atoi(weightStr); err != nil {
				panic(fmt.Errorf("invalid args:%v,err:%v", args, err))
			} else {
				weights[i] = value
				weightSum += value
			}
			if numbers[i], err = average(valuesStr, random); err != nil {
				panic(fmt.Errorf("invalid args:%v,err:%v", args, err))
			}
		}
		ranNum := random.Intn(weightSum)
		for i := 0; i < size; i++ {
			ranNum -= weights[i]
			if ranNum < 0 {
				return numbers[i]
			}
		}
		if value, err = strconv.Atoi(args); err != nil {
			panic(fmt.Errorf("invalid args:%v,err:%v", args, err))
		} else {
			return value
		}
	} else { // 1~10,44~89,2~5
		for i := 0; i < size; i++ {
			if numbers[i], err = average(values[i], random); err != nil {
				panic(fmt.Errorf("invalid args:%v,err:%v", args, err))
			}
		}
		return numbers[random.Intn(size)]
	}
}

func average(args string, random *rand.Rand) (int, error) { // 1~10
	if strings.Index(args, "~") > 0 {
		tmp := strings.Split(args, "~")
		if v1, err1 := strconv.Atoi(tmp[0]); err1 != nil {
			return 0, err1
		} else if v2, err2 := strconv.Atoi(tmp[1]); err2 != nil {
			return 0, err2
		} else {
			return v1 + random.Intn(v2+1-v1), nil
		}
	} else {
		return strconv.Atoi(args)
	}
}
