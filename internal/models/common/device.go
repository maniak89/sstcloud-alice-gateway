package common

type Device struct {
	ID               string
	Name             string
	Model            string
	Enabled          bool
	Connected        bool
	Tempometer       *Tempometer
	AdditionalFields map[string]string
}

type Tempometer struct {
	Degressess int
}
