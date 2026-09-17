package services

import (
	"path/filepath"
	"testing"
)

func TestImageURLFor(t *testing.T) {
	cases := []struct {
		name       string
		baseURL    string
		uploadPath string
		fsPath     string
		want       string
	}{
		{
			name:       "local relative path",
			uploadPath: "./uploads",
			fsPath:     filepath.Join("uploads", "products", "a.jpg"),
			want:       "/uploads/products/a.jpg",
		},
		{
			name:       "local relative path with baseURL",
			baseURL:    "http://localhost:8080",
			uploadPath: "./uploads",
			fsPath:     filepath.Join("uploads", "products", "a.jpg"),
			want:       "http://localhost:8080/uploads/products/a.jpg",
		},
		{
			name:       "docker absolute path",
			baseURL:    "http://localhost:8080",
			uploadPath: "/app/uploads",
			fsPath:     "/app/uploads/products/a.jpg",
			want:       "http://localhost:8080/uploads/products/a.jpg",
		},
		{
			name:       "baseURL trailing slash trimmed",
			baseURL:    "http://localhost:8080/",
			uploadPath: "/app/uploads",
			fsPath:     "/app/uploads/products/a.jpg",
			want:       "http://localhost:8080/uploads/products/a.jpg",
		},
		{
			name:       "path outside root falls back to filename",
			uploadPath: "/app/uploads",
			fsPath:     "/tmp/evil.jpg",
			want:       "/uploads/products/evil.jpg",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := NewProductImageService(nil, nil, nil, tc.baseURL, tc.uploadPath)

			if got := svc.imageURLFor(tc.fsPath); got != tc.want {
				t.Fatalf("imageURLFor(%q) = %q, want %q", tc.fsPath, got, tc.want)
			}
		})
	}
}
