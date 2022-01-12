package common

type Device struct {
	ID               string
	Name             string
	Model            string
	Tempometer       *Tempometer
	AdditionalFields map[string]string
}

type Tempometer struct {
	Degressess map[string]int
}
