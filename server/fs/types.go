// This interface shows which functions emmer needs to interface with a different
// filetype/filesytem. An implementation of this interface can be accessed using
// EM_CONNECTOR.

package server

type FileSystem interface {
	Ls() ([]string, error)
	Put(filename string, value any) error
	Get(filename string) (any, error)
	Del(filename string) error
}
