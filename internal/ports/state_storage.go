package port

import "context"

// StateStorage abstracts the remote storage of the session state file.
// Consumers depend on this interface and not on the concrete
// implementation (e.g. GitHub Gist).

type StateStorage interface {
	DownloadState(ctx context.Context, destinationPath string) error
	UploadState(ctx context.Context, sourcePath string) error
}