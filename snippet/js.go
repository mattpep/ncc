package snippet

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

type templateParams struct {
	Endpoint string
	PostRef  string
}

var endpoint string

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

		function pluralise_comments(count) {
			if (count == 1) {
				return "comment";
			}
			return "comments";

		}

		function insert_single_comment(comment) {
			var ncc = document.getElementById("ncc");
			elem = commentify(comment);
			ncc.insertAdjacentElement('beforeend', elem);
		}

		function show_comments() {
			var comment_url = "{{ .Endpoint }}" + "/comments/" +  "{{ .PostRef }}";
			var xmlHttp = new XMLHttpRequest();
			xmlHttp.open( "GET", comment_url, false );
			xmlHttp.send( null );
			var comment_info = JSON.parse(xmlHttp.responseText);
			var comments = comment_info["Comments"];
			var comment_count = comment_info["Count"];
			var ncc = document.getElementById("ncc");
			if (null === ncc) {
				console.log("ncc - could not find a div with id ncc")
				return
			}
			else {
				var elem = document.createElement("hr");
				ncc.insertAdjacentElement('afterbegin', elem);
				ncc.insertAdjacentHTML('beforeend', '<div class="ncc-banner">'+ comment_count +" " + pluralise_comments(comment_count) +'</div>');
				for (let com = 0; com < comments.length ; com++) {
					insert_single_comment(comments[com])
				}
			}
		}
		function submitComment(form) {
			var comment_url = "{{ .Endpoint }}" + "/comments/" +  "{{ .PostRef }}";
			var commentData = new FormData(form);

			var c = Object.fromEntries(commentData);

			const response = fetch(comment_url, { // await fetch(url, ... )

				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify(c),

			});
			c.id = '00';
			insert_single_comment(c);
			document.getElementById("submit_comment").reset()
		}

		function show_comment_form() {
			var ncc = document.getElementById("ncc");
			var form_url = "{{ .Endpoint }}" + "/comments/" +  "{{ .PostRef }}";
			var comment_form = '<form id="submit_comment">' +
			'<div class="field" style="padding-bottom: 1em;"><label for="display_name" style="display: block;">Name</label><input type="text" size=30 name="display_name" /></div>' +
			'<div class="field" style="padding-bottom: 0.6em;"><label for="body" style="display: block;">Comment</label><textarea cols=45 rows=6 name="body" placeholder="Enter comment here…"></textarea></div>' +
			'<div class="field" style="padding-bottom: 1em;"><input type="submit" value="submit" name="submit"></div>' +
			'</form>';
			ncc.insertAdjacentHTML('beforeend', comment_form);
		}
		show_comments();
		show_comment_form();
		document.getElementById("submit_comment").addEventListener("submit", function (e) {
			e.preventDefault();
			submitComment(e.target);
		});
	};
	`
	buf := &bytes.Buffer{}
	tmpl, err := template.New("tmpl").Parse(snippet)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(buf, templateParams{Endpoint: endpoint, PostRef: post_ref})
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ServeWebsiteInsert(w http.ResponseWriter, r *http.Request) {
	postref := r.URL.Query().Get("postref")

	ext_endpoint, present := os.LookupEnv("EXT_ENDPOINT")
	if present {
		endpoint = ext_endpoint
	} else {
		endpoint = "http://localhost:8080"
	}

	log.Printf("serving js insert: post ref is >%v<", postref)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	io.WriteString(w, snippet)
}
