package alg

import (
	"fmt"
	"testing"
)

/**
输入: nums = [0,1,0,3,12]
输出: [1,3,12,0,0]
*/
func TestMove(t *testing.T) {
	var nums = []int{0, 1, 0, 3, 12}
	z := make([]int, 0)
	n := make([]int, 0)
	for _, num := range nums {
		if num == 0 {
			z = append(z, 0)
			continue
		}
		n = append(n, num)
	}
	nums = append(n, z...)
	t.Log(nums)
}

func TestMove2(t *testing.T) {
	var nums = []int{0, 1, 0, 3, 12}
	doMoveZeroes(nums)
}
func doMoveZeroes(nums []int) {
	var slow = 0
	var fast = 0
	for fast < len(nums) {
		if nums[fast] != 0 {
			nums[slow], nums[fast] = nums[fast], nums[slow]
			slow++
		}
		fast++
	}
	fmt.Println(nums)
}
