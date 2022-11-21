package quatsch

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	b "github.com/SimonRichardson/colossus/quatsch/pool/bson"
	"gopkg.in/mgo.v2/bson"
)

func create(amount int, pool Pool, dst chan<- []bson.ObjectId) {
	var idents []bson.ObjectId
	for i := 0; i < amount; i++ {
		ident, err := b.Bson(pool.Get())
		if err != nil {
			log.Fatal(err)
		}
		idents = append(idents, ident)
	}
	dst <- idents
}

func TestCollision(t *testing.T) {
	var (
		amount   = 99999
		routines = 12

		wg  = sync.WaitGroup{}
		dst = make(chan []bson.ObjectId, routines)

		maxBuffer               = 99999
		maxInsertionPerDuration = int64(1000000)

		pool = New(b.New(maxBuffer, time.Second, maxInsertionPerDuration))
	)

	wg.Add(routines)
	go func() { wg.Wait(); close(dst) }()

	for i := 0; i < routines; i++ {
		go func() {
			defer wg.Done()

			create(amount, pool, dst)
		}()
	}

	idents := make([]bson.ObjectId, 0)
	for k := range dst {
		idents = append(idents, k...)
	}

	if expected := amount * routines; expected != len(idents) {
		t.Errorf("Total mismatch (expected %d, recieved %d)", expected, len(idents))
	}

	if err := isUnique(idents); err != nil {
		log.Fatal(err)
	}
}

func isUnique(idents []bson.ObjectId) error {
	m := make(map[bson.ObjectId]int, len(idents))
	for k, v := range idents {
		if index, ok := m[v]; ok {
			return fmt.Errorf("Collision found at %d and %d with key %s", index, k, v)
		}
		m[v] = k
	}
	return nil
}
