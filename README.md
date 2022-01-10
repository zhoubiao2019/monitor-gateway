## monitor-gateway Service


### 部署说明

## 编译

- 编译文件
```
make build //本地编译
make linux_build //服务端编译
```

## 初始化数据库
```
create database monitor;
```

## 运行
- 启动服务
```
bin/monitor-gateway -f conf/config.yaml
```

### 依赖的基础服务
- mysql

### 依赖的上游服务
暂无

### Design Considerations
1，对于这个题目更加合理的设计应该是不让监控服务耦合具体的业务，而是在被监控的任务上进行埋点(心跳检测，任务执行成功次数，执行失败次数)，将所有的 metrics 采集到监控中心(这里的场景更加适合 push 的方式)，然后单独在监控中心对于具体的指标做一些聚合的计算并设置一些预警规则，最后就是选择如何推送通知了
