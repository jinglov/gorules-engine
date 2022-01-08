### 数据转换相关

#### 字符串的长度 Len(args string) int
```
    Len(args string) int
```

#### 求和 Sum
```
    Sum(args ...int) int
```

#### 减 Subtract
```
    Subtract(a, b int) int
```

#### 乘法 Multipy
```
    Multiply(args ...int)int
```

#### 除法 Div（只取商）

```
    Div(a,b int)int
```

注意：返回的是 toString(int),舍去小数部分。需要小数部份时请把分母乘一定系数后再除
示例:  a/b > 0.5
```
    gt(div(multiply(a,10,b), 5)
```

### 字符串连接 Concat
```
    Concat(args ...string) string
```

#### 用一个符号隔开，取其中一个 split() string
```
    Split(str,delimiter string, index int) string
    split("a,b,c", ",", "0") -> "a"
```

#### 全转大写 Upper
 ```
    Upper(str string) string
```

#### 全转小写 Lower
```
    Lower(str string) string 
```

#### 取左第N个字符 Left() string
```
Left(k string, len int) string
Left("abc", 1) -> "c"
Left("abc", 4) -> "abc"
Left("abc", -1) -> "ab"
```

#### 取右边第N个字符 Right() string
```
Right(k string, len int) string
Right("abc", 1) -> "c"
Right("abc", -1) -> "bc"
```

#### 时间格式化成时间戳 dttounix string 传空时取当前时间戳
```
DtToUnix(a string) int 
```