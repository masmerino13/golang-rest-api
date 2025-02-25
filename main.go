package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"lens.com/m/v2/controllers"
	"lens.com/m/v2/migrations"
	"lens.com/m/v2/models"
	"lens.com/m/v2/templates"
	"lens.com/m/v2/views"
)

func main() {
	port := "3002"
	fmt.Println("Starting server at ", port)

	// Setup database
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")

	if err != nil {
		panic(err)
	}

	// Setup Services
	usersService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}

	// Setup middleware
	umw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// TODO: Set this ke in a secure way
	csrfKey := "gTHGA14aqlIiv6F6FhTMliv0AoW3ju71lOvnRJ9PWzZ8ML8aHVyXoziTsO3pkfDc05Y7AFH3Y5IPARnU7mtPzVHtju07wWUSf0Vi"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup controllers
	usersC := controllers.Users{
		// NOTE: this way we set the model to the controller
		UserService:    &usersService,
		SessionService: &sessionService,
	}
	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "singup.gohtml", "tailwind.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "singin.gohtml", "tailwind.gohtml"))

	// Router and routes
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(csrfMw)
	r.Use(umw.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", usersC.New)
	r.Post("/users", usersC.Create)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.Auth)
	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})
	r.Post("/users/signout", usersC.SingOut)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found blah", http.StatusNotFound)
	})

	fmt.Println("Running at ", port)

	http.ListenAndServe(":"+port, csrfMw(umw.SetUser(r)))
}
