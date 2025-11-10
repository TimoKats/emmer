// This interface shows which functions emmer needs to include a filesystem.
// E.g., if you can fulfill these functions in S3, then this interface allows
// you to include it in the main program. Note, you can set specific shared
// properties (like AWS_REGION, etc.). Also, you need to create a setup function
// that applies potential settings etc.

package server

type FileSystem interface {
	Ls() ([]string, error)
	Put(filename string, value any) error
	Get(filename string) (map[string]any, error)
	Del(filename string) error
}
