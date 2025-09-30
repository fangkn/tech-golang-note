

# 并发基础

在 Go 语言中最重要的一个特性，那就是 go 关键字。

并发的实现有哪些：

- **多进程**。多进程是在操作系统层面进行并发的基本模式。
- **多线程**。多线程在大部分操作系统上都属于系统层面的并发模式。
- **基于回调的非阻塞/异步IO**。在很多 高并发服务器开发实践中，使用多线程模式会很快耗尽服务器的内存和CPU资源。
- **协程**。本质上是一种用户态线程，不需要操作系统来进行抢占式调度， 且在真正的实现中寄存于线程中，因此，系统开销极小，可以有效提高线程的任务并发性，而避免多线程的缺点。

我们习惯于串行模式，串行方式有确定性，但是，为了提高处理效率，我们引入了并发模式。 并发模式会存在不确定性，导致处理行为的意外发生，所以我们采用到共享内存的方式。为了保证共享内存的有效性，我们采取了很多措施，比如加锁等，来避免死锁或资源竞争。实践证明，我们很难面面俱到，往往会在工程中遇到各种奇怪的故障和问题。

所以，又研究另一种新的系统模型，称之为“**消息传递系统**”。



# 协程

协程，是一种轻量级线程。之后叫“轻量级”是因为它可以轻松创建上百万个而不会导致系统资源衰竭。而线程和进程通常最多也不能超过1万个。

Go 语言在语言级别支持轻量级线程，叫 `goroutine`。

Go 语言标准库提供的所有系统调用操作 (当然也包括所有同步 IO 操作)，都会出让 CPU 给其他 `goroutine`。这让事情变得非常简单，让轻
量级线程的切换管理不依赖于系统的线程和进程，也不依赖于CPU的核心数量。

# goroutine

`goroutine` 是 Go 语言中的轻量级线程实现，由Go运行时(`runtime`)管理。

在一个函数调用前加上 `go` 关键字，这次调用就会在一个新的 `goroutine` 中并发执行。

如：

```go 
package main

import "fmt"

func Add(x, y int) { z := x + y
        fmt.Println(z)
}

func main() {
    for i := 0; i < 10; i++ {
        go Add(i, i) 
    }
}
```

Go 程序从初始化 main package 并执行 main() 函数开始，当 main() 函数返回时，程序退出， 且程序并不等待其他`goroutine`(非主`goroutine`)结束。把这个例子什么也不会输出。



# 并发通信

从上面的例子，可以看到我们需要处理的是并发之间的通信。 

有两种最常见的并发通信模型:共享数据和消息。

共享数据是指多个并发单元分别保持对同一个数据的引用，实现对该数据的共享。被共享的数据可能有多种形式，比如内存数据块、磁盘文件、网络数据等。

```go
package main
import "fmt"
import "sync"
import "runtime"

var counter int = 0

func Count(lock *sync.Mutex) { lock.Lock()
    counter++
    fmt.Println(z)
    lock.Unlock()
}

func main() {
    lock := &sync.Mutex{}

    for i := 0; i < 10; i++ { 
        go Count(lock)
    }

    for { 
        lock.Lock()
        c := counter
        lock.Unlock()
        runtime.Gosched() 
        if c >= 10 {
            break
        } 
    }
}
```

这种方式是引进了锁的方式来完成通信，这个试比较晦涩，如C++ 相同的方式， 并没有什么太大的优势。 

Go 语言提供的是另一种以消息机制的通信模型，被称为 `channel`。

# channel

`channel` 是 Go 语言在语言级别提供的 `goroutine` 间的通信方式。 我们可以使用 `channel` 在两个或多个`goroutine`之间传递消息。

`channel`  是**进程内**的通信方式，因此通过 `channel `传递对象的过程和调用函数时的参数传递行为比较一致，比如也可以传递指针等。

如果需要跨进程通信，我们建议用分布式系统的方法来解决，比如使用 `Socket` 或者 HTTP 等通信协议。Go 语言对于网络方面也有非常完善的支持。

我们可以用 `channel` 重写上例子。

```go 
package main
import "fmt"

func Count(ch chan int) { 
    ch <- 1
    fmt.Println("Counting")
}

func main() {

    chs := make([]chan int， 10)  
    for i := 0; i < 10; i++ {
        chs[i] = make(chan int) // 没有缓冲区
        go Count(chs[i]) 
    }

    for _, ch := range(chs) { 
        <-ch
    }
}
```
定义了一个包含10个 channel 的数组(名为`chs`)，并把数组中的每个 channel 分配给 10 个不同的 `goroutine`。

在每个 `goroutine` 的 Add() 函数完成后，我们通过`ch <- 1`语句向对应的 `channel` 中写入一个数据。在这个 `channel` 被读取前，这个操作是阻塞的。在所有的 `goroutine` 启动完成后，我们通过 `<-ch` 语句从10个 `channel` 中依次读取数据。

## 基本语法

channel 的声明形式为:

```go
var chanName chan ElementType
```

`ElementType` 指定这个`channel` 所能传递的元素类型。如：

```go 
var ch chan int
var m map[string] chan bool
```

