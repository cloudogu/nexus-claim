package infrastructure

import (
	"os"

	"io/ioutil"

	"github.com/cloudogu/nexus-claim/domain"
	"github.com/hashicorp/hcl"
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
	if _, err := os.Stat(dao.path); os.IsNotExist(err) {
		return model, errors.Wrapf(err, "could not find model at %s", dao.path)
	}

	bytes, err := ioutil.ReadFile(dao.path)
	if err != nil {
		return model, errors.Wrapf(err, "failed to read file %s", dao.path)
	}

	err = hcl.Unmarshal(bytes, &model)
	if err != nil {
		return model, errors.Wrapf(err, "failed to unmarshal model %s", dao.path)
	}

	return model, nil
}
