package terraformcore

// TfArgs is an interface to validate terraform arguments
// It's shared between different terraform commands,
// And it's the common behavior for all terraform command's arguments
type TfArgs interface {
	AreValid() error
}
