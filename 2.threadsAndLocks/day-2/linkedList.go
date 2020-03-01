package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
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
	sample := make([]int, listSize)
	randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < listSize; i++ {
		sample[0] = randGen.Int()
	}
	sample2 := make([]int, listSize)
	copy(sample2, sample)
	var wg sync.WaitGroup
	wg.Add(2)
	sample1Lock := concurrentSortedList{}
	go func() {
		defer wg.Done()
		startTime := time.Now()
		fmt.Println("1 lock start:", startTime)
		for _, v := range sample {
			sample1Lock.insert1Lock(v)
		}
		fmt.Println("1 lock elapse:", time.Now().Sub(startTime))
	}()
	sampleH2HLock := concurrentSortedList{}
	go func() {
		defer wg.Done()
		startTime := time.Now()
		fmt.Println("H2H lock start:", startTime)
		for _, v := range sample {
			sampleH2HLock.insertH2H(v)
		}
		fmt.Println("H2H lock elapse:", time.Now().Sub(startTime))
	}()
	wg.Wait()
}
