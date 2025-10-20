

你让AI 帮你写代码，然后再用同一个 AI 帮你 Review 代码。 最后注定就会坑你的。
所以，最好还是自己 Review 代码。

问题如下： 


```go 
    applyTime = time.Unix(0, 0)
    if applyTime.IsZero() {
        //todo something
    }
```
由于历史原因， applyTime 是 `1970-01-01 08:00:00`， 是用了`time.Unix(0, 0)` 而不是 `time.Time{}`。

我的目的是： 当到 applyTime = 1970-01-01 08:00:00 执行 下面  todo something 的逻辑 。 

AI 给我写的代码是  `if applyTime.IsZero() {}`  IsZero的用法， 我有点怀疑。所以我就还单独的问了一下AI， IsZero() 是否可以判断出

`1970-01-01 08:00:00`。它给我很自信，很能肯定的回答。 我没有多想，相信了它。

最后，就出了问题。 

IsZero() 的实现方法就是判断 sec 和 nsec 是否都为 0。 这种情况下 `1970-01-01 08:00:00` 肯定就是一个非0的值了。 我就是点进去看一下它的实现，也能发现端倪。 可以是由于有AI 的辅助，就懒了，没有进一步的验证。

```golang 
// IsZero reports whether t represents the zero time instant,
// January 1, year 1, 00:00:00 UTC.
func (t Time) IsZero() bool {
	return t.sec() == 0 && t.nsec() == 0
}
```

现在AI 用了越来越多。 然后人有越想偷懒了，然后，问题就来了。 

见代码： 

[](source/time-iszero/main.go)









