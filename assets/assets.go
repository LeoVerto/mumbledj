package assets

import (
    "bytes"
    "embed"
    "io/fs"
    "net/http"
    "path"
    "time"
)

// fsAssets embeds non-Go files from this directory recursively.
//go:embed **
var fsAssets embed.FS

// box is a minimal wrapper exposing the subset of packr.Box API used in the codebase.
type box struct{}

// Assets provides access to embedded asset files.
var Assets = &box{}

// List returns a list of file paths (relative to the assets package directory)
// for all embedded files. Only regular files are listed (directories are skipped).
func (b *box) List() []string {
    var files []string
    // Walk from "." because fsAssets is rooted at the package directory.
    _ = fs.WalkDir(fsAssets, ".", func(p string, d fs.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        // Skip the current directory and any directories; we only collect files.
        if d.IsDir() {
            return nil
        }
        // Exclude Go source files from the listing.
        if path.Ext(p) == ".go" {
            return nil
        }
        files = append(files, p)
        return nil
    })
    return files
}

// HasDir reports whether a directory with the given relative name exists in the embedded assets.
func (b *box) HasDir(name string) bool {
    info, err := fs.Stat(fsAssets, name)
    if err != nil {
        return false
    }
    return info.IsDir()
}

// Open opens a file from the embedded assets and returns it as fs.File.
func (b *box) Open(name string) (http.File, error) {
    // Read all to allow seeking over an in-memory buffer.
    data, err := fs.ReadFile(fsAssets, name)
    if err != nil {
        return nil, err
    }
    info, err := fs.Stat(fsAssets, name)
    if err != nil {
        return nil, err
    }
    return &memHTTPFile{
        Reader:  bytes.NewReader(data),
        name:    name,
        size:    int64(len(data)),
        modTime: info.ModTime(),
        isDir:   false,
    }, nil
}

// Find reads and returns the content of the named file from the embedded assets.
func (b *box) Find(name string) ([]byte, error) {
    return fs.ReadFile(fsAssets, name)
}

// memHTTPFile is an in-memory implementation of http.File backed by bytes.Reader.
type memHTTPFile struct {
    *bytes.Reader
    name    string
    size    int64
    modTime time.Time
    isDir   bool
}

func (f *memHTTPFile) Close() error { return nil }

// Readdir returns directory entries. For files, it returns an empty slice.
func (f *memHTTPFile) Readdir(count int) ([]fs.FileInfo, error) {
    if !f.isDir {
        return []fs.FileInfo{}, nil
    }
    return []fs.FileInfo{}, nil
}

// Stat returns file info for this in-memory file.
func (f *memHTTPFile) Stat() (fs.FileInfo, error) { return memFileInfo{f}, nil }

type memFileInfo struct{ f *memHTTPFile }

func (fi memFileInfo) Name() string       { return path.Base(fi.f.name) }
func (fi memFileInfo) Size() int64        { return fi.f.size }
func (fi memFileInfo) Mode() fs.FileMode  { return 0444 }
func (fi memFileInfo) ModTime() time.Time { return fi.f.modTime }
func (fi memFileInfo) IsDir() bool        { return fi.f.isDir }
func (fi memFileInfo) Sys() any           { return nil }
