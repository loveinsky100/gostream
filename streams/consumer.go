/**
 * @author leo
 * @date 2023/4/17 12:10
 */
package streams

type Consumer[T any] interface {
	Accept(t T)
}
