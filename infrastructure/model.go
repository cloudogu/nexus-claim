package infrastructure

import (
	"os"

	"io/ioutil"

	"strings"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/pkg/errors"
)

// NewFileModelDAO creates a new file based ModelDAO
func NewFileModelDAO(path string) domain.ModelDAO {
	return &fileModelDAO{path}
}

type fileModelDAO struct {
	path string
}

func (dao *fileModelDAO) Get() (domain.Model, error) {
	model := domain.Model{}
	repositories := []domain.ModelRepository{}

	file, err := dao.parseFile()
	if err != nil {
		return model, err
	}

	repositoryNodes, err := dao.findRepositoryNodes(file)
	if err != nil {
		return model, err
	}

	for _, repositoryNode := range repositoryNodes {
		repository, err := dao.parseRepositoryNode(repositoryNode)
		if err != nil {
			return model, err
		}

		repositories = append(repositories, repository)
	}

	model.Repositories = repositories

	return model, nil
}

func (dao *fileModelDAO) parseFile() (*ast.File, error) {
	if _, err := os.Stat(dao.path); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "could not find model at %s", dao.path)
	}

	bytes, err := ioutil.ReadFile(dao.path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", dao.path)
	}

	file, err := hcl.ParseBytes(bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse file %s", dao.path)
	}

	return file, nil
}

func (dao *fileModelDAO) findRepositoryNodes(file *ast.File) ([]*ast.ObjectItem, error) {
	if objectList, ok := file.Node.(*ast.ObjectList); ok {
		return objectList.Filter("repository").Items, nil
	}
	return nil, errors.Errorf("file does not start with an objectlist")
}

func (dao *fileModelDAO) parseRepositoryNode(repositoryNode *ast.ObjectItem) (domain.ModelRepository, error) {
	repository := domain.ModelRepository{}

	id, err := dao.parseID(repositoryNode)
	if err != nil {
		return repository, err
	}
	repository.ID = id

	properties, err := dao.parseProperties(id, repositoryNode)
	if err != nil {
		return repository, err
	}
	repository.Properties = properties

	state, err := dao.parseState(id, properties)
	if err != nil {
		return repository, err
	}
	repository.State = state

	return repository, nil
}

func (dao *fileModelDAO) parseID(repositoryNode *ast.ObjectItem) (string, error) {
	keyCount := len(repositoryNode.Keys)
	if keyCount < 0 || keyCount > 1 {
		return "", errors.New("found repository with more or less than one id")
	}

	key := repositoryNode.Keys[0]

	id := strings.Trim(key.Token.Text, "\"")
	if id == "" {
		return "", errors.New("repository with empty id found")
	}

	return id, nil
}

func (dao *fileModelDAO) parseProperties(id string, repositoryNode *ast.ObjectItem) (domain.Properties, error) {
	properties := make(domain.Properties)

	err := hcl.DecodeObject(&properties, repositoryNode.Val)
	if err != nil {
		return properties, errors.Wrapf(err, "failed to parse properties of repository %s", id)
	}

	return properties, nil
}

func (dao *fileModelDAO) parseState(repositoryID string, properties domain.Properties) (domain.State, error) {
	state := domain.StatePresent
	if stateValue, ok := properties["_state"]; ok {

		if stateString, ok := stateValue.(string); ok {
			switch domain.State(stateString) {
			case domain.StatePresent:
				state = domain.StatePresent
			case domain.StateAbsent:
				state = domain.StateAbsent
			default:
				return state, errors.Errorf("state %s of repository %s is not a valid state", stateString, repositoryID)
			}

		} else {
			return state, errors.Errorf("state of repository %s is not a string", repositoryID)
		}

	} else {
		return state, errors.Errorf("repository %s has no _state field", repositoryID)
	}

	return state, nil
}
