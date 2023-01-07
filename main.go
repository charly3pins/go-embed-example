package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

const (
	layoutsDir   = "templates/layouts"
	templatesDir = "templates"
	extension    = "/*.html"
)

var (
	//go:embed templates/* templates/layouts/*
	files     embed.FS
	templates map[string]*template.Template
)

func LoadTemplates() error {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}
	tmplFiles, err := fs.ReadDir(files, templatesDir)
	if err != nil {
		return err
	}
	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(files, templatesDir+"/"+tmpl.Name(), layoutsDir+extension)
		if err != nil {
			return err
		}

		templates[tmpl.Name()] = pt
	}
	return nil
}

const userProfile = "userProfile.html"

func UserProfile(w http.ResponseWriter, r *http.Request) {
	t, ok := templates[userProfile]
	if !ok {
		log.Printf("template %s not found", userProfile)
		return
	}

	data := make(map[string]interface{})
	data["Name"] = "John Doe"
	data["Email"] = "johndoe@email.com"
	data["Address"] = "Fake Street, 123"
	data["PhoneNumber"] = "654123987"

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}

func main() {
	err := LoadTemplates()
	if err != nil {
		log.Fatal(err)
	}
	r := http.NewServeMux()
	r.HandleFunc("/user-profile", UserProfile)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println(err)
	}
}
