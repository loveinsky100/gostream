/**
 * @author leo
 * @date 2023/4/17 11:47
 */
package streams

import (
	"fmt"
	"github.com/loveinsky100/gostreams/routines"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	isRun := false
	Of(1, 2, 3, 4, 6).Foreach(func(v int) {
		isRun = true
	})

	if !isRun {
		t.Error("not run")
	}
}

func TestFilter(t *testing.T) {
	Of(1, 2, 3, 4, 6).Filter(func(v int) bool {
		return v > 1
	}).Foreach(func(v int) {
		if v <= 1 {
			t.Error("error")
		}
	})
}

func TestPeek(t *testing.T) {
	count := 0
	Of(1, 2, 3, 4, 6).Peek(func(v int) {
		count++
	}).Foreach(func(v int) {
	})

	if count != 5 {
		t.Error("peek error")
	}
}

func TestFindFirst(t *testing.T) {
	f := Of(1, 2, 3, 4, 6).FindFirst()
	if f != 1 {
		t.Error("error")
	}
}

func TestLimit(t *testing.T) {
	count := 0
	Of(1, 2, 3, 4, 6).Limit(3).Foreach(func(v int) {
		count++
	})
	if count != 3 {
		t.Error("error")
	}
}

func TestAnyMatch(t *testing.T) {
	v := Of(1, 2, 3, 4, 6).AnyMatch(func(v int) bool {
		return v == 3
	})

	if !v {
		t.Error("error")
	}
}

func TestAllMatch(t *testing.T) {
	v := Of(1, 2, 3, 4, 6).AllMatch(func(v int) bool {
		return v < 3
	})

	if v {
		t.Error("error")
	}
}

func TestNoneMatch(t *testing.T) {
	v := Of(1, 2, 3, 4, 6).NoneMatch(func(v int) bool {
		return v > 1000
	})

	if !v {
		t.Error("error")
	}
}

func TestMax(t *testing.T) {
	v := Of(1, 2, 3, 4, 6).Max(Numeric[int]())
	if v != 6 {
		t.Error("error")
	}
}

func TestMin(t *testing.T) {
	v := Of(1, 2, 3, 4, 6).Min(Numeric[int]())
	if v != 1 {
		t.Error("error")
	}
}

func TestSort(t *testing.T) {
	max := 6
	Of(1, 2, 3, 4, 5, 6).Sort(NumericDesc[int]()).Foreach(func(v int) {
		if max != v {
			t.Error("error")
		}

		max--
	})
}

func TestCollect(t *testing.T) {
	d := Of(1, 2, 3, 4, 6).
		Sort(NumericDesc[int]()).
		MapTo(Mapper[int, string]).(MapperStream[int, string]).
		Map(func(v int) string {
			return fmt.Sprintf("a_%d", v)
		}).
		CollectToList()
	fmt.Println(fmt.Sprintf("%+v", d))
}

func TestCollectMap(t *testing.T) {
	d := Of(1, 2, 3, 4, 6).Sort(NumericDesc[int]()).MapTo(Mapper[int, string]).(MapperStream[int, string]).Map(func(v int) string {
		return fmt.Sprintf("a_%d", v)
	}).GroupBy(func(k string) string {
		return fmt.Sprintf("k_%s", k)
	})

	fmt.Println(fmt.Sprintf("%+v", d))
}

func TestDistinct(t *testing.T) {
	d := Of(1, 2, 3, 4, 6, 6, 6, 6).Distinct().CollectToList()
	if len(d) != 5 {
		t.Error("TestDistinct")
	}
}

type Demo struct {
	Name string
}

func TestFilter2(t *testing.T) {
	d := Of(&Demo{"hello"}, &Demo{"value"}, nil).Distinct().Filter(NotNil[*Demo]()).CollectToList()
	if len(d) != 2 {
		t.Error("error")
	}
}

func TestSkip(t *testing.T) {
	d := Of(0, 1, 2, 3, 4, 5).Skip(1).CollectToList()
	if len(d) != 5 {
		t.Error("error")
	}
}

func TestCount(t *testing.T) {
	d := Of(0, 1, 2, 3, 4, 5).Skip(1).Count()
	if d != 5 {
		t.Error("error")
	}
}

type RejectHandler struct {
}

func (r *RejectHandler) Reject(callable routines.Callable) error {
	fmt.Println("reject")
	return nil
}

