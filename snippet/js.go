package snippet

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"text/template"
)

type templateParams struct {
	Endpoint string
	PostRef  string
}

const Endpoint = "http://localhost:8080"

func websiteInsert(post_ref string) (string, error) {
	snippet := `var ncc = function(){
		const author_style = "color: darkgray;";
		const body_style = "color: black;";
		const comment_style = "margin-bottom: 1.5em;"
// document.addEventListener("DOMContentLoaded", function(){
		function commentify(comment) {
			var elem = document.createElement("div");
			elem.id = "ncc-comment-" + comment.id;
			elem.setAttribute("style", comment_style);
			
			var author = document.createElement("div");
			author.classList.add('author');
			author.setAttribute("style", author_style);
			author.innerHTML = comment.display_name;

			var combody = document.createElement("div");
			combody.classList.add('body');
			combody.setAttribute("style", body_style);
			combody.innerHTML = comment.body;

			elem.insertAdjacentElement('beforeend', author);
			elem.insertAdjacentElement('beforeend', combody);
			return elem;
		}
		var comment_url = "{{ .Endpoint }}" + "/comments/" +  "{{ .PostRef }}";
		var xmlHttp = new XMLHttpRequest();
		xmlHttp.open( "GET", comment_url, false );
		xmlHttp.send( null );
		var comment_info = xmlHttp.responseText;
		var ncc = document.getElementById("ncc");
		if (null === ncc) {
			console.log("ncc - could not find a div with id ncc")
			return
		}
		else {
			var elem = document.createElement("hr");
			ncc.insertAdjacentElement('afterbegin', elem);
			ncc.insertAdjacentHTML('beforeend', '<div class="ncc-banner">Comments</div>');
			var comments = JSON.parse(comment_info)["Comments"];
			for (let com = 0; com < comments.length ; com++) {
				elem = commentify(comments[com]);
				ncc.insertAdjacentElement('beforeend', elem);
			}
		}

		//})
	}
	`
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: Endpoint, PostRef: post_ref})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ServeWebsiteInsert(w http.ResponseWriter, r *http.Request) {
	postref := r.URL.Query().Get("postref")

	log.Printf("post ref is >%v<", postref)
	if postref == "" {
		log.Println("{\"status\":\"error\",\"message\":\"postref must be provided\"}")
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"status\":\"error\",\"message\":\"postref must be provided\"}")
		return
	}
	snippet, err := websiteInsert(postref)
	if err != nil {
		log.Printf("Error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "{\"status\":\"error\",\"message\":\"Cannot serve snippet\"}")
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	io.WriteString(w, snippet)
}
