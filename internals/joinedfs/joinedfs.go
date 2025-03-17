package joinedfs

import "io/fs"

func Join(fsys ...fs.FS) fs.FS {
	return &joinedFS{fsys}
}

type joinedFS struct {
	fsys []fs.FS
}

var _ fs.FS = (*joinedFS)(nil)

func (j *joinedFS) Open(name string) (fs.File, error) {
	var err error
	var f fs.File
	for _, fsys := range j.fsys {
		f, err = fsys.Open(name)
		if err == nil {
			return f, nil
		}
	}
	return f, err
}
