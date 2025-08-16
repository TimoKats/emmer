package server

func (entry *EntryPayload) add() error {
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.UpdateJSON(path, entry.Key, entry.Value)
}

func (entry *EntryPayload) del() error {
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.DelJSON(path, entry.Key)
}
