package main

import (
	"fmt"
	"sync"
)

func main() {
	// Create a pool with a maximum size of 2 objects
	pool := NewMyObjectPool(2)

	// Get an object from the pool
	obj1 := pool.Get()
	obj1.Value = 42
	fmt.Printf("Object 1 Value: %d\n", obj1.Value)

	// Get another object (this should reuse obj1)
	obj2 := pool.Get()
	fmt.Printf("Object 2 Value: %d\n", obj2.Value) // Should print 0 because we reset it

	// Get a third object (this should reuse obj2)
	obj3 := pool.Get()
	obj3.Value = 100
	fmt.Printf("Object 3 Value: %d\n", obj3.Value)

	// Return obj3 to the pool
	pool.Put(obj1)
	pool.Put(obj2)
	pool.Put(obj3)

	// Get a fourth object (this should reuse obj3)
	obj4 := pool.Get()
	fmt.Printf("Object 4 Value: %d\n", obj4.Value) // Should print 0 because we reset it

	// Return obj4 to the pool
	pool.Put(obj4)

	// Exceed the pool's capacity
	obj5 := pool.Get()
	pool.Put(obj5) // This will be discarded because the pool is full

	// Put another object to see the discard message
	obj6 := pool.Get()
	pool.Put(obj6) // This will also be discarded
}

// MyObject represents the object we're pooling.
type MyObject struct {
	Value int
}

// MyObjectPool is a custom pool for managing MyObject instances.
type MyObjectPool struct {
	mu      sync.Mutex
	pool    []*MyObject // Slice to store pooled objects
	maxSize int         // Maximum size of the pool
}

// NewMyObjectPool creates a new MyObjectPool with the specified maximum size.
func NewMyObjectPool(maxSize int) *MyObjectPool {
	return &MyObjectPool{
		pool:    make([]*MyObject, 0, maxSize),
		maxSize: maxSize,
	}
}

// Get retrieves an object from the pool or creates a new one if the pool is empty.
func (p *MyObjectPool) Get() *MyObject {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.pool) == 0 {
		// Pool is empty, create a new object
		fmt.Println("Creating a new object. Pool size: ", len(p.pool))
		return &MyObject{}
	}

	// Get an object from the pool
	fmt.Println("Reusing an object from the pool. Pool size: ", len(p.pool))
	obj := p.pool[len(p.pool)-1]
	p.pool = p.pool[:len(p.pool)-1] // Remove object from the pool slice

	return obj
}

// Put returns an object to the pool. If the pool is full, the object is discarded.
func (p *MyObjectPool) Put(obj *MyObject) {
	p.mu.Lock()
	defer p.mu.Unlock()
	fmt.Println("Put Called, before Pool size: ", len(p.pool))
	if len(p.pool) >= p.maxSize {
		// Pool is full, discard the object
		fmt.Println("Pool is full. Discarding the object.")
		return
	}

	// Reset object state if necessary
	obj.Value = 0

	// Add the object back to the pool
	p.pool = append(p.pool, obj)
	fmt.Println("Object returned to the pool. Pool size: ", len(p.pool))
}
