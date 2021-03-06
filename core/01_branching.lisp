;;;; sputter core: branching

(defmacro not
  {:doc "logically inverts the truthiness of the provided value"}
  [val]
  `(if ~val false true))

(defmacro when
  {:doc-asset "when"}
  [test & body]
  `(if ~test (do ~@body) nil))

(defmacro when-not
  {:doc-asset "when"}
  [test & body]
  `(if ~test nil (do ~@body)))

(defmacro cond
  {:doc-asset "cond"}
  ([] nil)
  ([clause] clause)
  ([& clauses]
    `(if ~(clauses 0)
         ~(clauses 1)
         (sputter:cond ~@(rest (rest clauses))))))

(defmacro and
  {:doc-asset "and"}
  ([] true)
  ([clause] clause)
  ([& clauses]
    `(let [and# ~(clauses 0)]
      (if and# (sputter:and ~@(rest clauses)) and#))))

(defmacro or
  {:doc-asset "or"}
  ([] nil)
  ([clause] clause)
  ([& clauses]
    `(let [or# ~(clauses 0)]
      (if or# or# (sputter:or ~@(rest clauses))))))
