package domain

// State describes the state of an repository. State is inspired by ansible state.
// The state can be present or absent.
type State string

const (
	// StatePresent ensures that the repository exists
	StatePresent State = "present"
	// StateAbsent ensures that the repository does not exists
	StateAbsent State = "absent"
)

// ModelRepository represents the repository from a hcl file
type ModelRepository struct {
	Repository
	State
}

// Model is the structure which can be read from a hcl file
type Model struct {
	Repositories []ModelRepository
}

// ModelDAO contains methods to read and parse hcl files
type ModelDAO interface {
	Get() (Model, error)
}
