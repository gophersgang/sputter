Sputter is a simple Lisp Interpreter written in
[Go](https://golang.org/). Basically, it's just me having some fun
and trying to improve my Go skills. That means you're unlikely to
find something you'd want to use in production here. On the other
hand, if you want to join in on the fun, you're more than welcome
to.

## How To Install
Make sure your `GOPATH` is set, then run `go get` to retrieve the
package.

```bash
go get github.com/kode4food/sputter
```

## How To Invoke The Interpreter
Once you've installed the package, you can run it from `GOPATH/bin`
like so:

```bash
sputter somefile.lisp

# or

cat somefile.lisp | sputter
```

## How To Invoke The REPL
Sputter has a very crude Read-Eval-Print Loop that will be more than
happy to start if you call it with no arguments from the terminal:

<img src="img/repl.jpeg" />

## Current Status
I just started this thing and it's still pretty fragile, but
that will change rapidly. The current built-in forms are:

  * Control and Branching: `if`, `cond`, `quote`, `do`
  * Numeric: `+`, `-`, `*`, `/`, `!=`, `=`, `<`, `<=`, `>`, `>=`
  * Variables: `def`, `let`, `ns`, `with-ns`
  * Functions: `defn`, `fn`, `apply`
  * Predicates: `eq`, `!eq`, `nil?`, `!nil?`
  * Sequences: `cons`, `conj`, `first`, `rest`, `seq?`, `!seq?`
  * Lists: `list`, `list?`, `!list?`, `to-list`
  * Vectors: `vector`, `vector?`, `!vector?`, `to-vector`
  * Associative Arrays: `assoc`, `assoc?`, `!assoc?`, `to-assoc`
  * Indexed Sequences: `nth`, `indexed?`, `!indexed?`
  * Mapped Sequences: `mapped?`, `!mapped?`
  * Metadata: `meta`, `with-meta`
  * Comprehensions: `concat`, `map`, `filter`
  * Concurrency: `channel`, `async`, `promise`, `future`
  * I/O: `print`, `println`, `pr`, `prn`

Documentation for most of these forms may be viewed in the
REPL using the `doc` function.