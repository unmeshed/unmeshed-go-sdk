package common

type StepDependency struct {
	Ref  string `json:"ref"`
	Deps string `json:"deps"`
}
