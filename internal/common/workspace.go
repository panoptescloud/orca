package common

type ProjectRepository struct {
	SSH string
}

type Project struct {
	Name       string
	Repository ProjectRepository
	Path       string
}

type Workspace struct {
	Name     string
	Projects []Project
}

func (self *Workspace) GetProject(name string) (*Project, error) {
	for _, project := range self.Projects {
		if project.Name == name {
			return &project, nil
		}
	}

	return nil, ErrUnknownProject{
		Name: name,
	}
}

type WorkspaceLocation struct {
	Name string
	Path string
}

type WorkspaceLocations []WorkspaceLocation
