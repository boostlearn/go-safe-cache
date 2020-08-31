## 简介
为优化对热点数据的缓存，本项目设计两个LRU队列，一个LRU队列存储数据，另一个LRU队列统计数据键热度变化。

![avatar](https://github.com/boostlearn/go-safe-cache/raw/master/doc/safe_cache.png)


## 测试
本次测试，针对符合帕累托分布(二八原则)数据进行，生成方法如下：

    import numpy as np
    np.random.pareto(4, 1000000))*10000 + 1
 
数据分布如下：
 ![avatar](https://github.com/boostlearn/go-safe-cache/raw/master/doc/pareto_4.png)
 
测试结果：
### 缓存大小100时，KEY命中率
|最小热度K值|LRU算法|2Q算法|ARC算法|
|:----|----:|----:|----:|
|0(基准)|	0.0875463|	0.118036|	0.1210772|
|1|	0.1283985|	0.1350703|	0.1397713|
|2|	0.1285002|	0.1351739|	0.1398301|
|4|	0.1595406|	0.1438613|	0.1510483|
|8|	-|	0.1617143|	0.1640282|


### 缓存大小500时，KEY命中率

|最小热度K值|LRU算法|2Q算法|ARC算法|
|:----|----:|----:|----:|
|0(基准)|	0.3960788|	0.4558791|	0.4760088|
|1|	0.4819089|	0.4843694|	0.5070587|
|2|	0.4819796|	0.4841934|	0.5071921|
|4|	0.554208|	0.4959567|	0.52678|
|8|	0.568441|	0.5585039|	0.5572612|
|16|		|0.5664294|	0.5685396|


### 缓存大小1000时，KEY命中率

|最小热度K值|LRU算法|2Q算法|ARC算法|
|:----|----:|----:|----:|
|0(基准)|	0.7231351|	0.693408|	0.7209555|
|1|	0.7227088|	0.7102965|	0.7432944|
|2|	0.7691043|	0.7101937|	0.7432019|
|4|	0.7855555|	0.7160656|	0.7601185|
|8|	-|	0.7740888|	0.7720391|
|16|	-|	0.7813305|	0.7814856|


### 缓存大小2000时，KEY命中率

|最小热度K值|LRU算法|2Q算法|ARC算法|
|:----|----:|----:|----:|
|0(基准)|	0.8959815|	0.8899972|	0.9080057|
|1|	0.906264|	0.8898667|	0.9145529|
|2|	0.9063081|	0.8897457|	0.9143771|
|4|	0.918635|	0.8899214|	0.9170311|
|8|	0.9227886|	0.9190385|	0.9192692|


## 算法流程
访问流程：
![avatar](https://github.com/boostlearn/go-safe-cache/raw/master/doc/safe_cache_query.png)

更新流程
![avatar](https://github.com/boostlearn/go-safe-cache/raw/master/doc/safe_cache_insert.png)

## BenchMark
    BenchmarkBucketLru_Single-4      1000000              1268 ns/op             140 B/op          4 allocs/op
    BenchmarkBucketLru_K-4           1000000              1257 ns/op             118 B/op          3 allocs/op
    BenchmarkBucket2Q_Single-4       1000000              1971 ns/op             166 B/op          4 allocs/op
    BenchmarkBucket2Q_K-4            1000000              1624 ns/op             112 B/op          3 allocs/op
    BenchmarkBucketArc_Single-4      1000000              2028 ns/op             167 B/op          4 allocs/op
    BenchmarkBucketArc_K-4           1000000              1692 ns/op             113 B/op          3 allocs/op

## 示例