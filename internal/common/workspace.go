package common

type ProjectRepository struct {
	SSH string
}

type Project struct {
	Name       string
	Repository ProjectRepository
}

type Workspace struct {
	Name     string
	Projects []Project
}

type WorkspaceLocation struct {
	Name string
	Path string
}

type WorkspaceLocations []WorkspaceLocation
