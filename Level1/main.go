package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type Human struct {
	Name  string
	Age   int
	Phone string
}

type Action struct {
	Action string
	Human
}

func (h *Human) CelebrateBirthday() {
	h.Age++
}

func (a Action) SaySomething(str string) string {
	return fmt.Sprintf("Human with name %s %s word %s!", a.Name, a.Action, str)
}

func task1() {
	human := Action{
		Action: "say something",
		Human: Human{
			Name:  "Jack",
			Age:   22,
			Phone: "1337",
		},
	}
	human.CelebrateBirthday()
	fmt.Println(human.Age)
	human.SaySomething("Hello World")
}

func workerCreation(ch <-chan int, wg *sync.WaitGroup) {
	for n := range ch {
		fmt.Println(n * n)
		wg.Done()
	}
}

func task2() {
	arr := []int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	ch := make(chan int)
	// запускаем пул воркеров
	for i := 0; i < 5; i++ {
		go workerCreation(ch, &wg)
	}
	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		ch <- arr[i]
	}
	wg.Wait()
}

func task2_1() {
	arr := []int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	mtx := sync.Mutex{}
	ans := make([]int, 0)
	for i := 0; i < len(arr); i++ {
		go func(num int) {
			wg.Add(1)
			square := num * num
			mtx.Lock()
			ans = append(ans, square)
			defer mtx.Unlock()
			defer wg.Done()
		}(arr[i])
	}
	wg.Wait()
	for i := 0; i < len(ans); i++ {
		fmt.Println(ans[i])
	}
}

func task3() {
	arr := []int{2, 4, 6, 8, 10}
	sum := 0
	ch := make(chan int)
	for _, num := range arr {
		go func(i int) {
			ch <- i * i
		}(num)
	}

	for range arr {
		sum += <-ch
	}
	fmt.Println(sum)
}

func task3_1() {
	arr := []int{2, 4, 6, 8, 10}
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}
	sum := 0
	for _, num := range arr {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			square := i * i
			mtx.Lock()
			defer mtx.Unlock()
			sum += square
		}(num)
	}
	wg.Wait()
	fmt.Println(sum)
}

func workerPool(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Worker with id %d received %d\n", id, num)
		time.Sleep(1 * time.Second)
	}
}

func task4() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ch := make(chan int, 100)
	var wg sync.WaitGroup

	var n int
	fmt.Print("Print number of pool workers: ")
	fmt.Scanln(&n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go workerPool(i+1, ch, &wg)
	}

	go func() {
		counter := 1
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Stopping job generator")
				close(ch)
				return
			default:
				ch <- counter
				counter++
				time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
			}
		}
	}()

	sig := <-sigChan
	fmt.Println("\nReceived signal:", sig)
	cancel()

	wg.Wait()
	fmt.Println("All workers finished.")
}

func worker5(id int, inputCh <-chan int, outputCh chan<- string) {
	for num := range inputCh {
		time.Sleep(500 * time.Millisecond)
		outputCh <- fmt.Sprintf("Channel with id %d and value %d\n", id, num)
	}
}

func task5() {
	const (
		poolSize = 5
	)
	timer := time.NewTimer(10 * time.Second)
	inputCh := make(chan int)
	outputCh := make(chan string)
	wg := sync.WaitGroup{}
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			worker5(i, inputCh, outputCh)
		}(i)
	}
	go func() {
		counter := 0
		for {
			select {
			case <-timer.C:
				fmt.Println("Timeout")
				close(inputCh)
				return
			case inputCh <- counter:
				counter++
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(outputCh)
	}(&wg)
	for v := range outputCh {
		fmt.Println(v)
	}
	fmt.Print("All workers finished.")
}

func worker6_1(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping job generator")
			return
		}
	}
}

func task6_1() {
	// timer with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go worker6_1(ctx, &wg)

	wg.Wait()
}

