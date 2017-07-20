;;;; sputter core: i/o

(defn pr-map-with-nil
  {:private true}
  [func seq]
  (map (lambda [val]
    (if (nil? val) val (func val))) seq))

(defn pr [& forms]
  (let [s (pr-map-with-nil str! forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write *space* e))))

(defn prn [& forms]
  (apply pr forms)
  (. sputter:*stdout* :write *newline*))

(defn print [& forms]
  (let [s (pr-map-with-nil str forms)]
    (if (seq? s) (. sputter:*stdout* :write (first s)))
    (for-each [e (rest s)]
      (. sputter:*stdout* :write *space* e))))

(defn println [& forms]
  (apply print forms)
  (. sputter:*stdout* :write *newline*))

(defmacro with-open [bindings & body]
  (cond
    (!vector? bindings)
      (println "not a vector") ; really explode

    (= (len bindings) 0)
      `(sputter:do ~@body)

    (>= (len bindings) 2)
      `(let [~(bindings 0) ~(bindings 1)]
        ~`(with-open [~@(rest (rest bindings))] ~@body)
        (let [close# (get ~(bindings 0) :close nil)]
          (when (apply? close#) (close#))))))
