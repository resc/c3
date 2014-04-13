package c3

type generatorIterable struct {
	g Generator
}

func (i *generatorIterable) Iterator() Iterator {
	return &generateIterator{i.g(), defaultElementValue}
}
