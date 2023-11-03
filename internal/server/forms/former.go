package forms

type Former interface {
	ParseAndValidate(message interface{}) error
	ConvertToMap() map[string]interface{}
}
