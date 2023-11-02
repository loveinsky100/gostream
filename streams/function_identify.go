/**
 * @author leo
 * @date 2023/4/17 12:54
 */
package streams

func Identity[T any]() FunctionHandler[T, T] {
	return func(t T) T {
		return t
	}
}
