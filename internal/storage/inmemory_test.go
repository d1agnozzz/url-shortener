package storage

import (
	"context"
	"slices"
	"sync"
	"testing"

	"github.com/d1agnozzz/url-shortener/internal/types"
)

func Test_inMemStorage_CreateDuplicateURL(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage()

	toInsert := []types.URLMapping{
		{
			Url:   "SAME",
			Alias: "different_1",
		},
		{
			Url:   "SAME",
			Alias: "different_2",
		},
	}

	err1 := storage.InsertURLMapping(ctx, toInsert[0])

	if err1 != nil {
		t.Fatalf("first insertion error: %v", err1)
	}

	err2 := storage.InsertURLMapping(ctx, toInsert[1])

	if err2 == nil {
		t.Fatalf("duplicate insertion didn't return error")
	}
}

func Test_inMemStorage_GetNotFound(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage()

	_, err := storage.GetByAlias(ctx, "non existent alias")

	if err == nil {
		t.Fatalf("non nil error on failed get")
	}
}

func Test_inMemStorage_ConcurrentInsert(t *testing.T) {
	ctx := context.Background()
	storage := NewInMemoryStorage()

	toInsert := []types.URLMapping{
		{
			Url:   "example1.com",
			Alias: "1",
		},
		{
			Url:   "example2.com",
			Alias: "2",
		},
		{
			Url:   "example3.com",
			Alias: "3",
		},
		{
			Url:   "example4.com",
			Alias: "4",
		},
		{
			Url:   "example5.com",
			Alias: "5",
		},
	}

	errCh := make(chan error, len(toInsert))

	var wg sync.WaitGroup
	wg.Add(len(toInsert))

	for i := range len(toInsert) {
		go func(n int) {
			defer wg.Done()
			err := storage.InsertURLMapping(ctx, toInsert[n%len(toInsert)])

			if err != nil {
				_, dupl := err.(InsertDuplicateError)

				// ignore duplicate for test
				if !dupl {
					errCh <- err
				}

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
	var urls []string

	for _, v := range toInsert {
		urls = append(urls, v.Url)
	}

	// cast to concrete for testing
	storage_, _ := storage.(*inMemStorage)

	unique_aliases := make([]string, 5)
	for al, val := range storage_.urlMappings {
		if al != val.Alias {
			t.Fatalf("data is corrupted!")
		}

		if slices.Contains(unique_aliases, al) {
			t.Fatalf("data is corrupted!")
		}
		unique_aliases = append(unique_aliases, al)

		if !slices.Contains(urls, val.Url) {
			t.Fatalf("data is corrupted!")
		}

	}
}
