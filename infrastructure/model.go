package infrastructure

import (
	"fmt"
	"os"

	"io/ioutil"

	"strings"

	"reflect"

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

	file, err := dao.parseFile()
	if err != nil {
		return model, err
	}

	repositories, err := dao.findAndParse(file, "repository", domain.TypeRepository)
	if err != nil {
		return model, err
	}

	groups, err := dao.findAndParse(file, "repository_group", domain.TypeGroup)
	if err != nil {
		return model, err
	}

	model.Repositories = append(repositories, groups...)

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

func (dao *fileModelDAO) findAndParse(file *ast.File, nodeName string, repositoryType domain.RepositoryType) ([]domain.ModelRepository, error) {
	nodes, err := dao.findNodes(file, nodeName)
	if err != nil {
		return nil, err
	}

	return dao.parseNodes(nodes, repositoryType)
}

func (dao *fileModelDAO) findNodes(file *ast.File, nodeName string) ([]*ast.ObjectItem, error) {
	if objectList, ok := file.Node.(*ast.ObjectList); ok {
		return objectList.Filter(nodeName).Items, nil
	}
	return nil, errors.Errorf("file does not start with an objectlist")
}

func (dao *fileModelDAO) parseNodes(nodes []*ast.ObjectItem, repositoryType domain.RepositoryType) ([]domain.ModelRepository, error) {
	repositories := []domain.ModelRepository{}
	for _, repositoryNode := range nodes {
		repository, err := dao.parseRepositoryNode(repositoryNode)
		if err != nil {
			return nil, err
		}

		repository.Type = repositoryType
		repositories = append(repositories, repository)
	}
	return repositories, nil
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

	dao.normalizeProperties(repository.Properties)

	return repository, nil
}
func (dao *fileModelDAO) parseID(repositoryNode *ast.ObjectItem) (domain.RepositoryID, error) {
	keyCount := len(repositoryNode.Keys)
	if keyCount < 0 || keyCount > 1 {
		return "", errors.New("found repository with more or less than one id")
	}

	key := repositoryNode.Keys[0]

	id := strings.Trim(key.Token.Text, "\"")
	if id == "" {
		return "", errors.New("repository with empty id found")
	}

	return domain.RepositoryID(id), nil
}

func (dao *fileModelDAO) parseProperties(id domain.RepositoryID, repositoryNode *ast.ObjectItem) (domain.Properties, error) {
	properties := make(domain.Properties)

	err := hcl.DecodeObject(&properties, repositoryNode.Val)
	if err != nil {
		return properties, errors.Wrapf(err, "failed to parse properties of repository %s", id)
	}

	return properties, nil
}

func (dao *fileModelDAO) parseState(repositoryID domain.RepositoryID, properties domain.Properties) (domain.State, error) {
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

func (dao *fileModelDAO) normalizeProperties(properties domain.Properties) {
	// remove _state from repository property.
	// The state property can cause problems on create or modify, because nexus tries to interpret it.
	delete(properties, "_state")
	dao.unwrapNestedProperties(properties)
}

func (dao *fileModelDAO) unwrapNestedProperties(properties map[string]interface{}) {
	for key, value := range properties {
		if dao.isNestedProperty(value) {
			fmt.Println("unwrap", key)

			unwrapped := dao.unwrapNestedProperty(value)

			dao.unwrapNestedProperties(unwrapped.(map[string]interface{}))

			properties[key] = unwrapped
		}
	}
}

func (dao *fileModelDAO) isNestedProperty(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice && v.Len() == 1 {
		elem := v.Type().Elem()
		return elem.Kind() == reflect.Map
	}
	return false
}

func (dao *fileModelDAO) unwrapNestedProperty(value interface{}) interface{} {
	return value.([]map[string]interface{})[0]
}
