package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

func Pic(dx, dy int) [][]uint8 {
	rs := make([][]uint8, dy)
	for i := 0; i < dy; i++ {
		rs[i] = make([]uint8, dx)
		for j := range rs[i] {
			rs[i][j] = 255
		}
	}
	return rs
}

func WordCount(s string) map[string]int {
	rs := make(map[string]int)
	words := strings.Fields(s)
	for _, word := range words {
		rs[word] = len(word)
	}
	return rs
}

//(0, 1, 1, 2, 3, 5, ...).
func fibonacci() func() int {
	pre := 0
	cur := 0
	return func() int {
		nex := cur + pre
		if cur == 0 {
			cur = 1
		} else if cur == 1 {
			pre = 1
		} else {
			pre = cur
			cur = nex
		}

		return nex
	}
}

type IPAddr [4]byte

// TODO: Add a "String() string" method to IPAddr.

func (iPAddrp IPAddr) String() string {
	str := make([]string, len(iPAddrp))
	for i, v := range iPAddrp {
		str[i] = fmt.Sprint(v)
	}
	return strings.Join(str, ",")
}

type rot13Reader struct {
	r io.Reader
}

func (r rot13Reader) Read(bytes []byte) (int, error) {
	n, err := r.r.Read(bytes)
	for i := range bytes {
		bytes[i] = bytes[i] + 13
	}
	return n, err
}

func sum(s []int, c chan int, id int) {
	sum := 0
	for _, v := range s {
		if id == 1 {
			time.Sleep(1 * time.Microsecond)

		}
		fmt.Printf("Thread %v - V = %v \n", id, v)
		sum += v
	}
	c <- sum // send sum to c
}

func producer(c chan int) {
	s := []int{7, 2, 8, -9, 4, 0}
	for _, v := range s {
		fmt.Println(v)
		c <- v
		time.Sleep(1 * time.Second)
	}
	close(c)
}

func receiver(c chan int) {
	for v := range c {
		fmt.Println(v)
		time.Sleep(2 * time.Second)
	}
}

func main1() {
	c := make(chan int)
	go producer(c)
	receiver(c)
}
