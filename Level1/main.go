package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	TEST_ARR     = []int{3, 4, 5, 6, 3, 2, 1, 12, 2, 5, 7, 1, 9}
	TEST_STR     = "abcdefghijklmnop"
	TEST_ARR_STR = "snow dog sun"
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

func task13() {
	var a, b int32 = 4, 6
	fmt.Printf("Variables before swapping: %d %d\n", a, b)
	a, b = b, a
	fmt.Printf("Variables after swapping: %d %d\n", a, b)
}

func ValType(arr []interface{}) []string {
	var slice []string
	for _, v := range arr {
		switch reflect.TypeOf(v).Kind().String() {
		case "int":
			slice = append(slice, "int")
		case "string":
			slice = append(slice, "string")
		case "bool":
			slice = append(slice, "bool")
		case "chan":
			slice = append(slice, "chan int")
		}
	}
	return slice
}

func task14() []string {
	// int, string, bool, channel
	var slice = []interface{}{3, "stringVal", true, make(chan struct{})}
	return ValType(slice)
}

func createHugeString(n int) string {
	return string(make([]byte, n))
}

func someFunc() string {
	var justString string
	v := createHugeString(1 << 10)
	justString = v[:100]
	v = string([]byte(justString))
	return v
}

func task15() {
	ans := someFunc()
	fmt.Println(ans)
}

func partition(arr []int, low, high int) int {
	pivot := arr[high] // выбор опорного элемента
	i := low - 1       // индекс меньшего элемента

	for j := low; j <= high-1; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}

// функция quickSort реализует алгоритм быстрой сортировки
func quickSort(arr []int, low, high int) {
	if low < high {
		// pi - индекс опорного элемента, arr[pi] находится на правильном месте
		pi := partition(arr, low, high)

		// Рекурсивно сортируем элементы перед разделением и после разделения
		quickSort(arr, low, pi-1)
		quickSort(arr, pi+1, high)
	}
}

func task16() {
	fmt.Printf("Array before sorting: %v", TEST_ARR)
	quickSort(TEST_ARR, 0, len(TEST_ARR)-1)
	fmt.Printf("Array after sorting: %v", TEST_ARR)
}

func binarySearch(l, r, target int) (int, error) {
	if l >= r && TEST_ARR[l] != target {
		return -1, errors.New("not found such position")
	}
	mid := l + (r-l)/2
	if TEST_ARR[mid] == target {
		return mid, nil
	} else if TEST_ARR[mid] > target {
		return binarySearch(l, mid-1, target)
	}
	return binarySearch(mid+1, r, target)
}

func task17() int {
	target, n := 5, len(TEST_ARR)-1
	quickSort(TEST_ARR, 0, n)
	sort.Ints(TEST_ARR)
	pos, err := binarySearch(0, n, target)
	if err != nil {
		log.Printf("binary search err: %v", err)
		return pos
	}
	return pos
}

func task18() int {
	counter, goroutinesSize := 0, 1000
	ch := make(chan int)
	wg, wgReader := sync.WaitGroup{}, sync.WaitGroup{}
	mtx := sync.Mutex{}

	wgReader.Add(1)
	go func() {
		defer wgReader.Done()
		for v := range ch {
			mtx.Lock()
			counter += v
			mtx.Unlock()
		}
	}()

	for i := 1; i <= goroutinesSize; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			ch <- num
		}(i)
	}

	wg.Wait()
	close(ch)
	wgReader.Wait()
	return (2*counter)/goroutinesSize - 1
}

func swapString(runeArr []string) {
	l, r := 0, len(runeArr)-1
	for l < r {
		runeArr[l], runeArr[r] = runeArr[r], runeArr[l]
		l++
		r--
	}
}

func task19() string {
	runeArr := strings.Split(TEST_STR, "")
	swapString(runeArr)
	return strings.Join(runeArr, "")
}

func task20() string {
	runeArr := strings.Split(TEST_ARR_STR, " ")
	swapString(runeArr)
	return strings.Join(runeArr, " ")
}

