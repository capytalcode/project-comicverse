package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"forge.capytal.company/capytalcode/project-comicverse/internals/randstr"
	"forge.capytal.company/capytalcode/project-comicverse/service"
	"forge.capytal.company/loreddev/x/smalltrip/exception"
)

func (router *router) pages(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	// TODO: Check if project exists
	id := r.PathValue("ID")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "ID" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	pageID := r.PathValue("PageID")

	switch getMethod(r) {
	case http.MethodGet, http.MethodHead:
		if pageID == "" {
			exception.
				BadRequest(fmt.Errorf(`a valid path value of "PageID" must be provided`)).
				ServeHTTP(w, r)
			return
		}

		router.getPage(w, r)

	case http.MethodPost:
		router.addPage(w, r)

	case http.MethodDelete:
		if pageID == "" {
			exception.
				BadRequest(fmt.Errorf(`a valid path value of "PageID" must be provided`)).
				ServeHTTP(w, r)
			return
		}

		router.deletePage(w, r)

	default:
		exception.
			MethodNotAllowed([]string{
				http.MethodGet,
				http.MethodHead,
				http.MethodPost,
				http.MethodDelete,
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

	pageID := r.PathValue("PageID")
	router.assert.NotZero(pageID, "This method should be used after the path values are checked")

	page, err := router.service.GetPage(id, pageID)
	if errors.Is(err, service.ErrPageNotExists) {
		exception.NotFound(exception.WithError(err)).ServeHTTP(w, r)
		return
	}
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	if i, ok := page.Image.(io.WriterTo); ok {
		_, err = i.WriteTo(w)
	} else {
		_, err = io.Copy(w, page.Image)
	}

	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
}

func (router *router) deletePage(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	id := r.PathValue("ID")
	router.assert.NotZero(id, "This method should be used after the path values are checked")

	pageID := r.PathValue("PageID")
	router.assert.NotZero(pageID, "This method should be used after the path values are checked")

	err := router.service.DeletePage(id, pageID)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%s/", id), http.StatusSeeOther)
}

func (router *router) interactions(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)

	// TODO: Check if the project exists
	id := r.PathValue("ID")
	if id == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "ID" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	// TODO: Check if page exists
	pageID := r.PathValue("PageID")
	if pageID == "" {
		exception.
			BadRequest(fmt.Errorf(`a valid path value of "PageID" must be provided`)).
			ServeHTTP(w, r)
		return
	}

	interactionID := r.PathValue("InteractionID")

	switch getMethod(r) {
	case http.MethodPost:
		router.addInteraction(w, r)

	default:
		exception.
			MethodNotAllowed([]string{
				http.MethodPost,
			}).
			ServeHTTP(w, r)
	}
}

func (router *router) addInteraction(w http.ResponseWriter, r *http.Request) {
	router.assert.NotNil(w)
	router.assert.NotNil(r)
	router.assert.NotNil(router.service)

	id := r.PathValue("ID")
	router.assert.NotZero(id, "This method should be used after the path values are checked")

	pageID := r.PathValue("PageID")
	router.assert.NotZero(pageID, "This method should be used after the path values are checked")

	// TODO: Methods to manipulate interactions, instead of router need to do this logic
	page, err := router.service.GetPage(id, pageID)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}
	page.Image = nil // HACK: Prevent image update on S3

	x, err := strconv.ParseUint(r.FormValue("x"), 10, 0)
	if err != nil {
		exception.
			BadRequest(errors.Join(errors.New(`value "x" should be a valid non-negative integer`), err)).
			ServeHTTP(w, r)
		return
	}
	y, err := strconv.ParseUint(r.FormValue("y"), 10, 0)
	if err != nil {
		exception.
			BadRequest(errors.Join(errors.New(`value "y" should be a valid non-negative integer`), err)).
			ServeHTTP(w, r)
		return
	}

	link := r.FormValue("link")
	if link == "" {
		exception.BadRequest(errors.New(`missing parameter "link" in request`)).ServeHTTP(w, r)
		return
	}

	intID, err := randstr.NewHex(6)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	page.Interactions[intID] = service.PageInteraction{
		X:   uint16(x),
		Y:   uint16(y),
		URL: link,
	}

	err = router.service.UpdatePage(id, page)
	if err != nil {
		exception.InternalServerError(err).ServeHTTP(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/projects/%s/", id), http.StatusSeeOther)
}

