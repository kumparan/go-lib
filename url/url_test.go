package url

import (
	"fmt"
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
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url have a scheme", func(t *testing.T) {
		url := "https://magneto.com/feed"
		resURL, err := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is a directory", func(t *testing.T) {
		url := "/category/news"
		resURL, err := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "https://kumparan.com/category/news", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is not a directory", func(t *testing.T) {
		baseURL := "https://kumparan.com/feed/category/news"
		url := "kumparan.com/feed"
		resURL, err := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "https://kumparan.com/feed/category/kumparan.com/feed", resURL)
		assert.NoError(t, err)
	})

	t.Run("if url is a directory and base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "/category/news"
		resURL, err := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
		assert.NoError(t, err)
	})

	t.Run("if base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "category/news"
		resURL, err := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "kumparan.com/category/news", resURL)
		assert.NoError(t, err)
	})
}
