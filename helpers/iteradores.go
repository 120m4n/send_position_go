package helpers

type CircularIterator struct {
    data   [][][2]float64
    index  int
}

func NewCircularIterator(data [][][2]float64) *CircularIterator {
    return &CircularIterator{
        data:   data,
        index:  0,
    }
}

func (ci *CircularIterator) Next() [][2]float64 {
    if len(ci.data) == 0 {
        return nil
    }
    currentData := ci.data[ci.index]
    ci.index = (ci.index + 1) % len(ci.data)
    return currentData
}