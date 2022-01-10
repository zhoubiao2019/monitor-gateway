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
