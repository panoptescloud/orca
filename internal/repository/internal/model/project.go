package model

type LoaderPropertyCondition struct {
	Name  string
	Value any
}

type LoaderCondition struct {
	OS       string
	Arch     string
	Property LoaderPropertyCondition
}

type ExtraComposeFile struct {
	Path string
	When []LoaderCondition
}

type ComposeFiles struct {
	Primary string
	Extras  []ExtraComposeFile
}

type Property struct {
	Name    string
	Type    string
	Default any
}

type Extension struct {
	Name    string
	Chdir   string
	Command string
	Service string
}

type ProjectConfig struct {
	ComposeFiles ComposeFiles `yaml:"composeFiles"`
	Properties   []Property
	Hosts        []string
	Extensions   []Extension
}
