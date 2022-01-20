package main

import "C"

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"

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

//export decrypt
func decrypt(encString string, secretKeyring string, passphrase string) *C.char {
    //... your code here
    var str string = "returning string"
    return C.CString(str)
}

//export Eval
func Eval(in string) *C.char {

	ctx := context.Background()

	// Raw input data that will be used in evaluation.
	raw := fmt.Sprintf(`{"user": "%s"}`, in)
	d := json.NewDecoder(bytes.NewBufferString(raw))

	// Numeric values must be represented using json.Number.
	d.UseNumber()

	var input interface{}

	if err := d.Decode(&input); err != nil {
		panic(err)
	}

	// Create a simple query over the input.
	rego := rego.New(
		rego.Query("data.example.x"),
		rego.Input(input),
		rego.Module("example.rego",
			`package example

default x = false
	
x = input.user == "alice"
`,
		))

	//Run evaluation.
	rs, err := rego.Eval(ctx)

	if err != nil {
		// Handle error.
		panic(err)
	}

	// Inspect results.
    var str string = fmt.Sprintf("%v", rs[0].Expressions[0].Value)
    return C.CString(str)

}

func main() {}