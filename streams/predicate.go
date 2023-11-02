/**
 * @author leo
 * @date 2023/4/17 11:51
 */
package streams

type Predicate[T any] interface {
	//
	// Test
	//  @Description: 预测结果
	//  @param data
	//  @return bool
	//
	Test(data T) bool
}
