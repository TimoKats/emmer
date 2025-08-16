// This interface shows which functions emmer needs to include a filesystem.
// E.g., if you can fulfill these functions in S3, then this interface allows
// you to include it in the main program. Note, you can set specific shared
// properties (like AWS_REGION, etc.)

package server

type IO interface {
	// file level
	Fetch(filename string) (string, error)
	DeleteFile(filename string) error
	Info() string

	// data level
	CreateJSON(path string) error
	ReadJSON(path string) (map[string]any, error)
	UpdateJSON(path string, key string, value any) error
	DelJSON(path string, key string) error
}