func worker6_2(ch <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ch:
			fmt.Println("Stopping job generator ")
			return
		default:
			fmt.Println("It still working!!!")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func closeChannel(ch chan struct{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Stopping background  context")
			close(ch)
			return
		}
	}
}

func task6_2() {
	// done channel
	doneCh := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(1)
	go worker6_2(doneCh, &wg)
	closeChannel(doneCh)
	wg.Wait()
	fmt.Println("All workers finished.")
}

func worker6_3(ch chan int) {
	for v := range ch {
		fmt.Printf("Received %d\n", v)
	}
}

func task6_3() {
	// closing channels
	// create unbuffered channel
	ch := make(chan int)
	go worker6_3(ch)
	for i := 0; i < 5; i++ {
		ch <- i * i
		time.Sleep(1 * time.Second)
	}
	close(ch)
	fmt.Println("All workers finished.")
}

type SafeMap struct {
	m   map[int]int
	mtx *sync.Mutex
}

func (s *SafeMap) Write(k int, counter *atomic.Int32) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[k] += int(counter.Load())
	counter.Add(1)
}

func task7() {
	m := SafeMap{m: make(map[int]int), mtx: &sync.Mutex{}}
	wg := sync.WaitGroup{}
	var counter atomic.Int32
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Write(i, &counter)
		}(i)
	}
	wg.Wait()
	for k, v := range m.m {
		fmt.Println(k, v)
	}
}

func setBit(bit, val byte, n int64) int64 {
	if val == 0 {
		n &= ^(1 << bit)
	} else {
		n |= (1 << bit)
	}
	return n
}

func task8() {
	var n int64 = 29
	fmt.Println(setBit(2, 0, n))
}

func filter(x int) int {
	return x * 2
}

func worker9(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range ch {
		fmt.Println(v)
	}
}

func task9() {
	inCh := make(chan int)
	outputCh := make(chan int)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	wg := sync.WaitGroup{}
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker9(outputCh, &wg)
	}

	go func() {
		defer close(outputCh)
		for v := range inCh {
			filtered := filter(v)
			outputCh <- filtered
		}
	}()

	cnt := 0
loop:
	for {
		select {
		case <-ctx.Done():
			fmt.Print("End of input numbers\n")
			break loop
		case inCh <- cnt:
			cnt++
			fmt.Printf("Number %d pulled into input channel\n", cnt)
			time.Sleep(500 * time.Millisecond)
		}
	}
	close(inCh)
	wg.Wait()
}

func task10() {
	arr := []float64{-25.4, -27.0, 13.0, 19.0,
		15.5, 24.5, -21.0, 32.5}
	m := make(map[int][]float64)
	for i := 0; i < len(arr); i++ {
		key := (int)(arr[i]/10) * 10
		m[key] = append(m[key], arr[i])
	}
	for k, v := range m {
		fmt.Printf("Key %d with value %v\n", k, v)
	}
}

func intersection(arr1, arr2 []int) map[int]struct{} {
	m1, ansMap := make(map[int]struct{}), make(map[int]struct{})
	for i := 0; i < len(arr1); i++ {
		m1[arr1[i]] = struct{}{}
	}
	for i := 0; i < len(arr2); i++ {
		if _, ok := m1[arr2[i]]; ok {
			ansMap[arr2[i]] = struct{}{}
		}
	}
	return ansMap
}

func task11() {
	set1 := []int{3, 4, 5, 6, 7, 8, 9}
	set2 := []int{1, 2, 3, 5, 6, 7, 10, 11, 12}
	intersectionMap := intersection(set1, set2)
	fmt.Println(intersectionMap)
}

type linuxCommandsSet = map[string]struct{}

func task12() linuxCommandsSet {
	set := linuxCommandsSet{}
	str := []string{"cat", "cat", "dog", "cat", "tree"}
	for _, v := range str {
		set[v] = struct{}{}
	}
	return set
}

func main() {
	fmt.Print(task12())
}
