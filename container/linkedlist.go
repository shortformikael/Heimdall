package container

import (
	"errors"
	"fmt"
)

type ListNode struct {
	data interface{}
	next *ListNode
}

type LinkedList struct {
	head *ListNode
}

func (l *LinkedList) Append(data interface{}) {
	newNode := &ListNode{data: data}
	if l.head == nil {
		l.head = newNode
		return
	}
	current := l.head
	for current.next != nil {
		current = current.next
	}
	current.next = newNode
}

func (l *LinkedList) Prepend(data interface{}) {
	newNode := &ListNode{data: data, next: l.head}
	l.head = newNode
}

func (l *LinkedList) Pop() {
	current := l.head
	if current == nil {
		return
	} else if current.next == nil {
		current = nil
		return
	}

	previous := current
	current = current.next
	for current != nil {
		if current.next == nil {
			previous.next = nil
			return
		} else {
			previous = current
			current = current.next
		}
	}

}

func (l *LinkedList) Print() {
	current := l.head
	for {
		if current != nil {
			fmt.Println(current.data)
			current = current.next
		} else {
			return
		}
	}
}

func (l *LinkedList) PrintList() {
	current := l.head
	for current != nil {
		//fmt.Println(reflect.TypeOf(current.data))
		switch c := current.data.(type) {
		case string:
			fmt.Printf(" - %s \n", c)
		case *MenuItem:
			fmt.Printf("  *  %s \n", c.Name)
		case *LinkedList:
			c.PrintList()
		default:
			fmt.Printf(" - Unknown \n")
		}
		current = current.next
	}
}

func (l *LinkedList) Delete(s string) error {
	if s == l.head.data {
		l.head = l.head.next
		return nil
	}
	current := l.head.next
	previous := l.head
	for current != nil {
		if current.data == s {
			previous.next = current.next
			return nil
		}
		previous = current
		current = current.next
	}
	return errors.New("ITEM DOESNT EXIST IN LIST")
}

func (l *LinkedList) DeleteAt(n int) error {
	if 0 > n || n >= l.Length() {
		return errors.New("INTEGER OUT OF BOUNDS")
	} else if n == 0 {
		l.head = l.head.next
		return nil
	}
	current := l.head.next
	previous := l.head
	for i := 0; i < n; i++ {
		previous = current
		current = current.next
	}
	previous.next = current.next
	return nil
}

func (l *LinkedList) Get(n int) (interface{}, error) {
	if 0 > n || n >= l.Length() {
		return "", errors.New("INTERGER OUT OF BOUNDS")
	}
	current := l.head
	for i := 0; i < n; i++ {
		if current != nil {
			current = current.next
		}
	}
	return current.data, nil
}

func (l *LinkedList) Search(data string) (interface{}, error) {
	current := l.head
	for current != nil {
		switch i := current.data.(type) {
		case string:
			if i == data {
				return i, nil
			}
		case *MenuItem:
			if i.Name == data {
				return i, nil
			}
		}
		current = current.next
	}
	return nil, errors.New("DOES NOT EXIST IN ARRAY")
}

func (l *LinkedList) Exists(s interface{}) bool {
	current := l.head
	for current != nil {
		if current.data == s {
			return true
		}
		current = current.next
	}
	return false
}

func (l *LinkedList) Length() int {
	count := 0
	current := l.head
	for current != nil {
		count++
		current = current.next
	}
	return count
}

func (l *LinkedList) GetArray() *[]any {
	var rArray []any = []any{}
	current := l.head
	for current.data != nil {
		rArray = append(rArray, current.data)
		current = current.next
	}
	return &rArray
}

// DOESNT WORK
func (l *LinkedList) AddArray(arr []interface{}) {
	for i := 0; i < len(arr); i++ {
		l.Append(arr[i])
	}
}
