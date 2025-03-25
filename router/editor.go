package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
)

func (router *router) pages(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	id := r.PathValue("ID")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "ID" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	switch getMethod(r) {
	case http.MethodGet, http.MethodHead:
		imgID := r.PathValue("PageID")
		if imgID == "" {
			exception.
				BadRequest(fmt.Errorf(`a valid path value of "PageID" must be provided`)).
				ServeHTTP(w, r)
			return
		}

		router.getPage(w, r)

	case http.MethodPost:
		router.addPage(w, r)

	default:
		exception.
			MethodNotAllowed([]string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
			}).
			ServeHTTP(w, r)
	}
}

func (router *router) addPage(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	id := r.PathValue("ID")
	router.assert.NotZero(id, "This method should be used after the path values are checked")

	img, _, err := r.FormFile("image")
	if err != nil {
		// TODO: Handle if the file is bigger than allowed by ParseForm (10mb)
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	err = router.service.AddPage(id, img)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%s/", id), http.StatusSeeOther)
}

func (router *router) getPage(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	id := r.PathValue("ID")
	router.assert.NotZero(id, "This method should be used after the path values are checked")

	imgID := r.PathValue("PageID")
	router.assert.NotZero(imgID, "This method should be used after the path values are checked")

	img, err := router.service.GetPage(id, imgID)
	if errors.Is(err, service.ErrPageNotExists) {
		exception.NotFound(exception.WithError(err)).ServeHTTP(w, r)
		return
	}
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	if i, ok := img.(io.WriterTo); ok {
		_, err = i.WriteTo(w)
	} else {
		_, err = io.Copy(w, img)
	}

	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
}

