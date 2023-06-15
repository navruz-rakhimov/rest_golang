package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/navruz-rakhimov/sarkortelecom/internal/data"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := data.User{
		Login:    r.FormValue("login"),
		Password: r.FormValue("password"),
		Name:     r.FormValue("name"),
	}
	age, err := strconv.Atoi(r.FormValue("age"))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user.Age = age

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	user.Password = string(hashedPassword)

	err = app.models.Users.Insert(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, http.Header{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) authorizeUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &credentials)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByLogin(credentials.Login)
	if err != nil {
		err = app.writeJSON(w, http.StatusUnauthorized, envelope{"credentials": credentials}, nil)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			err = app.writeJSON(w, http.StatusUnauthorized, envelope{"credentials": credentials}, nil)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	accessToken, err := app.GenerateJwtToken(*user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Name:     "SESSTOKEN",
		Value:    accessToken,
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"access_token": accessToken}, http.Header{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) UpdateUserPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	var phone data.Phone
	err := app.readJSON(w, r, &phone)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	_, err = app.models.Phones.Get(phone.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Phones.Update(&phone)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"phone": phone}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) GetUserPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("q")
	phones, err := app.models.Phones.GetAllWithNumber(phoneNumber)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"phones": phones}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserByNameHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	user, err := app.models.Users.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CreateUserPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(int)
	var phone data.Phone

	err := app.readJSON(w, r, &phone)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	phone.UserId = userId

	err = app.models.Phones.Insert(&phone)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"phone": phone}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) DeleteUserPhoneNumberHandler(w http.ResponseWriter, r *http.Request) {
	phoneId := chi.URLParam(r, "id")
	id, err := strconv.Atoi(phoneId)
	if err != nil || id < 1 {
		app.badRequestResponse(w, r, errors.New("invalid id format"))
		return
	}

	err = app.models.Phones.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "phone successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
