# monadgo

MonadGo is implements Scala monadic operations, like map, flatMap, fold, foreach, forall and etc.

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

## Traversable

**Traversable** represents Traversable in Scala.

[Traversable in Scala](https://www.scala-lang.org/api/current/scala/collection/Traversable.html)

### Tuple

**Tuple** represents Tuple in Scala. More than one outputs from function or method invoked in monadgo will be converted to Tuple. The min size of monadgo Tuple is **2**.

[Tuple2 in Scala](https://www.scala-lang.org/api/current/scala/Tuple2.html)

### Pair

**Pair** represents Pair in Scala. It is from Tuple2 and used in Map(K,V). Pair in Scala is deprecated and use Tuple2 instead.

[Try in Scala](https://www.scala-lang.org/api/current/scala/util/Try.html)

### Slice

**Slice** wraps GO slice and implements monad functions like in List of Scala.

[List in Scala](https://www.scala-lang.org/api/current/scala/collection/immutable/List.html)

### Map

**Map** wraps Go map and implements monad functions like in Map of Scala.

[Map in Scala](https://www.scala-lang.org/api/current/scala/collection/Map.html)

### Option

**Option** represents Option in Scala. **Some** and **None** are subtype and value of Option.

[Option in Scala](https://www.scala-lang.org/api/current/scala/Option.html)

### Try

**Try** represents Try in Scala. **Success** and **Failure** are subtypes of Try. TryOf or TryxOf returns Failure if last arguement is error or false.

### Either

**Either** represents Either in Scala. **Right** and **Left** are subtyes of Either, EitherOf or EitherxOf returns Left if last arguement is error or false.

[Either in Scala](https://www.scala-lang.org/api/current/scala/util/Either.html)

#### LeftProjection

**LeftProjection** represents LeftProjection in Scala. No RightProjection in monadog because Either of monadgo also is **Right-Biased**.

[LeftProjection in Scala](https://www.scala-lang.org/api/current/scala/util/Either$$LeftProjection.html)

### PartialFunc

**PartialFunc** represents PartialFunction in Scala. It consists of Condition and Action funtions. **Condition** checks input is valid or not, and then returns result from invoking **Action** on input if input is valid.

[PartialFunction in Scala](https://www.scala-lang.org/api/current/scala/PartialFunction.html)