func TestRecover(t *testing.T) {
	pool := routines.NewQueuedGoRoutinePool(3, 1, &RejectHandler{})
	d2 := Of(0, 1, 2, 3, 4, 5).
		Pool(pool).
		Debug(Println[string]()).
		Error(func(source string, err error) {
			fmt.Println(fmt.Sprintf("source: %+v, err: %+v", source, err))
		}).
		ParallelWith(time.Second * 6).
		Filter(func(v int) bool {
			if v == 3 {
				panic("error")
			}
			time.Sleep(time.Second)
			return true
		}).
		MapTo(Mapper[int, string]).(MapperStream[int, string]).
		Map(func(v int) string {
			time.Sleep(time.Second * time.Duration(v))
			return fmt.Sprintf("p_%d", v)
		}).
		Map(func(v string) string {
			time.Sleep(time.Second * 1)
			return fmt.Sprintf("o_%s", v)
		}).
		Sort(func(o1, o2 string) bool {
			return o1 > o2
		}).
		CollectToList()
	fmt.Println(d2)
}

func TestParallel(t *testing.T) {
	pool := routines.NewQueuedGoRoutinePool(2, 10, nil)
	d := Of(1, 2, 3, 4, 6).
		Pool(pool).
		Debug(Println[string]()).
		Error(func(source string, err error) {
			fmt.Println(fmt.Sprintf("Error: source: %+v, err: %+v", source, err))
		}).
		Parallel().
		MapTo(Mapper[int, string]).(MapperStream[int, string]).
		Map(func(v int) string {
			return fmt.Sprintf("a_%d", v)
		}).
		Peek(func(v string) {
			time.Sleep(time.Second * 3)
		}).
		Sort(func(o1, o2 string) bool {
			return o1 > o2
		}).
		CollectToList()

	fmt.Println(d)
}

func TestParallel2(t *testing.T) {
	pool := routines.NewQueuedGoRoutinePool(5, 10, nil)
	d := Of(1, 2, 3, 4, 6).
		Pool(pool).
		Debug(Println[string]()).
		Error(func(source string, err error) {
			fmt.Println(fmt.Sprintf("Error: source: %+v, err: %+v", source, err))
		}).
		Parallel().
		MapTo(Mapper[int, string]).(MapperStream[int, string]).
		Map(func(v int) string {
			return fmt.Sprintf("a_%d", v)
		}).
		Peek(func(v string) {
			time.Sleep(time.Second * 3)
		}).
		Sort(func(o1, o2 string) bool {
			return o1 > o2
		}).
		GroupBy(func(v string) string {
			time.Sleep(time.Second * 1)
			return v
		})

	fmt.Println(d)
}

type City struct {
	Name string
}

type Province struct {
	Name   string
	Cities []*City
}

func (p *Province) GetCities() []*City {
	time.Sleep(time.Second)
	return p.Cities
}

func TestFlatMap(t *testing.T) {
	p1 := &Province{
		Name: "ZJ",
		Cities: []*City{
			{"HZ"}, {"TZ"}, {"WZ"},
		},
	}

	p2 := &Province{
		Name: "JX",
		Cities: []*City{
			{"WX"}, {"NJ"}, {"SZ"},
		},
	}

	cities := Of(p1, p2).
		Debug(Println[string]()).
		Parallel().
		MapTo(Mapper[*Province, *City]).(MapperStream[*Province, *City]).
		FlatMap(func(p *Province) []*City {
			return p.GetCities()
		}).
		MapTo(Mapper[*City, string]).(MapperStream[*City, string]).
		Map(func(p *City) string {
			return p.Name
		}).
		Filter(func(v string) bool {
			return strings.HasSuffix(v, "Z")
		}).
		CollectToList()

	if 4 != len(cities) {
		t.Error("error")
	}

	fmt.Println(cities)
}

func TestReduce(t *testing.T) {
	d := Of(0, 1, 2, 3, 4, 5).Reduce(Sum[int]()).(int)
	if d != 15 {
		t.Error("error")
	}

	data := []interface{}{
		1, 2, 3, "122",
	}

	Of(data...).Map(func(v any) any {
		return fmt.Sprintf("%+v", v)
	}).Count()
}

func BenchmarkOrigin(b *testing.B) {
	values := []string{"Apple", "Star", "Hello", "World", "Attribute"}
	for i := 0; i < b.N; i++ {
		maxLength := 0
		for _, value := range values {
			if !strings.HasPrefix(value, "A") {
				continue
			}

			length := len(value)
			if maxLength < length {
				maxLength = length
			}
		}
	}
}

func BenchmarkStream(b *testing.B) {
	values := []string{"Apple", "Star", "Hello", "World", "Attribute"}
	for i := 0; i < b.N; i++ {
		SliceOf(values).
			Filter(func(v string) bool {
				return strings.HasPrefix(v, "A")
			}).
			MapTo(Mapper[string, int]).(MapperStream[string, int]).
			Map(func(v string) int {
				return len(v)
			}).
			Max(Numeric[int]())
	}
}
