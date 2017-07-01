;;;; sputter core: concurrency

(defmacro async
  {:doc-asset "async"}
  [& forms]
  (list 'sputter:make-closure
    (vector)
    (cons 'sputter:do-async forms)))

(defmacro generate
  {:doc-asset "generate"}
  [& forms]
  (list 'sputter:let
    (vector 'sputter/ch (list 'sputter:channel)
            'sputter/cl (list :close 'sputter/ch)
            'emit (list :emit 'sputter/ch))
    (list 'sputter:async
      (list 'sputter:let (vector 'x (cons 'sputter:do forms))
        (list 'sputter/cl)
        'x))
    (list :seq 'sputter/ch)))

(defmacro future
  {:doc-asset "future"}
  [& forms]
  `(let [promise# (promise)]
    (async (promise# (do ~@forms)))
    promise#))
