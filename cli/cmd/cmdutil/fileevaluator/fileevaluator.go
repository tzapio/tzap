package fileevaluator

import (
	"os"
	"path"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/tzapio/tzap/internal/logging/tl"
)

type FileEvaluator struct {
	includeMatcher *ignore.GitIgnore
	excludeMatcher *ignore.GitIgnore
}

const BaseTzapIgnore = `# Tzap ignore file. Add extra files like test folders, or other files that interfere with search (embeddings) quality.
node_modules
.env
`

const BaseTzapInclude = `# Common languages. Example, remove .js if .js files are only compiled bundles.
*.js
*.tsx
*.ts
*.py
*.go
*.cpp
*.css
*.html
*.java
*.c
*.h
*.hpp
*.rs
*.rb
*.php
*.cob
*.COB
*.cbl
*.CBL
*.c
*.cpp
*.cc
*.cxx
*.hpp
*.h
*.lua
*.cs
*.rb
*.hlsl
*.swift
*.kts
*.kt
*.dart
*.groovy
*.gvy
*.gy
*.gsh
*.scala
*.sc
*.pl
*.pm
*.t
*.r
*.m
*.mm
*.f
*.f90
*.f95
*.for
*.hs
*.asm
*.s
*.m
*.v
*.vhdl
*.pro
*.lisp
*.cl
*.el
*.scm
*.ss
*.rkt
*.il
*.fs
*.fsx
*.fsi
*.ml
*.mli
*.purs
*.ex
*.exs
*.elm
*.erl
*.hrl
*.clj
*.cljs
*.cljc
*.edn
*.vb
*.vbs
*.bas
*.ada
*.adb
*.ads
*.pascal
*.pas
*.d
*.nim
*.jl
*.cr
`

var baseExcludePatterns = []string{".git", ".DS_Store", "desktop.ini"}

func New(baseDir string) (*FileEvaluator, error) {
	gitIgnorePath := path.Join(baseDir, ".gitignore")
	tzapIgnorePath := path.Join(baseDir, ".tzapignore")
	tzapIncludePath := path.Join(baseDir, ".tzapinclude")
	ignoreFiles := []string{tzapIgnorePath}
	tl.Logger.Println("gitIgnorePath", gitIgnorePath)
	tl.Logger.Println("tzapIgnorePath", tzapIgnorePath)
	tl.Logger.Println("tzapIncludePath", tzapIncludePath)
	if _, err := os.Stat(gitIgnorePath); err == nil {
		ignoreFiles = append(ignoreFiles, gitIgnorePath)
	}
	var excludePatterns []string
	excludePatternsFromFile, err := ReadFilterPatternFiles(ignoreFiles...)
	if err != nil {
		baseTzapIgnore, _ := ReadPatternString(BaseTzapIgnore)
		excludePatterns = append(baseExcludePatterns, baseTzapIgnore...)
		println("Using base tzapignore")
	} else {
		excludePatterns = append(baseExcludePatterns, excludePatternsFromFile...)
	}
	var includePatterns []string
	includePatternsFromFile, err := ReadFilterPatternFiles(tzapIncludePath)
	if err != nil {
		baseTzapInclude, _ := ReadPatternString(BaseTzapInclude)
		includePatterns = append(baseExcludePatterns, baseTzapInclude...)
		println("Using base tzapinclude")
	} else {
		includePatterns = append(includePatterns, includePatternsFromFile...)
	}
	return NewWithPatterns(excludePatterns, includePatterns), nil
}
func NewWithBasePatterns() *FileEvaluator {
	baseTzapIgnore, _ := ReadPatternString(BaseTzapIgnore)
	excludePatterns := append(baseExcludePatterns, baseTzapIgnore...)
	baseTzapInclude, _ := ReadPatternString(BaseTzapInclude)
	includePatterns := append(baseExcludePatterns, baseTzapInclude...)
	return NewWithPatterns(excludePatterns, includePatterns)
}
func NewWithPatterns(excludePatterns []string, includePatterns []string) *FileEvaluator {
	excludeMatcher := ignore.CompileIgnoreLines(excludePatterns...)
	includeMatcher := ignore.CompileIgnoreLines(includePatterns...)
	return &FileEvaluator{excludeMatcher: excludeMatcher, includeMatcher: includeMatcher}
}
func (e *FileEvaluator) ShouldKeepPath(path string) bool {
	exclude := e.excludeMatcher.MatchesPath(path)
	include := e.includeMatcher.MatchesPath(path)
	return include && !exclude
}
func (e *FileEvaluator) ShouldTraverseDir(path string) bool {
	return !e.excludeMatcher.MatchesPath(path)
}
