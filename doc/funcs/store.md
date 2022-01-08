### 数据库操作相关

当天数据相关方法 rowKey = thedate+key 有效期一天
----------------------------------
#### 去重LIST增加元素
```
StoreAppendList(rowKey,item string) int
```

#### 去重LIST取出所有元素(output输出专用)
```
StoreGetList(rowKey string) []string
```

#### 去重LIST取出所有元素(用delimit隔开)
```
StoreGetListDelimit(rowKey,delimit string) string
```

#### 去重LIST取出所有元素个数
```
StoreGetListLen(rowKey string) int
```

#### 设置某个KEY的值
```
StoreSetValue(rowKey,value string) bool
```

#### 初始化某个key的值（不覆盖已存在值）
```
StoreSetNxValue(rowKey,value string) bool
```

#### 原子加指定值
```
StoreAddValue(rowKey string,value int) int
```
#### 取某个值
```
StoreGetValue(rowKey string)string
```

指定有效期数据相关方法(大于1天有效期，需注意key重复问题)  expire格式 123s 123秒，10d 10天
-----------------------------------------------

#### 去重LIST增加元素
```
StoreAppendListExp(exp,key,item string) int
```

#### 去重LIST取出所有元素(output输出专用)
```
StoreGetListLenExp(key string) []string
```

#### 去重LIST取出所有元素(用delimit隔开)
```
StoreGetListDelimitExp(key,delimit string) string
```

#### 去重LIST取出所有元素个数
```
StoreGetListLenExp(key string) int
```

#### 设置某个KEY的值
```
StoreSetValueExp(exp,key,value string) bool
```

#### 初始化某个key的值（不覆盖已存在值）
```
StoreSetNxValueExp(exp,key,value string) bool
```

#### 原子加指定值
```
StoreAddValueExp(exp, key string,value int) int
```

#### 取某个值
```
StoreGetValueExp(key string)string
```

多天数据相关方法 存储时key会加上从1970-01-01到现在的天数
--------------------------------------------------
#### 去重LIST增加元素
```
StoreAppendListDays(days int,key,item string)int
```
StoreGetListDays

#### 去重LIST取出所有元素(output输出专用)
```
StoreGetListDelimitDays(days int,key string)[]string
```

#### 去重LIST取出所有元素(用delimit隔开)
```
StoreGetListLenDays(days int,key,delimit string) string
```
#### 原子加指定值
```
StoreAddValueDays(days int,key string, value int) int
```
#### 取多天原子加的值合并
```
StoreGetValueDays(days int,key string) int
```