package database

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/dgraph-io/badger/v3/options"
)

type Store interface {
	Add(name string, SHA []byte, data []byte) error
	Get(name string) ([]byte, error)
	Remove(name string) error
	Update(name string, SHA []byte, data []byte) error
	List(details bool) ([]interface{}, error)
	FileExists(name string) bool
	SHAExists(SHA []byte) bool
	WordCount() (int64, error)
	WordFrequency() (map[string]int64, error)
	Close() error
}

type database struct {
	file *badger.DB
	sha  *badger.DB
}

type Config struct {
	Diskless  bool
	Path      string
	CacheSize int
}

type File struct {
	Name      string
	SHA       string
	Size      int64
	WordCount int64
}

// Open or Creat a new database
func New(config *Config) (*database, error) {
	file, err := badger.Open(badger.DefaultOptions(config.Path + "/file").
		WithLoggingLevel(badger.ERROR).
		WithCompression(options.None))
	if err != nil {
		return nil, err
	}
	sha, err := badger.Open(badger.DefaultOptions(config.Path + "/sha").
		WithLoggingLevel(badger.ERROR).
		WithCompression(options.None))
	if err != nil {
		return nil, err
	}

	go runValueLogGC(file)
	go runValueLogGC(sha)

	return &database{file, sha}, nil
}

// Close database files
func (db *database) Close() error {
	if err := db.file.Close(); err != nil {
		return err
	}
	if err := db.sha.Close(); err != nil {
		return err
	}
	return nil
}

// Check if File exists in the Store
func (db *database) FileExists(name string) bool {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	_, err := txn.Get([]byte(name))
	return err == nil
}

// Check if SHA exists in the Store
func (db *database) SHAExists(SHA []byte) bool {
	txn := db.sha.NewTransaction(false)
	defer txn.Discard()
	_, err := txn.Get(SHA)
	return err == nil
}

// Returns an array of list object containing file details
func (db *database) List(details bool) ([]interface{}, error) {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	opts.PrefetchValues = details
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	files := []interface{}{}
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		if details {
			file := File{Name: string(key)}
			sha, _ := item.ValueCopy(nil)
			file.SHA = fmt.Sprintf("%x", sha)
			size, _ := db.getSize(sha)
			file.Size = size
			count, _ := db.getWordCount(key)
			file.WordCount = count
			files = append(files, file)
		} else {
			files = append(files, string(key))
		}
	}
	return files, nil
}

// Remove a file from the Store
func (db *database) Remove(name string) error {
	sha, err := db.getSHA([]byte(name))
	if err != nil {
		return err
	}
	if err := db.removeFile([]byte(name)); err != nil {
		return err
	}
	if err := db.removeSHA(sha); err != nil {
		return err
	}
	return nil
}

// Remove a file record
func (db *database) removeFile(key []byte) error {
	return db.file.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// Remove a SHA record
func (db *database) removeSHA(key []byte) error {
	return db.sha.Update(func(txn *badger.Txn) error {
		if _, err := db.getName(key); err == badger.ErrKeyNotFound {
			return txn.Delete([]byte(key))
		}
		return nil
	})
}

// Add a file to the store
func (db *database) Add(name string, SHA []byte, data []byte) error {
	if len(SHA) == 0 {
		newSHA := sha256.Sum256(data)
		SHA = newSHA[:]
	}
	if err := db.addFile([]byte(name), SHA); err != nil {
		return err
	}
	if err := db.addSHA(SHA, data); err != nil {
		return err
	}
	return nil
}

// Add a new file record
func (db *database) addFile(key []byte, value []byte) error {
	return db.file.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			return txn.Set(key, value)
		}
		return badger.ErrConflict
	})
}

// Add a new SHA record
func (db *database) addSHA(key []byte, value []byte) error {
	return db.sha.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			return txn.Set(key, value)
		}
		return nil
	})
}

// Update a file
func (db *database) Update(name string, SHA []byte, data []byte) error {
	if len(SHA) == 0 {
		newSHA := sha256.Sum256(data)
		SHA = newSHA[:]
	}
	oldSHA, err := db.getSHA([]byte(name))
	if err := db.updateFile([]byte(name), SHA); err != nil {
		return err
	}
	if err == nil && !bytes.Equal(oldSHA, SHA) {
		if err := db.removeSHA(oldSHA); err != nil {
			return err
		}
	}
	if err := db.updateSHA(SHA, data); err != nil {
		return err
	}
	return nil
}

