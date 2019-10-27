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
		resURL := Join(baseURL, url)
		assert.Equal(t, baseURL, resURL)
	})

	t.Run("base empty", func(t *testing.T) {
		url := "trending"
		resURL := Join(emptyString, url)
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
	})

	t.Run("if url have a scheme", func(t *testing.T) {
		url := "https://magneto.com/feed"
		resURL := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
	})

	t.Run("if url is a directory", func(t *testing.T) {
		url := "/category/news"
		resURL := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "https://kumparan.com/category/news", resURL)
	})

	t.Run("if url is not a directory", func(t *testing.T) {
		baseURL := "https://kumparan.com/feed/category/news"
		url := "kumparan.com/feed"
		resURL := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "https://kumparan.com/feed/category/kumparan.com/feed", resURL)
	})

	t.Run("if url is a directory and base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "/category/news"
		resURL := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, url, resURL)
	})

	t.Run("if base have no scheme", func(t *testing.T) {
		baseURL := "kumparan.com/trending"
		url := "category/news"
		resURL := Join(baseURL, url)
		fmt.Println(resURL)
		assert.Equal(t, "kumparan.com/category/news", resURL)
	})
}
