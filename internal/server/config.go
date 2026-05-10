package server

import "io/fs"

type Config struct {
	Addr string
	FS   fs.FS
}