// Update the File record
func (db *database) updateFile(key []byte, value []byte) error {
	return db.file.Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Update the SHA record
func (db *database) updateSHA(key []byte, value []byte) error {
	return db.sha.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get(key); err == badger.ErrKeyNotFound {
			return txn.Set(key, value)
		}
		return nil
	})
}

// Returns the data of a file identified by file name
func (db *database) Get(name string) ([]byte, error) {
	sha, err := db.getSHA([]byte(name))
	if err != nil {
		return nil, err
	}
	data, err := db.getData(sha)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Returns the data in a file identified by SHA
func (db *database) getData(key []byte) ([]byte, error) {
	txn := db.sha.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

// Returns the SHA of a file indentified by the File Name
func (db *database) getSHA(key []byte) ([]byte, error) {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		return nil, err
	}
	return item.ValueCopy(nil)
}

// Returns the size of a given file identified by SHA
func (db *database) getSize(key []byte) (int64, error) {
	txn := db.sha.NewTransaction(false)
	defer txn.Discard()
	item, err := txn.Get(key)
	if err != nil {
		return 0, err
	}
	return item.ValueSize(), nil
}

// Returns name a file identified by the SHA
func (db *database) getName(sha []byte) ([]byte, error) {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		err := item.Value(func(value []byte) error {
			if string(value) != string(sha) {
				return badger.ErrKeyNotFound
			}
			return nil
		})
		if err == nil {
			return key, nil
		}
	}
	return nil, badger.ErrKeyNotFound
}

// Returns total word count in a single file in Store
func (db *database) getWordCount(key []byte) (int64, error) {
	var count int64
	sha, err := db.getSHA(key)
	if err != nil {
		return 0, err
	}
	data, err := db.getData(sha)
	if err != nil {
		return 0, err
	}
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		count++
	}
	return count, nil
}

// Returns frequency of words in a single file in Store
func (db *database) getWordFrequency(key []byte, frequency map[string]int64, mutex *sync.RWMutex) error {
	sha, err := db.getSHA(key)
	if err != nil {
		return err
	}
	data, err := db.getData(sha)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		word := scanner.Text()
		word = strings.ToLower(word)
		mutex.Lock()
		frequency[word]++
		mutex.Unlock()
	}
	return nil
}

// Returns total word count in all the files in Store
func (db *database) WordCount() (int64, error) {
	var totalCount int64
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.PrefetchSize = 10
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	wg := &sync.WaitGroup{}
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		wg.Add(1)
		go func() {
			defer wg.Done()
			count, _ := db.getWordCount(key)
			totalCount += count
		}()
	}
	wg.Wait()
	return totalCount, nil
}

// Returns frequency of words in all the files in Store
func (db *database) WordFrequency() (map[string]int64, error) {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.PrefetchSize = 10
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	frequency := map[string]int64{}
	wg := &sync.WaitGroup{}
	mutex := &sync.RWMutex{}
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		wg.Add(1)
		go func() {
			defer wg.Done()
			db.getWordFrequency(key, frequency, mutex)
		}()
	}
	wg.Wait()
	return frequency, nil
}

// Run log garbage collector after every 5 seconds
func runValueLogGC(db *badger.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
	again:
		err := db.RunValueLogGC(0.6)
		if err == nil {
			goto again
		}
	}
}

// Development and Debugging purpose
func (db *database) ListHelpFiles() ([]interface{}, error) {
	txn := db.file.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	files := []interface{}{}
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		files = append(files, string(key))
	}
	return files, nil
}

func (db *database) ListHelpSHA() ([]interface{}, error) {
	txn := db.sha.NewTransaction(false)
	defer txn.Discard()
	opts := badger.DefaultIteratorOptions
	opts.PrefetchSize = 10
	iterator := txn.NewIterator(opts)
	defer iterator.Close()
	files := []interface{}{}
	for iterator.Rewind(); iterator.Valid(); iterator.Next() {
		item := iterator.Item()
		key := item.KeyCopy(nil)
		files = append(files, string(key))
	}
	return files, nil
}