定义一个 channel 也很简单，直接使用内置的函数 `make()`即可

```go
ch := make(chan int)
```

**在 channel 的用法中，最常见的包括写入和读出。**

写入，将一个数据写入(发送)至 channel：
```go
ch <- value
```
向 channel 写入数据通常会导致程序阻塞。直到有其他 goroutine 从这个channel 中读取数据。

```go
value := <-ch
```
如果 channel 之前没有写入数据，那么从 channel 中读取数据也会导致程序阻塞，直到 channel 中被写入数据为止。

## select

go 语言直接在语言级别支持 select 关键字，用于处理异步 IO 问题。

select 的用法与 switch 语言非常类似，由 select 开始一个新的选择块，每个选择条件由 case 语句来描述。

最大的一条限制就是每个 case 语句里必须是一个IO操作。

```go
select {
    case <-chan1:
    // 如果chan1成功读到数据，则进行该case处理语句 
    case chan2 <- 1:
    // 如果成功向chan2写入数据，则进行该case处理语句 
    default:
    // 如果上面都没有成功，则进入default处理流程 
}
```



`select` **默认是阻塞的**，只有当监听的 `channel` 中有发送或接收到可以进行时才会运行。多个 `channel` 都准备好时，`select` 是随机的选择一个执行。 

`select`  中有 `default` 语法，`default` 就是当监听 `channel` 都没有准备好时时候， 默认执行（`select` 不再阻塞等待 `channel`）

```go 
select {
  case i :=<-c
  // use i 
  default : 
  // 当 C阻塞的时候执行。
}
```



## 缓冲机制

当我们需要持续传输大量数据时，`channel` 就要带上缓存。

创建一个带缓冲的 `channel`:

```go
 c := make(chan int, 1024)  // 带缓冲的
```
在调用 `make()` 时将缓冲区大小作为第二个参数传入即可，比如上面这个例子就创建了一个大小 为 1024 的 int 类型 channel，即使没有读取方，写入方也可以一直往 channel 里写入，**在缓冲区被填完之前都不会阻塞**。

## 超时机制

在并发编程的通信过程中，最需要处理的就是超时问题。即向 `channel` 写数据时发现 `channel`已满，或者从 `channel` 试图读取数据时发现 `channel` 为空。如果不正确处理这些情况，很可能会导致整个`goroutine`锁死。

Go 语言没有提供直接的超时处理机制，但我们可以利用 `select` 机制。`select` 的特点是只要其中一个 `case` 已经完成，程序就会继续往下执行，而不会考虑其他 `case` 的情况。    

```go
// 首先，我们实现并执行一个匿名的超时等待函数 
timeout := make(chan bool, 1)
go func() {
    time.Sleep(1e9) // 等待1秒钟
    timeout <- true
}()

// 然后我们把timeout这个channel利用起来
select {
    case <-ch:
    // 从ch中读取到数据
    case <-timeout:
    // 一直没有从ch中读取到数据，但从timeout中读取到了数据
}
```

这样使用 `select` 机制可以避免永久等待的问题，因为程序会在 `timeout` 中获取到一个数据 后继续执行，无论对 ch 的读取是否还处于等待状态，从而达成1秒超时的效果。

## channel 的传递

在 Go 语言中 channel 本身也是一个原生类型，与 map类型地位一样，因此 channel 本身在定义后也可以通过 channel 来传递。

## 单向 channel

顾名思义，单向channel只能用于发送或者接收数据。

单向channel变量的声明非常简单，如下:

```go
var ch1 chan int        // ch1是一个正常的channel，不是单向的 
var ch2 chan<- float64  // ch2是单向channel，只用于写float64数据
var ch3 <-chan int      // ch3是单向channel，只用于读取int数据
```
基于 `ch4`，我们通过类型转换初始化了两个单向 `channel` 单向读的 `ch5` 和单向写的 `ch6`。

单向 `channel` 也是起到这样的一种契约作用。

关闭 `channel`

关闭`channel` 非常简单，直接使用Go语言内置的`close()`函数即可。

```go
close(ch)
```

如何判断一个 `channel` 是否已经被关 闭?我们可以在读取的时候使用多重返回值的方式:

```go 
 x, ok := <-ch
```
这个用法与 map 中的按键获取 value 的过程比较类似，只需要看第二个 `bool` 返回值即可，如果返回值是false则表示ch已经被关闭。

# runtime 

`runtime` 包中有几个处理`goroutine`的函数：

- **Goexit** : 退出当前执行的 goroutine，但是 defer 函数还会继续调用
- **Gosche** : 让出当前 goroutine 的执行权限，调度器安排其他等待的任务运行，并在下次某个时候从该位置恢复执行。
- **NumCPU** : 返回 CPU 核数量。
- **NumGoroutine** : 返回正在执行和排队的任务总数。
- **GOMAXPROCS** : 用来设置可以运行的CPU核数。

# 同步

Go 语言的设计者虽然对`channel`有极高的期望，但也提供了妥善的资源锁方案。

## 同步锁

Go语言包中的 sync 包提供了两种锁类型: `sync.Mutex` 和 `sync.RWMutex`。`Mutex `是最简单 的一种锁类型，

