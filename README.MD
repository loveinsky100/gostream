Golang stream
## How to use

```go
// distinct
Of(1, 2, 3, 4, 6, 6, 6, 6).Distinct().CollectToList()

// find min value
Of(1, 2, 3, 4, 6).Min(Numeric[int]())

// check if all macth
Of(1, 2, 3, 4, 6).AllMatch(func(v int) bool {
    return v < 3
})

// filter map
SliceOf(values).
Filter(func(v string) bool {
    return strings.HasPrefix(v, "A")
}).
MapTo(Mapper[string, int]).(MapperStream[string, int]).
Map(func(v string) int {
    return len(v)
}).
Max(Numeric[int]())
```

## Support
| Method        | Example                                  | Description                                        |
|---------------|------------------------------------------|----------------------------------------------------|
| Debug         | Debug(Println[string]())                 | 打印后续节点信息（debug之后的）                                 |
| Error         | Error(func(source string, err error) {}) | 打印错误信息，比如panic等                                    |
| Recover       | Recover()                                | 尝试恢复，此时会将错误数据给到 Error ， 另外如果设置并行则默认添加              |
| Pool          | Pool(xxx)                                | 配置并行的协程池，不设置则默认使用全局                                |
| Filter        | Filter(PredicateHandler)                 | 过滤数据                                               |
| Skip          | Skip(index)                              | 跳过前N个数据                                            |
| Sort          | Sort(ComparatorHandler)                  | 对数据进行排序                                            |
| Limit         | Limit(size)                              | 限制输出数量                                             |
| Distinct      | Distinct()                               | 去重，比如字符，数字等可以直接使用                                  |
| DistinctWith  | DistinctWith(FunctionHandler)            | 复杂去重，自行传入去重逻辑                                      |
| Parallel      | Parallel()                               | 开启并行，后续的节点会进行并行处理，注意：对于后续的有状态节点会自动添加串型 - 节点 - 并行链路 |
| ParallelWith  | ParallelWith(timeout)                    | 开启并行，并限制超时                                         |
| Sequential    | Sequential()                             | 开启并行后可以使用此方法重新开始串型处理                               |
| Peek          | Peek(ConsumeHandler)                     | 遍历数据，然后原样返回                                        |
| Foreach       | Foreach(ConsumeHandler)                  | 终止节点，遍历数据                                          |
| AnyMatch      | AnyMatch(PredicateHandler)               | 是否任意一个数据匹配                                         |
| AllMatch      | AllMatch(PredicateHandler)               | 是否所有数据都匹配                                          |
| NoneMatch     | NoneMatch(PredicateHandler)              | 是否没有任何数据匹配                                         |
| FindFirst     | FindFirst()                              | 返回第一个达到的数据                                         |
| Max           | Max(ComparatorHandler)                   | 返回最大的数据                                            |
| Min           | Min(ComparatorHandler)                   | 返回最小的数据                                            |
| Count         | Count()                                  | 统计数量                                               |
| Collect       | Collect(Collector)                       | 收集数据                                               |
| CollectToList | CollectToList()                          | 收集数据变成数组                                           |
| GroupBy       | GroupBy(FunctionHandler)                 | 转化成map                                             |
| Reduce        | Reduce(Reduce)                           | Reduce                                             |
| Map           | Map(FunctionHandler)                     | 转化数据                                               |
| MapTo         | MapTo()                                  | 转化为其他类型的数据                                         |
| Join          | Join(Pipeline)                           | 加入自定义节点                                            |
| FlatMap       | FlatMap(FunctionHandler)                 | 转化数据，将流入的数组数据转化为item数据                             |

## Other
由于Golang对于范型的支持不够，对于Map操作有些特殊处理。
```go
    Of(0, 1, 2, 3, 4, 5).
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
```
MapTo(Mapper[int, string]).(MapperStream[int, string])
这里的作用是将int类型转化为string类型，后续节点都变化成string的流




