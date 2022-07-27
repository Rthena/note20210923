package alg

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

// 方法一
func Test0(t *testing.T) {
	//输入：strs = ["flower","flow","flight"]
	//输出："fl"
	var strs = []string{"flower", "flow", "flight"}
	c := LongestCommonPrefixV1(strs)
	t.Log(c)
}

func LongestCommonPrefixV1(strs []string) string {
	strLen := len(strs)
	if strLen == 0 {
		return ""
	}
	if strLen == 1 {
		return strs[0]
	}
	var commonPrefix []string
	for i := 0; i < strLen; i++ {
		for j := i + 1; j < strLen; j++ {
			commonPrefix = append(commonPrefix, CompareString(strs[i], strs[j]))
		}
	}
	sort.Strings(commonPrefix)
	fmt.Println(commonPrefix)
	if len(commonPrefix) == 0 {
		return ""
	}
	return commonPrefix[0]
}

func CompareString(a, b string) string {
	var loopCount = math.Min(float64(len(a)), float64(len(b)))
	var prefix string
	for i := 0; i < int(loopCount); i++ {
		if a[i] == b[i] {
			prefix += string(a[i])
		} else {
			break
		}
	}
	return prefix
}

func Test1(t *testing.T) {
	ns := strconv.FormatInt(123, 10)
	t.Log(ns)
	n, _ := strconv.ParseInt("123", 10, 64)
	t.Log(n)
}

// 方法二
func Test2(t *testing.T) {
	//输入：strs = ["flower","flow","flight"]
	//输出："fl"
	var strs = []string{"flower", "flow", "flight"}
	c := LongestCommonPrefixV2(strs)
	t.Log(c)
}

func LongestCommonPrefixV2(strs []string) string {
	strLen := len(strs)
	if strLen == 0 {
		return ""
	}
	if strLen == 1 {
		return strs[0]
	}
	var prefix = strs[0]
	for i := 1; i < strLen; i++ {
		prefix = CompareString(prefix, strs[i])
		if len(prefix) == 0 {
			break
		}
	}
	return prefix
}

// 方法三
func Test3(t *testing.T) {
	//输入：strs = ["flower","flow","flight"]
	//输出："fl"
	var strs = []string{"flower", "flow", "flight", "", ""}
	c := LongestCommonPrefixV3(strs, 0, len(strs)-1)
	t.Log(c)
}

func LongestCommonPrefixV3(strs []string, start, end int) string {
	//strLen := len(strs)
	//if strLen == 0 {
	//	return ""
	//}
	//if strLen == 1 {
	//	return strs[0]
	//}

	if start == end {
		return strs[start]
	}
	mid := (end-start)/2 + start
	fmt.Printf("start=%d mid=%d end=%d\n", start, mid, end)
	_ = LongestCommonPrefixV3(strs, start, mid)
	_ = LongestCommonPrefixV3(strs, mid+1, end)

	return "" //CompareString(leftPrefix, rightPrefix)
}
