package alg

import "testing"

/**
输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
*/
func TestTowSum(t *testing.T) {
	var nums = []int{3, 3} //{3, 2, 4} //{2, 7, 11, 15}
	r := towSumV2(nums, 6)
	t.Log(r)
}

// 方法1
func towSumV1(nums []int, target int) (res []int) {
	l := len(nums)
	if l == 0 {
		return res
	}
	for i := 0; i < l; i++ {
		for j := i + 1; j < l; j++ {
			if target == nums[i]+nums[j] {
				res = append(res, i, j)
			}
		}
	}
	return res
}

// 方法2
func towSumV2(nums []int, target int) (res []int) {
	l := len(nums)
	if l == 0 {
		return res
	}
	var m = make(map[int]int)
	for i := 0; i < l; i++ {
		if idx, ok := m[target-nums[i]]; ok {
			return append(res, idx, i)
		}
		m[nums[i]] = i
	}
	return res
}