## 全局唯一性操作

于从全局的角度只需要运行一次的代码，比如全局初始化操作，Go语言提供了一个 Once 类型来保证全局的唯一性操作。

**sync.Once** 是 Golang package 中使方法只执行一次的对象实现，作用与 **init** 函数类似。但也有所不同。

-   **init** 函数是在文件包首次被加载的时候执行，且只执行一次
-   **sync.Once** 是在代码运行中需要的时候执行，且只执行一次。所以，单例就可以用这种方式

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once
	
	onceBody := func() {
		fmt.Println("Only once")
	}
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			once.Do(onceBody)
			done <- true
		}()
	}
	for i := 0; i < 10; i++ {
		<-done
	}
}
```

输出：

```sh 
# Output:
Only once
```

**sync.Once** 使用变量 `done` 来记录函数的执行状态，使用 `sync.Mutex` 和 `sync.atomic` 来保证线程安全的读取 `done` 。

例子：[golang-code/singletion](https://github.com/4kkane/4kane-code/tree/main/golang-code/singleton)

 

**源码**

```go
package sync

import (
	"sync/atomic"
)

// Once is an object that will perform exactly one action.
type Once struct {
	m    Mutex
	done uint32
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once. Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	config.once.Do(func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
//
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}·比·
}
```

# 原子操作



golang 原子操作指的 是  golang 中 sync.atomic 的常见操作。

atomic 提供的原子操作能够确保任一时刻只有一个goroutine对变量进行操作，可以减少用锁操作。

atomic常见操作有：

- 增减
- 载入
- 比较并交换
- 交换
- 存储

## 增减

```go
func AddInt32(addr *int32, delta int32) (new int32)
func AddInt64(addr *int64, delta int64) (new int64)
func AddUint32(addr *uint32, delta uint32) (new uint32)
func AddUint64(addr *uint64, delta uint64) (new uint64)
func AddUintptr(addr *uintptr, delta uintptr) (new uintptr)
```

第一个参数必须是指针类型的值，通过指针变量可以获取被操作数在内存中的地址，从而施加特殊的CPU指令，确保同一时间只有一个goroutine能够进行操作。



## 载入

```go
func LoadInt32(addr *int32) (val int32)
func LoadInt64(addr *int64) (val int64)
func LoadPointer(addr *unsafe.Pointer) (val unsafe.Pointer)
func LoadUint32(addr *uint32) (val uint32)
func LoadUint64(addr *uint64) (val uint64)
func LoadUintptr(addr *uintptr) (val uintptr)
```

载入操作能够保证原子的读变量的值，当读取的时候，任何其他CPU操作都无法对该变量进行读写，其实现机制受到底层硬件的支持。见上述例子中的`atomic.LoadInt64(&opts)`。

## 比较并交换

该操作简称 `CAS` (Compare And Swap)。

```go
func CompareAndSwapInt32(addr *int32, old, new int32) (swapped bool)
func CompareAndSwapInt64(addr *int64, old, new int64) (swapped bool)
func CompareAndSwapPointer(addr *unsafe.Pointer, old, new unsafe.Pointer) (swapped bool)
func CompareAndSwapUint32(addr *uint32, old, new uint32) (swapped bool)
func CompareAndSwapUint64(addr *uint64, old, new uint64) (swapped bool)
func CompareAndSwapUintptr(addr *uintptr, old, new uintptr) (swapped bool)
```

该操作在进行交换前首先确保变量的值未被更改，即仍然保持参数 `old` 所记录的值，满足此前提下才进行交换操作。`CAS`的做法类似操作数据库时常见的乐观锁机制。

需要注意的是，当有大量的`goroutine` 对变量进行读写操作时，可能导致`CAS`操作无法成功，这时可以利用for循环多次尝试。

```go
var value int64

func atomicAddOp(tmp int64) {
for {
       oldValue := value
       if atomic.CompareAndSwapInt64(&value, oldValue, oldValue+tmp) {
           return
       }
   }
}
```

## 交换



```go
func SwapInt32(addr *int32, new int32) (old int32)
func SwapInt64(addr *int64, new int64) (old int64)
func SwapPointer(addr *unsafe.Pointer, new unsafe.Pointer) (old unsafe.Pointer)
func SwapUint32(addr *uint32, new uint32) (old uint32)
func SwapUint64(addr *uint64, new uint64) (old uint64)
func SwapUintptr(addr *uintptr, new uintptr) (old uintptr)
```

相对于`CAS`，明显此类操作更为暴力直接，并不管变量的旧值是否被改变，直接赋予新值然后返回背替换的值。

## 存储



```go
func StoreInt32(addr *int32, val int32)
func StoreInt64(addr *int64, val int64)
func StorePointer(addr *unsafe.Pointer, val unsafe.Pointer)
func StoreUint32(addr *uint32, val uint32)
func StoreUint64(addr *uint64, val uint64)
func StoreUintptr(addr *uintptr, val uintptr)
```

此类操作确保了写变量的原子性，避免其他操作读到了修改变量过程中的脏数据。
