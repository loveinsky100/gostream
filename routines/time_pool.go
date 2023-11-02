/**
 * @author leo
 * @date 2021/4/26 9:24 下午
 */
package routines

import (
	"sync"
	"time"
)

var timerPool sync.Pool
var timerRw sync.RWMutex

func init() {
	timerPool = sync.Pool{
		New: func() interface{} {
			return nil
		},
	}
}

func NewTimer(duration time.Duration) *time.Timer {
	timerRw.RLock()
	cacheData := timerPool.Get()
	if nil != cacheData {
		timerRw.RUnlock()
		timer := cacheData.(*time.Timer)
		timer.Reset(duration)
		return timer
	}

	timerRw.RUnlock()

	timerRw.Lock()
	defer timerRw.Unlock()
	// double check
	cacheData = timerPool.Get()
	if nil != cacheData {
		timer := cacheData.(*time.Timer)
		timer.Reset(duration)
		return timer
	}

	timer := time.NewTimer(duration)
	return timer
}

func Stop(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
			{
			}
		default:
			{
			}
		}
	}

	timerPool.Put(timer)
}
