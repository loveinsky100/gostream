/**
 * @author leo
 * @date 2020/8/28 3:56 下午
 */
package routines

import (
	"runtime"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	pool := NewGoRoutinePool(10, &CallableRejectedHandler{
		Handler: func(callable Callable) {
		},
	})

	future, err := pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	v, err := future.GetWithTimeout(time.Millisecond * 200)
	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	if 1 != v {
		t.Errorf("value not equal")
	}
}

func TestCancel(t *testing.T) {
	pool := NewGoRoutinePool(10, &CallableRejectedHandler{
		Handler: func(callable Callable) {
		},
	})

	for index := 0; index < 10; index++ {
		future, err := pool.AddHandler(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 100)
			return 1, nil
		})

		if err != nil {
			t.Errorf("an error occured: %+v", err)
		}

		if index%2 == 0 {
			time.Sleep(time.Millisecond * 10)
		}

		cancelSuccess := future.Cancel()
		_, err = future.GetWithTimeout(time.Millisecond * 200)
		if cancelSuccess {
			if index%2 == 0 {
				t.Errorf("an error occured: cancel failed %t", cancelSuccess)
				break
			}

			if err == nil {
				t.Errorf("an error occured: cancel failed %t", cancelSuccess)
			}
		} else {
			if err != nil {
				t.Errorf("an error occured: cancel failed %t", cancelSuccess)
			}
		}

	}
}

func TestReject(t *testing.T) {
	pool := NewGoRoutinePool(1, &CallableRejectedHandler{
		Handler: func(callable Callable) {

		},
	})

	_, err := pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	_, err = pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err == nil {
		t.Errorf("an error occured: reject failed")
	}
}

func TestTimeout(t *testing.T) {
	pool := NewGoRoutinePool(1, &CallableRejectedHandler{
		Handler: func(callable Callable) {

		},
	})

	future, err := pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	_, err = future.GetWithTimeout(time.Millisecond * 50)
	if err == nil {
		t.Errorf("an error occured: timeout failed")
	}
}

func TestTimeout2(t *testing.T) {
	pool := NewGoRoutinePool(1, nil)

	future, err := pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	_, err = future.GetWithTimeout(time.Millisecond * 50)
	if err == nil {
		t.Errorf("an error occured: timeout failed")
	}
}

func TestIsDone(t *testing.T) {
	pool := NewGoRoutinePool(1, &CallableRejectedHandler{
		Handler: func(callable Callable) {

		},
	})

	future, err := pool.AddHandler(func() (interface{}, error) {
		time.Sleep(time.Millisecond * 100)
		return 1, nil
	})

	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	_, err = future.GetWithTimeout(time.Millisecond * 200)
	if err != nil {
		t.Errorf("an error occured: %+v", err)
	}

	if future.Status() != FINISH {
		t.Errorf("an error occured: is done error")
	}
}

func TestGoRouterCounter(t *testing.T) {
	startCount := runtime.NumGoroutine()
	pool := NewGoRoutinePool(10, &CallableRejectedHandler{
		Handler: func(callable Callable) {

		},
	})

	for index := 0; index < 15; index++ {
		pool.AddHandler(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 1000)
			return 1, nil
		})

		// sleep make sure router add into runtime
		time.Sleep(time.Millisecond * 20)
		addCount := runtime.NumGoroutine() - startCount
		if index < 10 {
			if addCount != (1 + index) {
				t.Errorf("an error occured: add count error: %d", runtime.NumGoroutine())
			}
		} else {
			if addCount != 10 {
				t.Errorf("an error occured: add count error %d", runtime.NumGoroutine())
			}
		}
	}
}

func TestGoRouterExit(t *testing.T) {
	startCount := runtime.NumGoroutine()
	pool := NewGoRoutinePool(20, &CallableRejectedHandler{
		Handler: func(callable Callable) {

		},
	})

	for index := 0; index < 15; index++ {
		pool.AddHandler(func() (interface{}, error) {
			time.Sleep(time.Millisecond * 200)
			return 1, nil
		})
	}

	addCount := runtime.NumGoroutine() - startCount
	if addCount != 15 {
		t.Errorf("an error occured: add count error: %d", runtime.NumGoroutine())
	}

	// sleep make sure router add into runtime
	time.Sleep(time.Millisecond * 500)
	addCount = runtime.NumGoroutine() - startCount
	if addCount != 0 {
		t.Errorf("an error occured: add count error %d", runtime.NumGoroutine())
	}
}

//
//
// func TestGoRouterGroupExecute(t *testing.T)  {
// 	group := NewGoRouterGroup()
//
// 	value := atomic.In64(0)
//
// 	for index := 0; index < 100; index ++ {
// 		group.Add(func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	err := group.Next()
// 	if value.Load() != 100 {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
//
// func TestGoRouterGroupError(t *testing.T)  {
// 	group := NewGoRouterGroup()
// 	group.Add(func() error {
// 		return errors.New("error")
// 	})
//
// 	group.Add(func() error {
// 		return nil
// 	})
//
// 	err := group.Next()
// 	if nil == err || strings.Compare(err.Error(), "error") != 0 {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
//
// func TestGoRouterGroupTimeout(t *testing.T)  {
// 	group := NewGoRouterGroup()
//
// 	value := atomic.NewInt64(0)
//
// 	for index := 0; index < 10; index ++ {
// 		group.Add(func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	group.AddWithTimeout(10 * time.Millisecond, func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	err := group.Next()
// 	if nil == err {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }

// func TestGoRouterGroupTimeout2(t *testing.T)  {
// 	group := NewGoRouterGroup()
//
// 	value := atomic.NewInt64(0)
//
// 	for index := 0; index < 10; index ++ {
// 		group.AddWithTimeout(50 * time.Millisecond, func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	err := group.Next()
// 	if nil == err {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
//
// func TestGoRouterGroupTimeout3(t *testing.T)  {
// 	group := NewGoRouterGroup()
//
// 	value := atomic.NewInt64(0)
//
// 	for index := 0; index < 10; index ++ {
// 		group.AddWithTimeout(110 * time.Millisecond, func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	err := group.Next()
// 	if nil != err {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
//
// func TestGoRouterGroupTimeout5(t *testing.T)  {
// 	group := NewGoRouterGroup()
//
// 	value := atomic.NewInt64(0)
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	for index := 0; index < 10; index ++ {
// 		timeout := 110 - index
// 		group.AddWithTimeout(time.Duration(timeout) * time.Millisecond, func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	group.AddWithTimeout(10 * time.Millisecond, func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	err := group.Next()
// 	if nil == err {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
//
// func TestGoRouterGroupReject(t *testing.T)  {
// 	pool := NewGoRoutinePool(5, &CallableRejectedHandler {
// 		Handler: func(callable Callable) {
//
// 		},
// 	})
//
// 	group := NewGoRouterGroup()
//
// 	value := atomic.NewInt64(0)
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	for index := 0; index < 10; index ++ {
// 		timeout := 120 - index
// 		group.AddWithTimeout(time.Duration(timeout) * time.Millisecond, func() error {
// 			time.Sleep(100 * time.Millisecond)
// 			value.Inc()
// 			return nil
// 		})
// 	}
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	group.AddWithTimeout(10 * time.Millisecond, func() error {
// 		return nil
// 	})
//
// 	group.Add(func() error {
// 		time.Sleep(100 * time.Millisecond)
// 		return nil
// 	})
//
// 	err := group.ExecuteInPool(pool)
// 	if nil == err {
// 		t.Errorf("an error occured: %+v", err)
// 	}
// }
