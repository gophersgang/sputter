# (let [<name form>+] form+) assigns local variables
Will create a new local scope, assigning the specified variables to that scope. It will then evaluate the specified forms within that scope and return the result of the last evaluation.

## An Example

  (let [x '(1 2 3 4)
        y [5 6 7 8]]
    (concat x y))

This example will create a list called 'x' and a vector called 'y' and return the lazy concatenation of those sequences. Note that the two variables do not exist outside of the assignment form.
