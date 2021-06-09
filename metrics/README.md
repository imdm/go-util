# metrics

对 prometheus 的封装，依赖环境变量 "SERVICE\_NAME"。  
### 简介
对外暴露四种 metrics: Counter, Gauge, Summary, Histogram.  
Counter 用于记录某指标的总数量，对时间求导可得到每秒的增量，通常用于统计接口的 qps。  
Gauge 用于记录某指标的值，每次调用会覆盖旧值，用于监控如机器内存等信息。  
Summary 用于统计一组值的百分位数，如接口响应时延的 pct\_50（接口响应时延的中位数，注意不是平均值）, pct\_99 (99% 的请求时延低于哪个值)。summary 指标可以帮助我们有目标的优化接口响应时间、优化服务稳定性。  
Histogram 作用与 Summary 类似，但底层实现不一样。Histogram 计算百分位数是在 prometheus server 端完成，Summary 则是在客户端完成。更详细的对比可看[这篇文档](https://prometheus.io/docs/practices/histograms/)。  
### example
```
reqCount = metrics.NewCounter("http_request_total", []string{"view", "status", "err_code"})
reqCount.Add(1, map[string]string{"view": "Login", "err_code": "0"})

