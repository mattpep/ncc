function pluraliseComments (count) {
  if (count === 1) {
    return 'comment'
  }
  return 'comments'
}

function urlToPostref (url) {
  const postref = new URL(url).pathname.replaceAll('/', '-').substring(1).slice(0, 40)
  return postref
}

function updateCountInfo (blogref) {
  const countInfoURL = '{{ .Endpoint }}' + '/v2/blog/' + blogref + '/commentcounts'
  fetch(countInfoURL, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }).then((response) => {
    if (response.status !== 200) {
      console.log('Unexpected status when getting comments count: ' + response.status)
      return
    }
    return response.json()
  }).then((json) => {
    var dict = {}
    json.counts.forEach((el, _) => dict[el.postref] = el.count);
    const tags = document.getElementsByClassName('ncc')
    for (let i = 0; i < tags.length; i++) {
      const postref = urlToPostref(tags[i].href)
      if ( postref in dict) {
	tags[i].innerHTML = dict[postref] + ' ' + pluraliseComments(dict[postref])
      }
    }
  })

}

function addCount (tag) {
  const postref = urlToPostref(tag.href)
  const countURL = '{{ .Endpoint }}' + '/commentcount/' + postref
  fetch(countURL, {

    method: 'GET',
    headers: {
      'Content-Type': 'application/json'
    }
  }).then((response) => {
    if (response.status !== 200) {
      console.log('Unexpected status when getting comments count: ' + response.status)
      return
    }
    return response.json()
  }).then((json) => {
    tag.innerHTML = json.count + ' ' + pluraliseComments(json.count)
  })
}

document.addEventListener('DOMContentLoaded', function () {
  const ncc = document.getElementById('ncc')
  if (ncc == null) {
    console.log('ncc - could not find a tag with id ncc')
    return
  }
  updateCountInfo("{{ .BlogRef }}")
})
