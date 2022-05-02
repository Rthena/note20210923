package linkedlist

import "testing"

var l *LinkedList

func init() {
	n5 := &ListNode{next: nil, data: 5}
	n4 := &ListNode{next: n5, data: 4}
	n3 := &ListNode{next: n4, data: 3}
	n2 := &ListNode{next: n3, data: 2}
	n1 := &ListNode{next: n2, data: 1}
	l = &LinkedList{
		head:   n1,
		length: 0,
	}
}

func TestReverse(t *testing.T) {
	l.Print()
	l.Reverse()
	l.Print()
}
