/**
 * @author leo
 * @date 2023/4/24 21:06
 */
package streams

type Collector[T any, R any] interface {
	Sink[T]
	Result() R
}
