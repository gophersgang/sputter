package api

// ArgumentProcessor is the standard signature for a function that is
// capable of processing an Iterable (like Lists)
type ArgumentProcessor func(Context, Sequence) Value

// Function is a Value that can be invoked
type Function struct {
	Name Name
	Exec ArgumentProcessor
}

func (f *Function) String() string {
	return string(f.Name)
}
