package sidecar


import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "io/ioutil"
	"github.com/fsnotify/fsnotify"
    v1 "k8s.io/api/core/v1"
)



// transform keys of Secret as file with sub-path
// ex: "a.b.c.d" -> "a/b/c.d"
// From the context
//		- secret *v1.Secret
//		- basePath String (ex : "/a/")
func SecretToPath(c *Callback) error {
	// Check context
	secret, ok := c.Context["secret"].(*v1.Secret)
	if !ok {
	    return fmt.Errorf("missing or incorrect type for the context: secret")
	}
	basePath, ok := c.Context["basePath"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: basePath")
	}

    for key, value := range secret.Data {
        // transformer la clé
        parts := strings.Split(key, ".")
        if len(parts) > 1 {
            // tout sauf le dernier devient des dossiers
            dir := filepath.Join(basePath, filepath.Join(parts[:len(parts)-1]...))
            if err := os.MkdirAll(dir, 0o755); err != nil {
                return fmt.Errorf("mkdir: %w", err)
            }
            filePath := filepath.Join(dir, parts[len(parts)-1])
            if err := os.WriteFile(filePath, value, 0o644); err != nil {
                return fmt.Errorf("write: %w", err)
            }
        } else {
            // clé sans point
            filePath := filepath.Join(basePath, key)
            if err := os.WriteFile(filePath, value, 0o644); err != nil {
                return fmt.Errorf("write: %w", err)
            }
        }
    }
    return c.CallNext()
}


// Read a file as a Secret
// From the context
//		- mountedVolume String (ex : "/a")
//		- filePathInVolume String (ex : "b/c.d")
//		- (optional) secretPrefix String (ex: "lnd-1")
// Add in the context :
//		- secretKey String (ex: "b.c.d")
//		- secretData []byte (content of the file /a/b/c.d)
func ReadFileAsSecret(c *Callback) error {
	// Check context
	mountedVolume, ok := c.Context["mountedVolume"].(string)
	if  !ok {
	    return fmt.Errorf("missing or incorrect type for the context: mountedVolume")
	}
	path, ok := c.Context["filePathInVolume"].(string)
	if !ok {
		return fmt.Errorf("missing or incorrect type for the context: filePathInVolume")
	}

	// full path of file
	fullPath := filepath.Join(mountedVolume, path)

	// get the file
	secretData, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", fullPath, err)
	}
	c.Context["secretData"] = secretData

	// transform "a/b/c.d" -> "a.b.c.d"
	key := strings.ReplaceAll(path, "/", ".")
	if secretPrefix, ok := c.Context["secretPrefix"].(string) ; ok {
		key = secretPrefix + "." + key
	}

	c.Context["secretKey"] = key

	return c.CallNext()
}

// Routine if a file (recursively in the folder) that match the pattern is updated
// From the context
//		- basePath String (ex : "/a/")
//		- pattern String (ex : "*.d")
// Add in the context :
//		- filePath String (ex: ""/a/b/c.d")
func WatchFilePattern(c *Callback) error {
	basePath, ok1 := c.Context["basePath"].(string)
    pattern, ok2 := c.Context["pattern"].(string)
	if !ok1 || !ok2 {
		return fmt.Errorf("missing or incorrect type for the context: basePath | pattern")
	}


    watcher, err := fsnotify.NewWatcher()
    if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
    }
    defer watcher.Close()

    // ajouter le dossier racine
    if err := watcher.Add(basePath); err != nil {
		return fmt.Errorf("error adding basePath: %w", err)
    }

	go func() {
		fmt.Println("[Sidecar : WatchFilePattern] Routine launched")
	    for {
	        select {
	        case event, ok := <-watcher.Events:
	            if !ok {
					fmt.Println("[Sidecar : WatchFilePattern] Routine stopped")
	                return
	            }

				if matched, err := filepath.Match(pattern, filepath.Base(event.Name)); err == nil && matched {
					c.Context["filePath"] = event.Name
					c.Context["event.Op"] = event.Op
					fmt.Println("File changed :", event.Name)
					err = c.Clone().CallNext()
					if err != nil {
				        fmt.Println("error on watching file callback : %s", err)
				    }
				}

	        case err, ok := <-watcher.Errors:
	            if !ok {
					fmt.Println("[Sidecar : WatchFilePattern] Routine stopped")
	                return
	            }
	            fmt.Println("watcher error:", err)
	        }
	    }
	}()
	return nil
}


// Call next only if the event is  Write
// From the context
//		- event.Op fsnotify.Op
func IfIsEventOfWriting(c *Callback) error {
	if op, ok := c.Context["event.Op"].(fsnotify.Op) ; ok {
		if op & fsnotify.Write == fsnotify.Write {
			return c.CallNext()
		} else {
			return nil
		}
	} else {
		return fmt.Errorf("missing or incorrect type for the context: event.Op")
	}
}


// Mapping context
// From the context
//		- mountedVolume String
// Add in the context :
//		- basePath String
func MountedVolumeToBasePath(c *Callback) error {
    if _, ok := c.Context["mountedVolume"].(string); ok {
        c.Context["basePath"] = c.Context["mountedVolume"]
    } else {
		return fmt.Errorf("missing or incorrect type for the context: mountedVolume")
	}
    return c.CallNext()
}

// Mapping context
// From the context
//		- mountedVolume String
//		- filePath String
// Add in the context :
//		- filePathInVolume String
func FilePathToFilePathInVolume(c *Callback) error {
	mountedVolume, ok1 := c.Context["mountedVolume"].(string)
	filePath, ok2 := c.Context["filePath"].(string)
	if ok1 && ok2 {
        c.Context["filePathInVolume"] = filePath[len(mountedVolume):]
    } else {
		return fmt.Errorf("missing or incorrect type for the context: mountedVolume | filePath")
	}
    return c.CallNext()
}
