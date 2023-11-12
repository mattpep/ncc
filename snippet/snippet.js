var ncc = function(){
  const author_style = "color: darkgray;";
  const body_style = "color: black;";
  const timestamp_style = "color: lightgray;padding-left: 1em;";
  const comment_style = "margin-bottom: 1.5em;"
  // document.addEventListener("DOMContentLoaded", function(){
  function commentify(comment) {
    var elem = document.createElement("div");
    elem.id = "ncc-comment-" + comment.id;
    elem.classList.add('comment-wrap');
    elem.setAttribute("style", comment_style);

    var meta = document.createElement("div");
    meta.classList.add('meta');

    var author = document.createElement("span");
    author.setAttribute("style", author_style);
    author.innerHTML = comment.display_name;

    var datestamp = document.createElement("span");
    datestamp.setAttribute("style", timestamp_style);
    if (comment.date_time === undefined ) {
      datestamp_str = "just now"
    } else {
      var date_obj = new Date(comment.date_time)
      var dayName = new Array("Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun");
      var monthName = new Array("Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dez");

      /* simulates some of the format strings of strptime() */
      // This functiomn taken from https://www.lzone.de/examples/Javascript%20strptime, Copyright 2012 Lars Windolf <lars.windolf@gmx.de>
      function strptime(format, date) {
        var last = -2;
        var result = "";
        var hour = date.getHours();

        /* Expand aliases */
        format = format.replace(/%D/, "%m/%d/%y");
        format = format.replace(/%R/, "%H:%M");
        format = format.replace(/%T/, "%H:%M:%S");

        /* Note: we fail on strings without format characters */

        while(1) {
          /* find next format char */
          var pos = format.indexOf('%', last + 2);

          if(-1 == pos) {
            /* dump rest of text if no more format chars */
            result += format.substr(last + 2);
            break;
          } else {
            /* dump text after last format code */
            result += format.substr(last + 2, pos - (last + 2));

            /* apply format code */
            formatChar = format.charAt(pos + 1);
            switch(formatChar) {
              case '%':
                result += '%';
                break;
              case 'C':
                result += date.getYear();
                break;
              case 'H':
              case 'k':
                if(hour < 10) result += "0";
                result += hour;
                break;
              case 'M':
                if(date.getMinutes() < 10) result += "0";
                result += date.getMinutes();
                break;
              case 'S':
                if(date.getSeconds() < 10) result += "0";
                result += date.getSeconds();
                break;
              case 'm':
                if(date.getMonth() < 10) result += "0";
                result += date.getMonth();
                break;
              case 'a':
              case 'A':
                result += dayName[date.getDay() - 1];
                break;
              case 'b':
              case 'B':
              case 'h':
                result += monthName[date.getMonth()];
                break;
              case 'Y':
                result += date.getFullYear();
                break;
              case 'd':
              case 'e':
                if(date.getDate() < 10) result += "0";
                result += date.getDate();
                break;
              case 'w':
                result += date.getDay();
                break;
              case 'p':
              case 'P':
                if(hour < 12) {
                  result += "am";
                } else {
                  result += "pm";
                }
                break;
              case 'l':
              case 'I':
                if(hour % 12 < 10) result += "0";
                result += (hour % 12);
                break;
            }
          }
          last = pos;
        }
        return result;
      }

      datestamp_str = strptime("on %d %h %Y at %R", date_obj)
    }

    datestamp.innerHTML = datestamp_str;

    meta.insertAdjacentElement('beforeend', author);
    meta.insertAdjacentElement('beforeend', datestamp);

    var combody = document.createElement("div");
    combody.classList.add('body');
    combody.setAttribute("style", body_style);
    combody.innerHTML = comment.body;

    elem.insertAdjacentElement('beforeend', meta);
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

  function postref_to_path(postref) {
    if (postref == "") {
      postref = window.location.pathname.replaceAll('/','-').substring(1).slice(0,40);
    }
    return postref;
  }

  function show_comments() {
    var comment_url = "{{ .Endpoint }}" + "/comments/" +  postref_to_path("{{ .PostRef }}");
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
    var comment_url = "{{ .Endpoint }}" + "/comments/" +  postref_to_path("{{ .PostRef }}");
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
    var elem = document.getElementById("ncc-comment-00")
    elem.setAttribute("style", elem.getAttribute("style") +"border: 1px dotted black;")
    elem.scrollIntoView();
    document.getElementById("submit_comment").reset()
  }

  function show_comment_form() {
    var ncc = document.getElementById("ncc");
    var form_url = "{{ .Endpoint }}" + "/comments/" +  postref_to_path("{{ .PostRef }}");
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
