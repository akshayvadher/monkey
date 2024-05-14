# Monkey programming language

Implementation of the Monkey programming language in Go

As described in the book [Writing An Interpreter In Go](https://interpreterbook.com/)

# Run

`go build .`

`./monkey`

# Example
```
let five = 5;
let a = true;

let add = fn(x, y) {
    x + y;
};

let result = add(five, 10);

let multipleExecution = fn(func, x) {
    func(func(x));
};
let r = multipleExecution(add, 10);
puts(r);

let a = ["a", "b", "c", 1, true];
a[0]

let h = {"a": 1, "b": 2, 3: "c"};
a["a"]

```
