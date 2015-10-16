package codeasist

type Package struct {
	Name  string
	Dir   string
	Files []File
}

type File struct {
	Path      string
	Imports   []Unit
	Vars      []Var
	Types     []Type
	Functions []Function
}

type Type struct {
	Unit
	Functions []Function
	Type      string
}

type Function struct {
	Args string
	Type string
	Recv string
	Unit
}

type Var struct {
	Type string
	Unit
}

type Unit struct {
	Name string
	Doc  string
	Pos
}

type Pos struct {
	FileName string
	Line     int
	Column   int
}
