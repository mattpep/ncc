function pluralise_comments (count) {
  if (count == 1) {
    return 'comment'
  }
  return 'comments'
}

function url_to_postref (url) {
  var postref = new URL(url).pathname.replaceAll('/', '-').substring(1).slice(0, 40)
  return postref
}

function add_count (a_tag) {
  var postref = url_to_postref(a_tag.href)
  var count_url = '{{ .Endpoint }}' + '/commentcount/' + postref
  var response = fetch(count_url, {

    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }).then((response) => {
    if (response.status != 200) {
      console.log("Unexpected status when getting comments count: " + response.status)
      return
    }
    return response.json();
  }).then((json) => {
    a_tag.innerHTML = json.count + " " + pluralise_comments(json.count)
  })

}

document.addEventListener('DOMContentLoaded', function() {
// function insert_counts() {
  var ncc = document.getElementById('ncc')
  if (ncc == null) {
    console.log('ncc - could not find a tag with id ncc')
    return
  }

  // a_tags = document.querySelectorAll('#ncc a.ncc')
  a_tags = document.getElementsByClassName('ncc')
  for (let i = 0; i < a_tags.length; i++) {
    add_count(a_tags[i])
  }
})

