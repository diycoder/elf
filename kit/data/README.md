### echo 返回结果封装

### 使用说明

若报错，json返回结果调用如下

```bazaar
data.Error(ctx, err)
```
若不报错，json调用如下

```bazaar
data.Success(ctx, result)
```

若返回文件，调用如下

```bazaar
data.File(ctx, f)
```