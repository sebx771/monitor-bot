package port

import "context"

// StateStorage abstrae el almacenamiento remoto del archivo de estado de
// sesión. Los consumidores dependen de esta interfaz y no de la
// implementación concreta (p.ej. GitHub Gist).

type StateStorage interface {
	DownloadState(ctx context.Context, destinationPath string) error
	UploadState(ctx context.Context, sourcePath string) error
}