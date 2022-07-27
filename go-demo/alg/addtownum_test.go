package alg

import (
	"fmt"
	"strconv"
	"testing"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbersV1(l1 *ListNode, l2 *ListNode) *ListNode {
	l1 = reverse(l1)
	l2 = reverse(l2)
	towSum := helper(l1) + helper(l2)
	l := revString2LinkList(strconv.Itoa(towSum))
	return l
}

func reverse(l *ListNode) *ListNode {
	if l == nil {
		return nil
	}
	var pre *ListNode
	var cur = l
	for cur != nil {
		//nxt := cur.Next
		//cur.Next = pre
		//cur = nxt
		//pre = cur

		nxt := cur.Next
		cur.Next = pre
		pre = cur
		cur = nxt
	}
	return pre
}

func helper(l *ListNode) int {
	var format string
	var cur = l
	for cur != nil {
		format += fmt.Sprintf("%+v", cur.Val)
		cur = cur.Next
	}
	num, err := strconv.Atoi(format)
	if err != nil {
		fmt.Println("err", err)
		return 0
	}
	return num
}

func revString2LinkList(str string) (l *ListNode) {
	strLen := len(str)
	if strLen < 0 {
		return nil
	}
	var head *ListNode
	var tail *ListNode
	for i := strLen - 1; i >= 0; i-- {
		v, _ := strconv.Atoi(string(str[i]))
		if head == nil {
			head = &ListNode{Val: v}
			tail = head
		} else {
			tail.Next = &ListNode{Val: v}
			tail = tail.Next
		}
	}
	return head
}

/**
输入：l1 = [2,4,3], l2 = [5,6,4]
输出：[7,0,8]
解释：342 + 465 = 807.
*/
func TestAdd(t *testing.T) {
	a1 := &ListNode{Val: 3, Next: nil}
	a2 := &ListNode{Val: 4, Next: a1}
	a3 := &ListNode{Val: 2, Next: a2}

	b1 := &ListNode{Val: 4, Next: nil}
	b2 := &ListNode{Val: 6, Next: b1}
	b3 := &ListNode{Val: 5, Next: b2}
	r := addTwoNumbersV2(a3, b3)
	t.Log(r)
}

// 方法二
func addTwoNumbersV2(l1 *ListNode, l2 *ListNode) *ListNode {
	var head *ListNode
	var tail *ListNode
	var carry int
	for l1 != nil || l2 != nil {
		n1, n2 := 0, 0
		if l1 != nil {
			n1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			n2 = l2.Val
			l2 = l2.Next
		}
		sum := n1 + n2 + carry
		sum, carry = sum%10, sum/10
		if head == nil {
			head = &ListNode{Val: sum}
			tail = head
		} else {
			tail.Next = &ListNode{Val: sum}
			tail = tail.Next
		}
	}
	if carry > 0 {
		tail.Next = &ListNode{Val: carry}
	}
	return head
}

func TestAdd2(t *testing.T) {
	var head *ListNode
	var tail *ListNode

	ints := []int{1, 2, 3}
	for i := 0; i < len(ints); i++ {
		if head == nil {
			head = &ListNode{Val: ints[i]}
			tail = head
		} else {
			tail.Next = &ListNode{Val: ints[i]}
			tail = tail.Next
		}
	}
	cur := head
	for cur != nil {
		t.Log(cur.Val)
		cur = cur.Next
	}
}
