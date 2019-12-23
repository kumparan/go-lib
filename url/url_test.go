package url

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UrlJoin(t *testing.T) {
	baseURL := "https://kumparan.com/trending"
	emptyString := ""

	t.Run("url empty", func(t *testing.T) {
		url := ""
		resURL, err := Join(baseURL, url)
		assert.Equal(t, baseURL, resURL)
		assert.NoError(t, err)
	})

	t.Run("base empty", func(t *testing.T) {
		url := "trending"
		resURL, err := Join(emptyString, url)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have a scheme", func(t *testing.T) {
		url := "https://magneto.com/feed"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is a directory", func(t *testing.T) {
		url := "/category/news"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, "https://kumparan.com/category/news", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is not a directory", func(t *testing.T) {
		baseURL := "https://kumparan.com/feed/category/news"
		url := "kumparan.com/feed"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, "https://kumparan.com/feed/category/kumparan.com/feed", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is a directory and base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "/category/news"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "category/news"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, "kumparan.com/category/news", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have \"..\"", func(t *testing.T) {
		baseURL := "https://kumparan.com/trending/feed"
		url := "../indo.nesia/news"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, "https://kumparan.com/indo.nesia/news", resURL)
		assert.NoError(t, err)

		t.Run("if url have \"..\" on the last element", func(t *testing.T) {
			baseURL := "https://kumparan.com/trending/feed"
			url := "sepakbola/.."
			resURL, err := Join(baseURL, url)

			assert.Equal(t, "https://kumparan.com/trending/", resURL)
			assert.NoError(t, err)

		})
	})

	t.Run("if url have \".\"", func(t *testing.T) {
		baseURL := "https://kumparan.com/trending/feed"
		url := "./indo.nesia/news"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, "https://kumparan.com/trending/indo.nesia/news", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have a different scheme than baseURL", func(t *testing.T) {
		baseURL := "https://kumparan.com/trending/feed"
		url := "http://kumparan.com/trending"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have no path", func(t *testing.T) {
		baseURL := "https://kumparan.com/trending/feed"
		url := "https://"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, baseURL, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have port", func(t *testing.T) {
		baseURL := "https://192.168.1:3000/trending/feed"
		url := "/story"
		expectedURL := "https://192.168.1:3000/story"
		resURL, err := Join(baseURL, url)
		assert.Equal(t, expectedURL, resURL)
		assert.NoError(t, err)
	})
}
