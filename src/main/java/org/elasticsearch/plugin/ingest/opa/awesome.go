package main

import "C"

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"context"

	"bytes"
	"encoding/json"
	"github.com/open-policy-agent/opa/rego"
)

var count int
var mtx sync.Mutex

//export Add
func Add(a, b int) int {
	return a + b
}

//export Cosine
func Cosine(x float64) float64 {
	return math.Cos(x)
}

//export Sort
func Sort(vals []int) {
	sort.Ints(vals)
}

//export SortPtr
func SortPtr(vals *[]int) {
	Sort(*vals)
}

//export Log
func Log(msg string) int {
	mtx.Lock()
	defer mtx.Unlock()
	fmt.Println(msg)
	count++
	return count
}

//export LogPtr
func LogPtr(msg *string) int {
	return Log(*msg);
}

//export Eval
func Eval() int {

	ctx := context.Background()

	// Raw input data that will be used in evaluation.
	raw := `{"users": [{"id": "bob"}, {"id": "alice"}]}`
	d := json.NewDecoder(bytes.NewBufferString(raw))

	// Numeric values must be represented using json.Number.
	d.UseNumber()

	var input interface{}

	if err := d.Decode(&input); err != nil {
		panic(err)
	}

	// Create a simple query over the input.
	rego := rego.New(
		rego.Query("input.users[idx].id = user_id"),
		rego.Input(input))

	//Run evaluation.
	rs, err := rego.Eval(ctx)

	if err != nil {
		// Handle error.
	}

	// Inspect results.
	fmt.Println("len:", len(rs))
	fmt.Println("bindings.idx:", rs[1].Bindings["idx"])
	return 3

}

func main() {}