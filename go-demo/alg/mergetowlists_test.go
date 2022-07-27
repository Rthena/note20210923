package alg

import (
	"fmt"
	"testing"
)

func TestTow(t *testing.T) {
	l1 := CreateLinkList([]int{1, 2, 4})
	l2 := CreateLinkList([]int{1, 3, 4})
	l := mergeTwoLists(l1, l2)
	Print(l)
}

func Print(l *ListNode) {
	cur := l
	var format string
	for cur != nil {
		format += fmt.Sprintf("%+v", cur.Val)
		cur = cur.Next
		if nil != cur {
			format += "->"
		}
	}
	fmt.Println(format)
}
func CreateLinkList(nodeVal []int) *ListNode {
	var head *ListNode
	var tail *ListNode

	for i := 0; i < len(nodeVal); i++ {
		if head == nil {
			head = &ListNode{Val: nodeVal[i]}
			tail = head
		} else {
			tail.Next = &ListNode{Val: nodeVal[i]}
			tail = tail.Next
		}
	}
	return head
}

func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l1
	} else if l2 == nil {
		return l2
	} else if l1.Val < l2.Val {
		l1.Next = mergeTwoLists(l1.Next, l2)
		return l1
	} else {
		l2.Next = mergeTwoLists(l1, l2.Next)
		return l2
	}
}
