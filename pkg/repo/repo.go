package repo

import (
	"context"
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

func CloneRepository(ctx context.Context, repoUrl string, repoUsername string, repoPassword string, defaultRepoUrl string) (*object.Tree, error) {
	cloneOptions := new(git.CloneOptions)
	cloneOptions.Depth = 1 // We are only interested in the last commit

	if repoUrl != "" {
		cloneOptions.URL = repoUrl

		if repoUsername != "" && repoPassword != "" {
			// If the user provided a repository URL, we use it with the provided credentials
			cloneOptions.Auth = &http.BasicAuth{
				Username: repoUsername,
				Password: repoPassword,
			}
		}
	} else {
		// The default repository that is used if the user doesn't specify another one
		cloneOptions.URL = defaultRepoUrl
	}

	// Filesystem abstraction based on memory
	fs := memfs.New()

	// Git objects memStorage based on memory
	memStorage := memory.NewStorage()

	// Clones the repository into the worktree (fs) and stores all the git content into the memStorage
	repo, err := git.CloneContext(ctx, memStorage, fs, cloneOptions)
	if err != nil {
		if err.Error() == "authentication required" {
			return nil, fmt.Errorf("failed to authenticate repository access - check if the repository exists and the access credentials are provided")
		}

		return nil, err
	}

	// Retrieve the HEAD reference
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	return tree, nil
}
