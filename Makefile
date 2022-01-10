# >>>>>>>>>>> 自定义常量 >>>>>>>>>>>>>
# 定义项目基本信息
COMMONENVVAR      ?= GOOS=linux GOARCH=amd64
BUILDENVVAR       ?= CGO_ENABLED=0

# >>>>>>>>>>> 必须包含的命令 >>>>>>>>>
# 定义环境变量
export GOBIN  := $(CURDIR)/bin

# 构建并编译出静态可执行文件
all: linux_build

# 代码风格检查
lint:
	golangci-lint run

# 生成可执行文件
build:
	go build -o $(GOBIN)/monitor-gateway

# 交叉编译出linux下的静态可执行文件
linux_build:
	$(COMMONENVVAR) $(BUILDENVVAR) make build

# 清除所有编译生成的文件
clean:
	@rm -rf bin

.PHONY: build linux_build all test clean
