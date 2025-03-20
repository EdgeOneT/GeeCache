package singleflight

import "sync"

type call struct {
	wg  sync.WaitGroup // 用于等待并发请求完成
	val interface{}    // 存储调用的结果
	err error          // 存储调用中的错误
}

type Group struct {
	mu sync.Mutex       // 保护 m 的并发安全
	m  map[string]*call // 记录每个 key 对应的 call 对象
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call) // 延迟初始化，节省内存
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()       // 解锁，允许其他协程进入
		c.wg.Wait()         // 等待已有请求完成
		return c.val, c.err // 直接共享结果
	}
	c := new(call)
	c.wg.Add(1)   // 标记请求开始
	g.m[key] = c  // 注册到 Group
	g.mu.Unlock() // 释放锁，允许其他协程处理其他 key

	c.val, c.err = fn() // 执行实际函数（如查询数据库）
	c.wg.Done()         // 标记请求完成，唤醒等待的协程

	g.mu.Lock()
	delete(g.m, key) // 清理已完成的请求
	g.mu.Unlock()

	return c.val, c.err
}
