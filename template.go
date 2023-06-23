package main

/*
* NOTE:
*	This can't be part of main unless the controller is shifted up to main
*		- or -
*	This and `static.go` are shifted down to their own packages... (including embeded htdocs)
 */

import (
	"html/template"
	"net/http"
	"os"

	"github.com/oxtoacart/bpool"
	log "github.com/sirupsen/logrus"
)

const DISABLE_CSP_ENV_KEY = "DISABLE_CSP"

var (
	tmpl        *template.Template
	bufpool     *bpool.BufferPool
	disable_csp bool = false
)

func init() {
	bufpool = bpool.NewBufferPool(48)
	var err error
	tmpl = template.New("pages")
	tmpl = tmpl.Funcs(template.FuncMap{
		"toJS":   toJS,
		"toHTML": toHTML,
	})
	tmpl, err = tmpl.ParseFS(Templates(), "*.html")
	if err != nil {
		panic(err)
	}
	if os.Getenv(DISABLE_CSP_ENV_KEY) != "" {
		log.Warnf("HTTP Content-Security-Policy disabled by %s != \"\" in ENV", DISABLE_CSP_ENV_KEY)
		disable_csp = true
	} else {
		log.Info("HTTP Content-Security-Policy is enabled.")
	}
}

func RenderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	switch data.(type) {
	case *Stash:
		stash := data.(*Stash)
		stash.TemplateName = name + ".html"
		if disable_csp {
			break
		}
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' ;"+
				"script-src 'self' 'nonce-"+stash.Nonce+"' https://cdnjs.cloudflare.com/ajax/libs/;"+
				"font-src https://cdnjs.cloudflare.com/ajax/libs/;"+
				"img-src 'self' data:;"+
				"frame-src 'self';"+
				"style-src 'self' 'nonce-"+stash.Nonce+"' https://cdnjs.cloudflare.com/ajax/libs/;"+
				"connect-src 'self';",
		)
	default:
	}
	w.Header().Set("Version", Version)

	buf := bufpool.Get()
	err := tmpl.ExecuteTemplate(buf, name+".html", data)
	if err != nil {
		return err
	}
	buf.WriteTo(w)
	bufpool.Put(buf)
	return nil
}

func toJS(s string) template.JS {
	return template.JS(s)
}

func toHTML(s string) template.HTML {
	return template.HTML(s)
}
