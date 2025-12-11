

来源： [https://www.cnblogs.com/yjf512/p/16308469.html](https://www.cnblogs.com/yjf512/p/16308469.html)

```go
package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// DistributeLockRedis 基于 Redis 的分布式可重入锁，自动续租
type DistributeLockRedis struct {
	key       string             // 锁的 key
	expire    int64              // 锁的过期时间（秒）
	status    bool               // 上锁成功标识
	cancelFun context.CancelFunc // 用于取消自动续租协程
	redis     redis.Conn         // Redis 连接
	mu        sync.Mutex         // 保护 status 状态的互斥锁
}

// NewDistributeLockRedis 创建分布式锁实例
func NewDistributeLockRedis(redisConn redis.Conn, key string, expire int64) *DistributeLockRedis {
	return &DistributeLockRedis{
		key:    key,
		expire: expire,
		redis:  redisConn,
	}
}

// TryLock 尝试加锁
func (dl *DistributeLockRedis) TryLock(ctx context.Context) error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	// 调用 lock 获取锁
	if err := dl.lock(ctx); err != nil {
		return err
	}

	// 创建取消函数，用于控制自动续租协程
	newCtx, cancelFun := context.WithCancel(ctx)
	dl.cancelFun = cancelFun

	// 启动后台续租协程
	dl.startWatchDog(newCtx)
	dl.status = true

	return nil
}

// lock 尝试通过 Redis 获取锁
func (dl *DistributeLockRedis) lock(ctx context.Context) error {
	res, err := redis.String(dl.redis.Do("SET", dl.key, 1, "NX", "EX", dl.expire))
	if err != nil {
		return fmt.Errorf("redis 操作失败: %v", err)
	}
	if res != "OK" {
		return fmt.Errorf("锁已被占用")
	}
	return nil
}

// startWatchDog 创建后台守护协程，定期自动续租
func (dl *DistributeLockRedis) startWatchDog(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				// 收到取消信号，退出协程
				return
			default:
				dl.mu.Lock()
				if dl.status {
					// 尝试续租
					_, err := redis.Int(dl.redis.Do("EXPIRE", dl.key, dl.expire))
					if err != nil {
						fmt.Printf("锁续租失败: %v\n", err)
					}
				}
				dl.mu.Unlock()

				// 休眠 expire/2 秒后再续租
				time.Sleep(time.Duration(dl.expire/2) * time.Second)
			}
		}
	}()
}

// Unlock 释放锁
func (dl *DistributeLockRedis) Unlock() error {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	// 停止续租协程
	if dl.cancelFun != nil {
		dl.cancelFun()
	}

	// 释放 Redis 锁
	if dl.status {
		res, err := redis.Int(dl.redis.Do("DEL", dl.key))
		if err != nil {
			return fmt.Errorf("释放锁失败: %v", err)
		}
		if res == 1 {
			dl.status = false
			return nil
		}
	}
	return fmt.Errorf("释放锁失败")
}

// 示例代码
func main() {
	// 模拟 Redis 连接（替换为实际连接）
	redisConn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Printf("连接 Redis 失败: %v\n", err)
		return
	}
	defer redisConn.Close()

	// 创建锁实例
	lock := NewDistributeLockRedis(redisConn, "my-lock", 10)

	// 尝试加锁
	ctx := context.Background()
	if err := lock.TryLock(ctx); err != nil {
		fmt.Printf("加锁失败: %v\n", err)
		return
	}
	fmt.Println("加锁成功")

	// 模拟处理逻辑
	time.Sleep(15 * time.Second)

	// 释放锁
	if err := lock.Unlock(); err != nil {
		fmt.Printf("释放锁失败: %v\n", err)
	} else {
		fmt.Println("释放锁成功")
	}
}

```