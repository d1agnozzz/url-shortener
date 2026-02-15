package storage

import (
	"slices"
	"sync"
	"testing"

	"github.com/d1agnozzz/url-shortener/internal/aliaser"
	"github.com/d1agnozzz/url-shortener/internal/types"
	"github.com/d1agnozzz/url-shortener/internal/urlsanitizer"
)

const NGOROUTINES = 16

func Test_inMemStorage_CreateAndGet(t *testing.T) {
	storage := NewInMemoryStorage(StorageConfig{
		maxCollisions: 5,
		aliaser:       aliaser.NewMd5Aliaser(),
	})

	sanitizer := urlsanitizer.NewUrlSanitizer()

	url, _ := sanitizer.Sanitize("youtube.com")

	mapping, err := storage.CreateURLMapping(*url)

	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if mapping.Url != url.String() {
		t.Fatalf("created mapping url is different from given: got %s, want: %s", mapping.Url, url.String())
	}

	retrieved, err := storage.GetByAlias(mapping.Alias)

	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if retrieved.Url != url.String() {
		t.Fatalf("retrieved url is different: got '%s', want '%s'", url.String(), retrieved.Url)
	}
}

func Test_inMemStorage_CollisionResolution(t *testing.T) {
	storage := NewInMemoryStorage(StorageConfig{
		maxCollisions: 1,
		aliaser:       aliaser.NewMd5Aliaser(),
	})

	sanitizer := urlsanitizer.NewUrlSanitizer()

	url, _ := sanitizer.Sanitize("first.com")

	mapping, err := storage.CreateURLMapping(*url)
	if err != nil {
		t.Fatal(err)
	}

	// manual url override
	storage.urlMappings[mapping.Alias] = types.URLMapping{
		Url:   "DIFFERENT URL",
		Alias: mapping.Alias,
	}

	newMapping, err := storage.CreateURLMapping(*url)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}

	if newMapping.Url != url.String() {
		t.Fatalf("wrong url mapping after collision")
	}

	if newMapping.Alias == mapping.Alias {
		t.Fatalf("different urls with the same alias")
	}

}

func Test_inMemStorage_CreateDuplicateURL(t *testing.T) {
	storage := NewInMemoryStorage(StorageConfig{
		maxCollisions: 5,
		aliaser:       aliaser.NewMd5Aliaser(),
	})
	sanitizer := urlsanitizer.NewUrlSanitizer()

	url, _ := sanitizer.Sanitize("reddit.com")

	firstMap, err := storage.CreateURLMapping(*url)
	if err != nil {
		t.Fatal(err)
	}

	secondMap, err := storage.CreateURLMapping(*url)
	if err != nil {
		t.Fatal(err)
	}

	if *firstMap != *secondMap {
		t.Fatalf("duplicate url returns different maps")
	}
}

func Test_inMemStorage_GetNotFound(t *testing.T) {
	storage := NewInMemoryStorage(StorageConfig{
		maxCollisions: 5,
		aliaser:       aliaser.NewMd5Aliaser(),
	})

	_, err := storage.GetByAlias("non existent alias")

	if err == nil {
		t.Fatalf("non nil error on failed get")
	}
}

func Test_inMemStorage_ConcurrentCreate(t *testing.T) {
	storage := NewInMemoryStorage(StorageConfig{
		maxCollisions: 3,
		aliaser:       aliaser.NewMd5Aliaser(),
	})
	sanitizer := urlsanitizer.NewUrlSanitizer()

	rawUrls := []string{
		"example1.com",
		"example2.com",
		"example3.com",
		"example4.com",
		"example5.com",
	}

	sanitizedUrls := make([]string, 0)

	for _, v := range rawUrls {
		sanitized, err := sanitizer.Sanitize(v)
		if err != nil {
			t.Fatal(err)
		}
		sanitizedUrls = append(sanitizedUrls, sanitized.String())
	}

	errCh := make(chan error, NGOROUTINES)

	var wg sync.WaitGroup
	wg.Add(NGOROUTINES)

	for i := range NGOROUTINES {
		go func(n int) {
			defer wg.Done()
			url, _ := sanitizer.Sanitize(sanitizedUrls[n%len(sanitizedUrls)])
			_, err := storage.CreateURLMapping(*url)

			if err != nil {
				errCh <- err
			}
		}(i)
	}

	wg.Wait()
	close(errCh)

	var errors []error

	for err := range errCh {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		t.Errorf("encountered %d errors during concurent creates:", len(errors))
		for _, err := range errors {
			t.Errorf(" - %v", err)
		}
	}

	// data validation
	unique_aliases := make([]string, 5)

	for al, val := range storage.urlMappings {
		if al != val.Alias {
			t.Fatalf("data is corrupted!")
		}

		if slices.Contains(unique_aliases, al) {
			t.Fatalf("data is corrupted!")
		}
		unique_aliases = append(unique_aliases, al)

		if !slices.Contains(sanitizedUrls, val.Url) {
			t.Fatalf("data is corrupted!")
		}

	}
}
