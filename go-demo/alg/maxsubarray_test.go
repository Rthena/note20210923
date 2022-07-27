package alg

import (
	"fmt"
	"testing"
)

/**
输入：nums = [-2,1,-3,4,-1,2,1,-5,4]
输出：6
解释：连续子数组 [4,-1,2,1] 的和最大，为 6 。
*/
func TestMaxSubArray(t *testing.T) {
	//nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	c := make(chan int)
	c <- 1
	select {
	case <-c:
		fmt.Println("aaa")
		//panic(1)
	}

}
