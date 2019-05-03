# monadgo

MonadGo is implements Scala monadic operations, like map, flatMap, fold, foreach, forall and etc.

## Data Types

### Unit

**Unit** represents Unit in Scala or C-like void.

### Null

**Null** represents Null in Scala, or nil in Go. Container has Null inside if nil is added into it.

### Nothing

**Nothing** represents Nothing in Scala.

## Traversable

**Traversable** represents Traversable in Scala.

### Tuple

**Tuple** represents Tuple in Scala. More than one outputs from function or method invoked in monadgo will be converted to Tuple. The min size of monadgo Tuple is **2**.

### Pair

**Pair** represents Pair in Scala. It is from Tuple2 and used in Map(K,V).

### Option

**Option** represents Option in Scala. **Some** and **None** are subtypes of Option.

## Try

**Try** represents Try in Scala. **Success** and **Failure** are subtypes of Try. TryOf or TryxOf returns Failure if last arguement is error or false because of no exception in GO.

## Slice

**Slice** wraps GO slice and implements monad functions like in List of Scala.

## Map

**Map** wraps Go map and implements monad functions like in Map of Scala.