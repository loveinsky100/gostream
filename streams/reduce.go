/**
 * @author leo
 * @date 2023/4/27 13:01
 */
package streams

type Reduce[T any, R any] interface {
	Collector[T, R]
}
