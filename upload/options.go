package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/meatballhat/artifacts/env"
	"github.com/mitchellh/goamz/s3"
)

var (
	// DefaultCacheControl is the default value for each artifact's Cache-Control header
	DefaultCacheControl = "private"
	// DefaultConcurrency is the default number of concurrent goroutines used during upload
	DefaultConcurrency = uint64(5)
	// DefaultMaxSize is the default maximum allowed bytes for all artifacts
	DefaultMaxSize = uint64(1024 * 1024 * 1000)
	// DefaultPaths is the default slice of local paths to upload (empty)
	DefaultPaths = []string{}
	// DefaultPerm is the default ACL applied to each artifact
	DefaultPerm = "private"
	// DefaultRetries is the default number of times a given artifact upload will be retried
	DefaultRetries = uint64(2)
	// DefaultTargetPaths is the default upload prefix for each artifact
	DefaultTargetPaths = []string{}
	// DefaultUploadProvider is the provider used to upload (nuts)
	DefaultUploadProvider = "s3"
	// DefaultWorkingDir is the default working directory ... wow.
	DefaultWorkingDir, _ = os.Getwd()
)

// Options is used in the call to Upload
type Options struct {
	AccessKey    string
	BucketName   string
	CacheControl string
	Concurrency  uint64
	MaxSize      uint64
	Paths        []string
	Perm         s3.ACL
	Provider     string
	Retries      uint64
	SecretKey    string
	TargetPaths  []string
	WorkingDir   string
}

func init() {
	DefaultTargetPaths = append(DefaultTargetPaths, filepath.Join("artifacts",
		env.Get("TRAVIS_BUILD_NUMBER", ""),
		env.Get("TRAVIS_JOB_NUMBER", "")))
}

// NewOptions makes some *Options with defaults!
func NewOptions() *Options {
	cwd, _ := os.Getwd()
	cwd = env.Get("TRAVIS_BUILD_DIR", cwd)

	targetPaths := env.ExpandSlice(env.Slice("ARTIFACTS_TARGET_PATHS", ":", []string{}))
	if len(targetPaths) == 0 {
		targetPaths = DefaultTargetPaths
	}

	return &Options{
		AccessKey: env.Cascade([]string{
			"ARTIFACTS_KEY",
			"ARTIFACTS_AWS_ACCESS_KEY",
			"AWS_ACCESS_KEY_ID",
			"AWS_ACCESS_KEY",
		}, ""),
		SecretKey: env.Cascade([]string{
			"ARTIFACTS_SECRET",
			"ARTIFACTS_AWS_SECRET_KEY",
			"AWS_SECRET_ACCESS_KEY",
			"AWS_SECRET_KEY",
		}, ""),
		BucketName: strings.TrimSpace(env.Cascade([]string{
			"ARTIFACTS_BUCKET",
			"ARTIFACTS_S3_BUCKET",
		}, "")),

		CacheControl: strings.TrimSpace(env.Get("ARTIFACTS_CACHE_CONTROL", DefaultCacheControl)),
		Concurrency:  env.Uint("ARTIFACTS_CONCURRENCY", DefaultConcurrency),
		MaxSize:      env.UintSize("ARTIFACTS_MAX_SIZE", DefaultMaxSize),
		Paths:        env.ExpandSlice(env.Slice("ARTIFACTS_PATHS", ":", DefaultPaths)),
		Perm:         s3.ACL(env.Get("ARTIFACTS_PERMISSIONS", DefaultPerm)),
		Retries:      env.Uint("ARTIFACTS_RETRIES", DefaultRetries),
		TargetPaths:  targetPaths,
		Provider:     env.Get("ARTIFACTS_UPLOAD_PROVIDER", DefaultUploadProvider),
		WorkingDir:   cwd,
	}
}

// Validate checks for validity!
func (opts *Options) Validate() error {
	if opts.BucketName == "" {
		return fmt.Errorf("no bucket name given")
	}

	if opts.AccessKey == "" {
		return fmt.Errorf("no access key given")
	}

	if opts.SecretKey == "" {
		return fmt.Errorf("no secret key given")
	}

	return nil
}
