package server

// Fetches path on table name, then updates JSON.
// NOTE: Maybe check the entry/value values.
func (entry *EntryPayload) add() error {
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.UpdateJSON(path, entry.Key, entry.Value)
}

// Fetches path for table name, then removes key from JSON.
// NOTE: do some extra checks
func (entry *EntryPayload) del() error {
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.DeleteJson(path, entry.Key)
}
