package linkedlist

import "fmt"

type ListNode struct {
	next *ListNode
	data interface{}
}

type LinkedList struct {
	head   *ListNode
	length uint
}

func (this *LinkedList) Reverse() {
	if this.head == nil {
		return
	}
	var pre *ListNode = nil
	cur := this.head
	for cur != nil {
		nxt := cur.next
		cur.next = pre
		pre = cur
		cur = nxt
	}
	this.head = pre
}

func (this *LinkedList) hasCycle() {

}

//打印链表
func (this *LinkedList) Print() {
	cur := this.head
	format := ""
	for nil != cur {
		format += fmt.Sprintf("%+v", cur.data)
		cur = cur.next
		if nil != cur {
			format += "->"
		}
	}
	fmt.Println(format)
}
