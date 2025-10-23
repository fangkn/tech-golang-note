
# 单例说明

创建100个协程，每个协程都会调用  GetInstance() 函数，该函数会创建一个单例对象，并且返回该对象。

只会创建一次。

运行结果如下：

```
=== RUN   TestSingleton
Creating Singleton instance
--- PASS: TestSingleton (0.00s)
PASS
ok      singleton       1.243s
```