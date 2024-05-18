package golang

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/smartystreets/goconvey/convey"
)

// https://github.com/bouk/monkey author bouke from Netherlands
// https://github.com/agiledragon/gomonkey author Xiaolong Zhang from China
// It seems like the difference between the two repo is tiny. Or the latter might have some improvement from the former.

// go test -v  golang.go -test.run Add #Test the single func `TestAdd`
// go test -v  golang.go -test.bench ForFun -test.run ForFun #Test a func in benchmark mode
// go test -bench=. -benchtime=3s #run all benchmarks, specify the test duration.

// -cpuprofile: profile.out #output the result of cpu analysis
// -memprofile: memprofile.out #Output the result of Mem analysis
// https://blog.csdn.net/weixin_34232617/article/details/91854391

/* An example for another way to use `go test`. And it's related to memory allocation from go101.
var t *[5]int64
var s []byte

func f(b *testing.B) {
	for i := 0; i < b.N; i++ {
		t = &[5]int64{}
	}
}

func g(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s = make([]byte, 32769)
	}
}

func main() {
	println(unsafe.Sizeof(*t))      // 40
	rf := testing.Benchmark(f)
	println(rf.AllocedBytesPerOp()) // 48
	rg := testing.Benchmark(g)
	println(rg.AllocedBytesPerOp()) // 40960
}
*/

// do unit test for application code:
// https://medium.com/@rishibhardwaj2010/writing-unit-test-cases-in-golang-4389b43dc57e

// stress test
func BenchmarkDirect(b *testing.B) {
	x, y := 1, 2
	for i := 0; i < b.N; i++ {
		_ = x + y
	}
}

var globalNum int

// see more at https://github.com/agiledragon/gomonkey/tree/master/test
// Add() is a simple function. you need add `-gcflags "all=-N -l"` to avoid inlining optimization.
// explain: -N disable optimizations, -l disable inline.
// run `go tool compile -help` to get more about `gcflag`.
func TestMock(t *testing.T) {
	convey.Convey("test_single", t, func() {
		// Mock Function
		// make it by replacing the func address
		patches := gomonkey.ApplyFunc(Add, func(x, y int) int {
			return 1
		})
		defer patches.Reset()

		patches.ApplyFunc(Sub, func(x, y int) int {
			return 1
		})

		convey.So(Add(1, 3), convey.ShouldEqual, 1)
		convey.So(Sub(1, 3), convey.ShouldEqual, 1)
	})

	convey.Convey("test_bulk", t, func() {
		convey.Convey(" mock method", func() {
			// Mock method
			patches := gomonkey.ApplyMethod(Adder{}, "Add", func(adder Adder, x, y int) int {
				return 1
			})
			defer patches.Reset()

			// mock private method
			patches.ApplyPrivateMethod(Adder{}, "add", func(adder Adder, x, y int) int {
				return 1
			})

			// pointer method
			patches.ApplyMethod(&PtrAdder{}, "Add", func(adder *PtrAdder, x, y int) int {
				return 1
			})

			convey.So((Adder{}).Add(1, 3), convey.ShouldEqual, 1)
			convey.So((Adder{}).add(1, 3), convey.ShouldEqual, 1)
			convey.So((&PtrAdder{}).Add(1, 3), convey.ShouldEqual, 1)
		})

		convey.Convey("mock variables", func() {
			var num int
			// Mock local variable
			patches := gomonkey.ApplyGlobalVar(&num, 666)
			defer patches.Reset()

			// Mock global variable
			patches.ApplyGlobalVar(&globalNum, 666)

			convey.So(num, convey.ShouldEqual, 666)
			convey.So(globalNum, convey.ShouldEqual, 666)
		})

		convey.Convey("mock func seq", func() {
			patches := gomonkey.ApplyFuncSeq(AddAndSub, []gomonkey.OutputCell{
				{
					Values: gomonkey.Params{1, 2},
					Times:  2,
				},
				{
					Values: gomonkey.Params{2, 1},
					Times:  2,
				},
			})
			defer patches.Reset()

			a, b := AddAndSub(1, 1)
			convey.So(a, convey.ShouldEqual, 1)
			convey.So(b, convey.ShouldEqual, 2)
			a, b = AddAndSub(1, 1)
			convey.So(a, convey.ShouldEqual, 1)
			convey.So(b, convey.ShouldEqual, 2)
			a, b = AddAndSub(1, 1)
			convey.So(a, convey.ShouldEqual, 2)
			convey.So(b, convey.ShouldEqual, 1)
		})
	})
}

func Add(x, y int) int {
	return x + y
}

func Sub(x, y int) int {
	return x - y
}

func AddAndSub(x, y int) (int, int) {
	return Add(x, y), Sub(x, y)
}

type Adder struct {
}

func (a Adder) Add(x, y int) int {
	return x + y
}

func (a Adder) add(x, y int) int {
	return x + y
}

type PtrAdder struct {
}

func (a *PtrAdder) Add(x, y int) int {
	return x + y
}

// bulk-cases test. template generated by goland-generate
func TestAdd(t *testing.T) {
	var tests = []struct {
		Name   string
		ArgA   int
		ArgB   int
		Result int
		Err    bool
	}{
		{"case1", 1, 2, 3, false},
		{"case2", 3, -2, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			got := Add(tt.ArgA, tt.ArgB)
			if got != tt.Result {
				t.Errorf("Add() = %v, want = %v", got, tt.Result)
			}
		})
	}
}
