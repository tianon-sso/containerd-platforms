module github.com/containerd/platforms

go 1.26.0

require (
	github.com/containerd/log v0.2.0
	github.com/opencontainers/image-spec v1.1.0
	golang.org/x/sys v0.48.0
)

require (
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/sirupsen/logrus v1.10.2 // indirect
)

// TODO remove this when v0.49.0+ is released containing the v2/v3/v4 feature bits (merged upstream 2026-09-18, commit 01b91195)
replace golang.org/x/sys => golang.org/x/sys v0.48.1-0.20260918203432-01b91195d9ae
