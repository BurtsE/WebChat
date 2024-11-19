package main

import (
	"WebChat/trace"
	"flag"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates",
			t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}

	if userData, err := getUserData(r); err == nil {
		data["UserData"] = userData
	}

	err := t.templ.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	configureOauth2()

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/", http.RedirectHandler("/chat", http.StatusFound))
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/room", r)
	http.HandleFunc("/auth/{action}/{provider}/", loginHandler)

	// Handle static css and js files
	//http.Handle("/assets/", http.StripPrefix("/assets",
	//	http.FileServer(http.Dir("/path/to/assets/"))))

	go r.run()

	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func configureOauth2() {
	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, "http://localhost:8080/auth/callback/google/"),
		// facebook.New(clientID, clientSecret, "http://localhost:8080/auth/callback/facebook/"),
		github.New(githubClientID, githubClientSecret, "http://localhost:8080/auth/callback/github/"),
	)
	gothic.GetProviderName = getProviderName
}
