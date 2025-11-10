package interfaces

// IMapper is a generic interface for serializing and deserializing objects of type T.
type IMapper[T any] interface {
	// SerializeToFile serializes an object of type T to a file in the specified format.
	SerializeToFile(format string)
	// DeserializeFromFile deserializes an object of type T from a file in the specified format.
	DeserializeFromFile(format string) (*T, error)
	// Serialize converts an object of type T to a byte array in the specified format.
	Serialize(format string) ([]byte, error)
	// Deserialize converts a byte array in the specified format to an object of type T.
	Deserialize(data []byte, format string) (*T, error)
}
