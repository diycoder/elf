# 运行环境

## 环境变量
默认为`RUN_ENV`，也可通过`SetRunEnvKey`设置运行环境的环境变量名。
#### example
```
SetRunEnvKey("RUN_ENV")
```

## 运行环境配置
运行环境由 `集群类型`+`环境类型` 两块组合而成，`集群类型` 为可选，`环境类型` 为必选，两者都采用`小写`并使用`下划线`分割。

`集群类型` 不做任何限制可自行定义，`环境类型` 除 `develop` `test` `gray` `product` 四个默认环境类型外也可自行定义。

#### example
```
// 本地开发
RUN_ENV = develop

// 虚机测试
RUN_ENV = vm_test

// tke灰度
RUN_ENV = tke_gray

// tke金丝雀
RUN_ENV = tke_canary
```

## 判断运行环境
提供四个默认运行环境判断方法 `IsDev` `IsTest` `IsGray` `IsProd`（后缀匹配），还提供了两个自定义运行环境判断方法`Is` `Not`（后缀匹配）。

#### example
```
// 本地开发
RUN_ENV = develop
IsDev() // true

// 虚机开发
RUN_ENV = vm_develop
IsDev() // true

// 本地开发
RUN_ENV = develop
IsTest() // false

// tke开发
RUN_ENV = tke_develop
IsTest() // false

// 金丝雀
RUN_ENV = canary
Is("canary") // true

// 虚机金丝雀
RUN_ENV = vm_canary
Is("canary") // true

// 虚机开发
RUN_ENV = vm_develop
Is("canary") // false

// 金丝雀
RUN_ENV = canary
Not("canary") // false
```