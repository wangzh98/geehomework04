package homework4

import (
	"sync"
	"time"
)

type Number struct {
	// key 值为timestamp value为次数
	Buckets map[int64]*numberBucket
	Mutex   *sync.RWMutex
}

type numberBucket struct {
	Value float64
}

// Increment increments the number in current timeBucket.
func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	// 从map获取当前numberBucket ，没有则新建
	b := r.getCurrentBucket()
	b.Value += i
	// 移除过去10s以前的数据
	r.removeOldBuckets()
}

// 获取当前时间戳对应的 numberBucket，不存在时创建
func (r *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - 10

	for timestamp := range r.Buckets {
		// TODO: configurable rolling window
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

func (r *Number) Sum(now time.Time) float64 {
	sum := float64(0)

	// 加的写锁提升效率
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		// TODO: configurable rolling window
		if timestamp >= now.Unix()-10 {
			sum += bucket.Value
		}
	}

	return sum
}
