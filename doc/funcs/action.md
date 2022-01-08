
###条件判断相关

#### 并 And
```
    And(args ...string) string
```

#### 或 Or
```
    Or(args ...string) string
```

#### 相等 Eq
```
    Eq(args ...string) string
```

#### 不相等 Neq
```
    Neq(args ...string) string
```

#### 大于 Gt
```
    Gt(v1, v2 string) string
```

#### 大于等于 Gte
```
    Gte(v1, v2 string) string
```

#### 小于 Lt
```
    Lt(v1, v2 string) string
```

#### 小于等于 Lte
```
    Lte(v1, v2 string) string
```

#### 相似 Like
```
    Like(v1, v2 string) string
```

#### 相似其中一个 LikeOr
```
    LikeOr(args ...string) string
```

#### 在之中 In
```
    In(args ...string) string
```

#### 不在之中 NotIn
```
    NotIn(args ...string) string
```

### 在之中，第二个参数是一个以第三个参数分隔的数组
```
    SplitIn(item, str,sep string)string
```


###  版本比较（仅比较前3位，如：2.8.17.2019 只有2.8.17参与比较）

#### 版本小于 LtVersion
```
    LtVersion(v1, v2 string)string
```

#### 版本等于 EqVersion
```
    EqVersion(v1, v2 string)string
```

#### 版本大于 GtVersion
```
    GtVersion(v1, v2 string)string
```