type Round interface {
	getRadius() float64
}

// round interface
type SquarePeg struct {
	Width float64
}

// round interface

type RoundHole struct {
	Radius float64
}

// round interface

type RoundPeg struct {
	radius float64
}

func (p RoundPeg) getRadius() float64 {
	return p.radius
}

type SquarePegAdapter struct {
	squarePeg SquarePeg
}

func (p SquarePegAdapter) getRadius() float64 {
	return p.squarePeg.Width * math.Sqrt(2) / 2
}

func (r RoundHole) fits(peg Round) error {
	if r.Radius < peg.getRadius() {
		return errors.New("too small roundhole")
	}
	return nil
}

func task21(width, radius float64) error {
	roundHole := RoundHole{Radius: radius}
	squarePeg := SquarePeg{Width: width}
	peg := SquarePegAdapter{squarePeg}
	if err := roundHole.fits(peg); err != nil {
		return fmt.Errorf("error during pushing: %v", err)
	}
	log.Print("Connected!!!")
	return nil
}

type Sign string

const (
	SUBTRACTION    Sign = "-"
	ADDITION       Sign = "+"
	MULTIPLICATION Sign = "*"
	DIVISION       Sign = "/"
)

func task22() (*big.Int, error) {
	var sign Sign
	first, second, result := new(big.Int), new(big.Int), new(big.Int)
	scanner := bufio.NewReader(os.Stdin)
	str, err := scanner.ReadString('\n')
	if err != nil {
		return nil, err
	}
	arr := strings.Split(str, " ")
	first.SetString(arr[0], 10)
	second.SetString(arr[2], 10)
	sign = Sign(arr[1])
	fmt.Scanf("%s %v %s", &first, &sign, &second)
	switch sign {
	case SUBTRACTION:
		result = result.Sub(first, second)
	case ADDITION:
		result = result.Add(first, second)
	case MULTIPLICATION:
		result = result.Mul(first, second)
	case DIVISION:
		result = result.Div(first, second)
	default:
		return nil, errors.New("unknown sign")
	}
	return result, nil
}

func task23(arr []int, pos int) error {
	if len(arr) == 0 || pos < 0 || pos >= len(arr) {
		return errors.New("wrong initialization parameters")
	}
	if pos == 0 {
		arr = arr[1:]
		return nil
	}
	if pos == len(arr)-1 {
		arr = arr[:pos]
		return nil
	}
	leftHalf, rightHalf := arr[:pos], arr[pos+1:]
	leftHalf = append(leftHalf, rightHalf...)
	arr = leftHalf
	return nil
}

type Point struct {
	x, y float64
}

func PointConstructor(x, y float64) Point {
	return Point{
		x, y,
	}
}

func (p Point) getX() float64 {
	return p.x
}

func (p Point) getY() float64 {
	return p.y
}

func getLength(p1, p2 Point) float64 {
	currX := p1.getX() - p2.getX()
	currY := p1.getY() - p2.getY()
	return math.Sqrt(currX*currX + currY*currY)
}

func sleep(n int, wg *sync.WaitGroup) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
	}
	wg.Done()
}

func task25() {
	var n int
	fmt.Scanf("%d\n", &n)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go sleep(n, &wg)
	wg.Wait()
	fmt.Printf("Sleep function completed for %d seconds", n)
}

func task26() []bool {
	scanner := bufio.NewScanner(os.Stdin)
	var ans []bool
	for scanner.Scan() {
		curr := scanner.Text()
		if len(curr) == 0 {
			break
		}
		curr = strings.ToLower(strings.TrimSpace(curr))
		runeStr := []rune(curr)
		m := make(map[rune]struct{})
		flag := true
		for _, v := range runeStr {
			if _, ok := m[v]; ok {
				flag = false
				break
			}
			m[v] = struct{}{}
		}
		ans = append(ans, flag)
	}
	return ans
}

func main() {
	var w, r float64
	fmt.Scanf("%f %f\n", &w, &r)
	if err := task21(w, r); err != nil {
		log.Fatal(err)
	}
}
