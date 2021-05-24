package zone

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func Test_New_Handle(t *testing.T) {
	handler := handler{
		index: 0,
		inUse: make(map[uint16]bool),
	}

	newHandler = make(chan *handlerPetition, 1500)
	removeHandler = make(chan *handlerPetition, 1500)
	queryHandler = make(chan *handlerPetition, 1500)

	go handler.handleWorker()

	h, err := newHandle()
	if err != nil {
		t.Fatal(err)
	}

	if h != 1 {
		t.Fatalf("unexpected newHandle value %v", h)
	}
}

func Test_Remove_Handle(t *testing.T) {
	handler := handler{
		index: 0,
		inUse: make(map[uint16]bool),
	}

	newHandler = make(chan *handlerPetition, 1500)
	removeHandler = make(chan *handlerPetition, 1500)
	queryHandler = make(chan *handlerPetition, 1500)

	go handler.handleWorker()

	h, err := newHandle()
	if err != nil {
		t.Fatal(err)
	}

	removeHandle(h)

	if handleExists(h) {
		t.Fatalf("newHandle %v should not exist", h)
	}
}

func Test_10000_New_Handles(t *testing.T) {
	handler := handler{
		index: 0,
		inUse: make(map[uint16]bool),
	}

	newHandler = make(chan *handlerPetition, 1500)
	removeHandler = make(chan *handlerPetition, 1500)
	queryHandler = make(chan *handlerPetition, 1500)

	go handler.handleWorker()

	var (
		wg  sync.WaitGroup
		sem = make(chan int, 1000)
	)

	rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		sem <- 1

		go func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(30)))

			defer wg.Done()
			_, err := newHandle()
			if err != nil {
				t.Error(err)
			}
			<-sem
		}()
	}

	wg.Wait()

	h, err := newHandle()
	if err != nil {
		t.Error(err)
	}

	if h != 10001 {
		t.Errorf("unexpected newHandle value %v", h)
	}
}

func Test_Create_10000_New_Handles_And_Remove_Them(t *testing.T) {
	handler := handler{
		index: 0,
		inUse: make(map[uint16]bool),
	}

	newHandler = make(chan *handlerPetition, 1500)
	removeHandler = make(chan *handlerPetition, 1500)
	queryHandler = make(chan *handlerPetition, 1500)

	go handler.handleWorker()

	var (
		wg  sync.WaitGroup
		sem = make(chan int, 1000)
	)

	rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		sem <- 1

		go func(iwg *sync.WaitGroup) {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(30)))

			defer iwg.Done()

			h, err := newHandle()
			if err != nil {
				t.Error(err)
			}

			iwg.Add(1)
			go func(h uint16) {
				defer iwg.Done()
				removeHandle(h)
			}(h)

			<-sem
		}(&wg)
	}

	wg.Wait()

	h, err := newHandle()
	if err != nil {
		t.Error(err)
	}

	for i := 0; i < 10000; i++ {
		if handleExists(uint16(i)) {
			t.Error("handle should not exist")
		}
	}

	if h != 10001 {
		t.Errorf("unexpected newHandle value %v", h)
	}
}
