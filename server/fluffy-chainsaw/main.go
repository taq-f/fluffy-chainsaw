package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

//go:generate go-assets-builder --strip-prefix="/assets" -o assets.go  assets

type Pong struct {
	Status int
	Result string
}

type Todo struct {
	ID   string `json:"id"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func isGitRepository(repository string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = repository
	output, _ := cmd.Output()
	return strings.TrimSpace(string(output)) == "true"
}

func getFiles(repository string) []string {
	files := []string{}

	cmd := exec.Command("git", "ls-files")
	cmd.Dir = repository
	output, _ := cmd.Output()
	for _, file := range strings.Split(string(output), "\n") {
		if file == "" {
			continue
		}
		files = append(files, file)
	}

	cmd = exec.Command("git", "ls-files", "--others", "--exclude-standard")
	cmd.Dir = repository
	output, _ = cmd.Output()
	for _, file := range strings.Split(string(output), "\n") {
		if file == "" {
			continue
		}
		files = append(files, file)
	}

	return files
}

func getExtensions(files []string) []string {
	extensions := []string{}

	for _, file := range files {
		extensions = append(extensions, path.Ext(file))
	}

	// remove duplicates
	m := make(map[string]bool)
	uniq := []string{}

	for _, ext := range extensions {
		if !m[ext] {
			m[ext] = true
			uniq = append(uniq, ext)
		}
	}

	return uniq
}

func findTodos(repository string, extension string) []Todo {
	var regex string
	var fileBlob string

	switch extension {
	case ".js":
		regex = "//\\s*todo"
	case ".ts":
		regex = "//\\s*todo"
	case ".py":
		regex = "#\\s*todo"
	case ".java":
		regex = "//\\s*todo"
	case ".vue":
		regex = "//\\s*todo"
	case ".go":
		regex = "//\\s*todo"
	default:
		// unsupported file type
		return []Todo{}
	}
	fileBlob = "*" + extension

	cmd := exec.Command("git", "grep", "-e", regex, "-i", "--full-name", "--line-number", "--untracked", "--no-color", "-w", "-I", "--null", "--", fileBlob)
	cmd.Dir = repository
	output, _ := cmd.Output()

	lines := strings.Split(string(output), "\n")
	todos := []Todo{}

	for _, line := range lines {
		components := strings.Split(line, "\x00")
		if len(components) != 3 {
			continue
		}

		todo := Todo{}
		todo.ID = components[0] + components[1]
		todo.File = components[0]
		todo.Line, _ = strconv.Atoi(components[1])
		todo.Text = strings.TrimSpace(components[2])

		todos = append(todos, todo)
	}

	return todos
}

func main() {
	flag.Parse()
	repository := flag.Arg(0)

	fmt.Println()
	fmt.Println("Initializing application...")

	if !isGitRepository(repository) {
		log.Fatal("FATAL! The provided path doesn't seem to be a git repository!")
	}

	fmt.Printf("Repository: %v\n", repository)

	http.HandleFunc("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		files := getFiles(repository)
		extensions := getExtensions(files)
		todos := []Todo{}

		for _, ext := range extensions {
			for _, todo := range findTodos(repository, ext) {
				todos = append(todos, todo)
			}
		}

		fmt.Printf("  [%v] API request %v. %v TODOs found\n", time.Now().Format("01/02 15:04:05"), r.URL.Path, len(todos))

		res, _ := json.Marshal(todos)

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		pong := Pong{http.StatusOK, "OK"}

		res, _ := json.Marshal(pong)

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	})

	http.Handle("/", http.FileServer(Assets))

	fmt.Printf("\n\n  Serving at %v !\n\n", "localhost:8123")
	log.Fatal(http.ListenAndServe("localhost:8123", nil))
}
