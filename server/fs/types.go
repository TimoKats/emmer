// This interface shows which functions emmer needs to include a filesystem.
// E.g., if you can fulfill these functions in S3, then this interface allows
// you to include it in the main program. Note, you can set specific shared
// properties (like AWS_REGION, etc.). Also, you need to create a setup function
// that applies potential settings etc.

package server

type FileSystem interface {
	Fetch(filename string) (string, error)
	DeleteFile(filename string) error
	List() (map[string]any, error)
	Info() string

	CreateJSON(filename string) error
	ReadJSON(filename string) (map[string]any, error)
	UpdateJSON(filename string, key []string, value any) error
	DeleteJson(filename string, key []string) error
}
