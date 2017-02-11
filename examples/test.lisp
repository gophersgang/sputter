(defvar blah "blah is stored")

(defun say-hello [name]
  (println "Hello there, " name "!"))

(let (a 10)
  (println "a = " a)
  (let (a 20)
    (println "a = " a))
  (println "a = " a))

(println
  "hello "
  (*
    (+ 9.5e+2 1.0)
    2.0)
  " "
  '(blah 99 "yep"))

(say-hello "Thom")

(println
  "howdy "
  "ho "
  blah)

(if (list? [1])
  (println "yep")
  (println "nope"))

(println [1 2 3])
