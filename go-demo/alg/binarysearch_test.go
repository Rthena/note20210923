package alg

import "testing"

func TestBinary(t *testing.T) {
	i := BinarySearchIterator([]int{1, 2, 3, 4, 5, 6, 7, 8}, 8)
	t.Log(i)
}

func BinarySearchIterator(arr []int, num int) int {
	var low int
	var high = len(arr)
	for low <= high {
		mid := low + (high-low)/2
		if num == arr[mid] {
			return mid
		}
		if num > arr[mid] {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}
	return -1
}

func BinarySearchRecursion(input []int, num int) int {
	//var res int
	//var l = len(input)
	//var mid = l / 2
	//var low =
	//var high
	//
	//// 1
	//if num == input[res] || input[low] == input[mid] || input[high] == input[mid] {
	//	return res
	//}
	//
	//// 右边
	//if num > input[mid] {
	//	BinarySearch(input[mid:], num)
	//} else {
	//	BinarySearch(input[:mid], num)
	//}
	//
	//return res
	return 0
}
