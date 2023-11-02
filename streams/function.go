/**
 * @author leo
 * @date 2023/4/17 12:00
 */
package streams

type Function[T any, R any] interface {
	Apply(data T) R
}
