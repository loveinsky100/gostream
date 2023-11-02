/**
 * @author leo
 * @date 2021/4/26 3:15 下午
 */
package routines

type Goroutine interface {
	Go(runnable func())
}
