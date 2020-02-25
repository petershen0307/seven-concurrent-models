package main

import (
	"flag"
	"sync"
)

type concurrentNode struct {
	value int
	prev  *concurrentNode
	next  *concurrentNode
	lock  sync.Mutex
}

type concurrentSortedList struct {
	head *concurrentNode
	lock sync.Mutex
}

func (list *concurrentSortedList) toList() []int {
	outList := make([]int, 0)
	c := list.head
	for c != nil {
		outList = append(outList, c.value)
		c = c.next
	}
	return outList
}

func (list *concurrentSortedList) insertH2H(v int) {
	newNode := &concurrentNode{
		value: v,
		prev:  nil,
		next:  nil,
	}
	list.lock.Lock()
	if list.head == nil {
		list.head = newNode
		list.lock.Unlock()
		return
	}
	list.lock.Unlock()
	current := &list.head
	for true {
		(*current).lock.Lock()
		if v < (*current).value {
			// head
			if (*current).prev == nil {
				newNode.next = (*current)
				(*current).lock.Unlock()
				(*current) = newNode
			} else {
				(*current).prev.lock.Lock()
				(*current).prev.next = newNode
				(*current).prev.lock.Unlock()
				(*current).lock.Unlock()
			}
			break
		}
		if nil == (*current).next {
			(*current).lock.Unlock()
			break
		}
		(*current).lock.Unlock()
		current = &(*current).next
	}
	if newNode.prev == nil && newNode.next == nil {
		(*current).lock.Lock()
		newNode.prev = *current
		(*current).next = newNode
		(*current).lock.Unlock()
	}
}

func (list *concurrentSortedList) insert1Lock(v int) {
	list.lock.Lock()
	defer list.lock.Unlock()
	newNode := &concurrentNode{
		value: v,
		prev:  nil,
		next:  nil,
	}
	if list.head == nil {
		list.head = newNode
		return
	}
	current := &list.head
	// small to big
	for true {
		if v < (*current).value {
			// head
			if (*current).prev == nil {
				newNode.next = (*current)
				(*current) = newNode
			} else {
				(*current).prev.next = newNode
			}
			break
		}
		if nil == (*current).next {
			break
		}
		current = &(*current).next
	}
	if newNode.prev == nil && newNode.next == nil {
		newNode.prev = *current
		(*current).next = newNode
	}
}

func main() {
	var listSize int
	flag.IntVar(&listSize, "s", 1000, "give a size for list")
	flag.Parse()
}
