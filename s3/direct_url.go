package s3

type DirectURLFactory struct {
	// Base should include transport
	base string
}

type DirectURL struct {
	base string

	// Path should start with a slash
	path string
}

func NewDirectURLFactory(base string) URLFactory {
	return DirectURLFactory{base: base}
}

func (f DirectURLFactory) New(path, fileName string) URL {
	return DirectURL{base: f.base, path: path}
}

func (u DirectURL) String() (string, error) { return u.base + u.path, nil }
