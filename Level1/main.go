package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
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

func Square(x int, ch chan int) {
	ch <- x * x
}

func task2() []int {
	arr := []int{2, 4, 6, 8, 10}
	ch := make(chan int, len(arr))
	for _, v := range arr {
		go Square(v, ch)
	}
	ans := make([]int, 0)
	for i := 0; i < len(arr); i++ {
		ans = append(ans, <-ch)
	}
	close(ch)
	return ans
}

func workerCreation(ch <-chan int, wg *sync.WaitGroup) {
	for v := range ch {
		fmt.Println(v * v)
		wg.Done()
	}
}

func task2_1() {
	arr := []int{2, 4, 6, 8, 10}
	wg := sync.WaitGroup{}
	jobs := make(chan int)

	for i := 1; i <= 5; i++ {
		go workerCreation(jobs, &wg)
	}

	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		jobs <- arr[i]
	}

	wg.Wait()
	close(jobs)
}

func task3_1(arr []int) int {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	return sum
}

func task3_2(arr []int) int {
	wg := &sync.WaitGroup{}
	mtx := &sync.Mutex{}
	sum := 0
	for i := 0; i < len(arr); i++ {
		wg.Add(1)
		func(num int) {
			wg.Done()

			mtx.Lock()
			defer mtx.Unlock()
			sum += (num * num)
		}(arr[i])
	}
	wg.Wait()
	return sum
}

func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Println("worker", id, "processing job", j)
		time.Sleep(1 * time.Second)
	}
}

func task4() {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scanln(&numWorkers)

	jobs := make(chan int, 100)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	go func() {
		counter := 1
		for {
			select {
			case <-stopChan:
				fmt.Println("\nПолучен сигнал завершения. Остановка отправки работ.")
				close(jobs)
				return
			default:
				jobs <- counter
				counter++
				time.Sleep(500 * time.Millisecond) // Имитация задержки при получении данных
			}
		}
	}()

	<-stopChan

	wg.Wait()
	fmt.Println("Программа завершила работу.")
}

func main() {
}
