package sessions

type SupportedLanguage struct {
	Name       string
	FolderPath string
	FileName   string
	BuildCmd   string
	RunCmd     string
	Code       string
}

var (
	supportedLanguages = []SupportedLanguage{
		{
			Name:       "go",
			FolderPath: "cmd",
			FileName:   "main.go",
			BuildCmd:   "go build -o cmd/main cmd/main.go",
			RunCmd:     "cmd/main",
			Code: `
package main

import "fmt"

func main() {
    fmt.Println("Hello world from Go.")
}
`,
		},
		{
			Name:       "python",
			FolderPath: "cmd",
			FileName:   "main.py",
			BuildCmd:   "python3 -m compileall cmd/main.py",
			RunCmd:     "python3 -u cmd/main.py",
			Code: `
print("Hello world from Python.")
`,
		},
		{
			Name:       "typescript",
			FolderPath: "cmd",
			FileName:   "main.ts",
			BuildCmd:   "tsc --outFile cmd/main.js cmd/main.ts",
			RunCmd:     "node cmd/main.js",
			Code: `
console.log('Hello world from TypeScript.');
`,
		},
		{
			Name:       "c",
			FolderPath: "cmd",
			FileName:   "main.c",
			BuildCmd:   "gcc -o cmd/main cmd/main.c",
			RunCmd:     "cmd/main",
			Code: `
#include <stdio.h>

int main(int argc, char *argv[]) {
	printf("Hello world from C.\n");

	return 0;
}
`,
		},
		{
			Name:       "rust",
			FolderPath: "cmd",
			FileName:   "main.rs",
			BuildCmd:   "rustc -o cmd/main cmd/main.rs",
			RunCmd:     "cmd/main",
			Code: `
fn main() {
    println!("Hello world from Rust.");
}
`,
		},
		{
			Name:       "java",
			FolderPath: "cmd",
			FileName:   "Main.java",
			BuildCmd:   "javac cmd/Main.java",
			RunCmd:     "java --class-path cmd Main",
			Code: `
class Main {
    public static void main(String[] args) {
        System.out.println("Hello world from Java.");
    }
}
`,
		},
	}

	supportedLanguageByName map[string]SupportedLanguage
)

func init() {
	supportedLanguageByName = make(map[string]SupportedLanguage)
	for _, supportedLanguage := range supportedLanguages {
		supportedLanguageByName[supportedLanguage.Name] = supportedLanguage
	}
}
