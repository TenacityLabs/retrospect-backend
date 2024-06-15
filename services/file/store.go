package file

type FileStore struct {
}

func NewFileStore() *FileStore {
	return &FileStore{}
}
func (fileStore *FileStore) UploadFile(file []byte, filename string) (string, error) {
	return "", nil
}
