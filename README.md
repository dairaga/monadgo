# monadgo

MonadGo is toolkit about implementing Scala monadic operations, like map, flatMap, fold, foreach, forall and etc. It is used in internal tools of [Cyberon Corporation](https://www.cyberon.com.tw/projects/cyberon_web/english/index.html.php). Refer to test codes and see how to use it.

[About Cyberon Corporation](https://www.cyberon.com.tw/projects/cyberon_web/english/index.html.php)  
[About Scala](https://www.scala-lang.org/)

## Data Types

### Unit

**Unit** represents Unit in Scala or C-like void.

[Unit in Scala](https://www.scala-lang.org/api/current/scala/Unit.html)

### Null

**Null** represents Null in Scala, or nil in Go. Container has Null inside if nil is added into it.

[Null in Scala](https://www.scala-lang.org/api/current/scala/Null.html)

### Nothing

**Nothing** represents Nothing in Scala.

[Nothing in Scala](https://www.scala-lang.org/api/current/scala/Nothing.html)


### Tuple

**Tuple** represents Tuple in Scala. More than one outputs from function or method invoked in monadgo will be converted to Tuple. The min size of monadgo Tuple is **2**.

[Tuple2 in Scala](https://www.scala-lang.org/api/current/scala/Tuple2.html)

#### Tuple2

```go
t := Tuple2Of("2", 1)
printGet(t.V1().(string))
printGet(t.V(0).(string))
printGet(t.V2().(int))
printGet(t.V(1).(int))
```

#### Tuple3

```go
t := Tuple3Of(1, "2", float64(3.14))
printGet(t.V1().(int))
printGet(t.V(0).(int))
printGet(t.V2().(string))
printGet(t.V(1).(string))
printGet(t.V3().(float64))
printGet(t.V(2).(float64))
```

#### Tuple4

```go
t := Tuple4Of(1, "2", float64(3.14), complex(5, 7))
printGet(t.V1().(int))
printGet(t.V(0).(int))
printGet(t.V2().(string))
printGet(t.V(1).(string))
printGet(t.V3().(float64))
printGet(t.V(2).(float64))
printGet(t.V4().(complex128))
printGet(t.V(3).(complex128))
```

#### TupleN

Dimension of tuple is more than 4.

```go
t := TupleOf([]interface{}{
    1, "2", 3.14, complex(5, 7), true,
})

printGet(t.V(0).(int))
printGet(t.V(1).(string))
printGet(t.V(2).(float64))
printGet(t.V(3).(complex128))
printGet(t.V(4).(bool))
printGet(t.Dimension())
```

### Pair

**Pair** represents Pair in Scala. It is from Tuple2 and used in Map(K,V). Pair in Scala is deprecated and use Tuple2 instead.

```go
p := PairOf(1, "100")

printGet(p.Key().(int))
printGet(p.Value().(string))
```

### Traversable

**Traversable** represents Traversable in Scala. **Slice** and **Map** are Traversable types.

[Traversable in Scala](https://www.scala-lang.org/api/current/scala/collection/Traversable.html)

### Slice

**Slice** wraps GO slice and implements monadic functions like in List of Scala.

```go
s := SliceOf([]int{1, 2, 3, 4, 5})

printGet(s.Get().([]int))
fmt.Println(s.Len())
fmt.Println(s.Cap())

// Foreach
s.Foreach(func(x int) {
    fmt.Println(x)
})

// sum = 16
sum := s.Fold(1, func(z, x int) int {
    return z + x
})
```

#### Map operation in Slice

```go
s1 := SliceOf([]Pair{PairOf(1, 11), PairOf(2, 22), PairOf(1, 111), PairOf(2, 222)})

s2 := s1.Map(func(p Pair) string {
    return p.String()
})

fmt.Println(s2)
printGet(s2.Get().([]string))

s2 = s1.Map(func(k, v int) (string, int) {
    return fmt.Sprintf("%d", k+v), k * v
})

fmt.Println(s2)
printGet(s2.Get().([]Pair))
```

MonadGo will convert `func(k, v int) (string, int)` to `func(Tuple2) Pair`.

#### FlatMap in Slice

```go
s1 := SliceOf([]Pair{PairOf(1, 11), PairOf(2, 22), PairOf(1, 111), PairOf(2, 222)})
printGet(s1.Get())

s2 := s1.FlatMap(func(p Pair) []int {
    return []int{p.Key().(int), p.Value().(int)}
})
printGet(s2.Get().([]int))

s2 = s1.FlatMap(func(p Pair) map[int]int {
    return map[int]int{
        p.Key().(int): p.Value().(int),
    }
})
printGet(s2.Get().([]Pair))

SliceOf([]int{1, 2, 3}).FlatMap(func(x int) {
    SliceOf([]int{1, 2, 3}).Map(func(y int) {
        fmt.Printf("%dx%d=%d\n", x, y, x*y)
    })
})
```

[List in Scala](https://www.scala-lang.org/api/current/scala/collection/immutable/List.html)

### Map

**Map** wraps Go map and implements monadic functions like in Map of Scala.

```go
m := MapOf(map[string]int{
    "a": 11,
    "b": 22,
})
fmt.Println(m)
printGet(m.Get().(map[string]int))

m = MapOf([]Pair{PairOf("a", 11), PairOf("b", 22)})
printGet(m.Get().(map[string]int))

m = MapOf(PairOf("a", 11))
printGet(m.Get())

m = MapOf(seqOf([]Pair{PairOf("a", 11), PairOf("b", 22)}))
printGet(m.Get().(map[string]int))

m.Foreach(func(k string, v int) {
    fmt.Printf("%s->%d\n", k+k, v*v)
})
```

#### Map Operation in Map

```go
s1 := MapOf(map[string][]int{
    "a": []int{11, 111},
    "b": []int{22, 222},
}).Map(func(_ string, x []int) []int {
    return x
}).Get().([][]int)

for _, x := range s1 {
    for _, y := range x {
        fmt.Println(y)
    }
}

s2 := MapOf(map[string][]Pair{
    "a": []Pair{PairOf(1, 11), PairOf(1, 111)},
    "b": []Pair{PairOf(2, 22), PairOf(2, 222)},
}).Map(func(_ string, x []Pair) []Pair {
    return x
}).Get().([][]Pair)
for _, x := range s2 {
    for _, y := range x {
        fmt.Println(y)
    }
}

s3 := MapOf(map[string]int{
    "a": 1,
    "b": 2,
}).Map(func(k string, v int) (string, int) {
    return k + k, v + v
}).Get().(map[string]int)

for k, v := range s3 {
    fmt.Println(k, v)
}

s4 := MapOf(map[string]int{
    "a": 1,
    "b": 2,
}).Map(func(k string, v int) Pair {
    return PairOf(k+k+k, v+v+v)
}).Get().(map[string]int)

for k, v := range s4 {
    fmt.Println(k, v)
}
```

#### FlatMap Operation in Map

```go
s1 := MapOf(map[string][]int{
    "a": []int{11, 111},
    "b": []int{22, 222},
}).FlatMap(func(_ string, x []int) []int {
    return x
}).Get().([]int)
for _, x := range s1 {
    fmt.Println(x)
}

s2 := MapOf(map[string][]Pair{
    "a": []Pair{PairOf(1, 11), PairOf(1, 111)},
    "b": []Pair{PairOf(2, 22), PairOf(2, 222)},
}).FlatMap(func(_ string, x []Pair) []Pair {
    return x
}).Get().(map[int]int)

for k, v := range s2 {
    fmt.Println(k, v)
}

s3 := MapOf(map[string]int{
    "a": 1,
    "b": 2,
}).FlatMap(func(k string, v int) map[string]int {
    return map[string]int{
        k:         v,
        k + k:     v + v,
        k + k + k: v + v + v,
    }
}).Get().(map[string]int)

for k, v := range s3 {
    fmt.Println(k, v)
}
```

[Map in Scala](https://www.scala-lang.org/api/current/scala/collection/Map.html)

### Option

**Option** represents Option in Scala. **Some** and **None** are subtypes of Option.

```go
```

[Option in Scala](https://www.scala-lang.org/api/current/scala/Option.html)

### Try

**Try** represents Try in Scala. **Success** and **Failure** are subtypes of Try. TryOf returns Failure if **last** arguement is **error** or **false**.

[Try in Scala](https://www.scala-lang.org/api/current/scala/util/Try.html)

### Either

**Either** represents Either in Scala. **Right** and **Left** are subtyes of Either, EitherOf returns Left if **last arguement** is **error** or **false**.

[Either in Scala](https://www.scala-lang.org/api/current/scala/util/Either.html)

#### LeftProjection

**LeftProjection** represents LeftProjection in Scala. No RightProjection in monadog because Either of monadgo also is **Right-Biased**.

[LeftProjection in Scala](https://www.scala-lang.org/api/current/scala/util/Either$$LeftProjection.html)

### PartialFunc

**PartialFunc** represents PartialFunction in Scala. It consists of Condition and Action funtions. **Condition** checks input is valid or not, and then returns result from invoking **Action** on input if input is valid.

[PartialFunction in Scala](https://www.scala-lang.org/api/current/scala/PartialFunction.html)

### Promise and Future

**Promise** and **Future** represent Promise and Future in Scala. Unlike scala throwing exceptions, assigning new result to completed Promise and Future in MonadGo will have no effect. Promise or Future can be canceled and all futures depending on it will be canceled, too.

[Promise in Scala](https://www.scala-lang.org/api/current/scala/concurrent/Promise.html)  
[Future in Scala](https://www.scala-lang.org/api/current/scala/concurrent/Future.html)
