package config

type Evaluatable interface {
	Evaluate() error
}

type Validatable interface {
	Validate() error
}

type SettingsInterface interface {
	Evaluatable
	Validatable
}
