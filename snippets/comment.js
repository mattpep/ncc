function flagComment (commentId) { // eslint-disable-line no-unused-vars
  const commentURL = '{{ .Endpoint }}' + '/flag/' + commentId

  fetch(commentURL, {
    method: 'POST',
    body: ''
  }).then((response) => {
    if (response.status === 204) {
      const commentEl = document.querySelector('#ncc #ncc-comment-' + commentId)
      commentEl.innerHTML = '<span style="color: green; font-size: 0.6em;">Comment flagged</span>'
      commentEl.scrollIntoView()
    } else {
      const actionEl = document.querySelector('#ncc #ncc-comment-' + commentId + ' div.actions form')
      actionEl.insertAdjacentHTML('beforeend', '<span style="color: red; font-size: 0.6em;">Could not flag comment</span>')
    }
  }).catch((err) => {
    console.log('Error: ' + err)
  })
}

document.addEventListener('DOMContentLoaded', function () {
  const authorStyle = 'color: darkgray;'
  const bodyStyle = 'color: black;'
  const timestampStyle = 'color: lightgray;padding-left: 1em;'
  const commentStyle = 'margin-bottom: 1.5em;'
  const actionStyle = 'font-weight: bold; text-size: 0.8em;'
  function commentify (comment) {
    const elem = document.createElement('div')
    let datestampStr = ''
    elem.id = 'ncc-comment-' + comment.id
    elem.classList.add('comment-wrap')
    elem.setAttribute('style', commentStyle)

    const meta = document.createElement('div')
    meta.classList.add('meta')

    const author = document.createElement('span')
    author.setAttribute('style', authorStyle)
    author.innerHTML = comment.display_name

    const datestamp = document.createElement('span')
    datestamp.setAttribute('style', timestampStyle)
    if (comment.date_time === undefined) {
      datestampStr = 'just now'
    } else {
      const dateObj = new Date(comment.date_time)
      const dayName = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']
      const monthName = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dez']

      /* simulates some of the format strings of strptime() */
      // This function taken from https://www.lzone.de/examples/Javascript%20strptime, Copyright 2012 Lars Windolf <lars.windolf@gmx.de>
      function strptime (format, date) {
        let last = -2
        let result = ''
        let formatChar
        const hour = date.getHours()

        /* Expand aliases */
        format = format.replace(/%D/, '%m/%d/%y')
        format = format.replace(/%R/, '%H:%M')
        format = format.replace(/%T/, '%H:%M:%S')

        /* Note: we fail on strings without format characters */

        while (1) {
          /* find next format char */
          const pos = format.indexOf('%', last + 2)

          if (pos === -1) {
            /* dump rest of text if no more format chars */
            result += format.substr(last + 2)
            break
          } else {
            /* dump text after last format code */
            result += format.substr(last + 2, pos - (last + 2))

            /* apply format code */
            formatChar = format.charAt(pos + 1)
            switch (formatChar) {
              case '%':
                result += '%'
                break
              case 'C':
                result += date.getYear()
                break
              case 'H':
              case 'k':
                if (hour < 10) result += '0'
                result += hour
                break
              case 'M':
                if (date.getMinutes() < 10) result += '0'
                result += date.getMinutes()
                break
              case 'S':
                if (date.getSeconds() < 10) result += '0'
                result += date.getSeconds()
                break
              case 'm':
                if (date.getMonth() < 10) result += '0'
                result += date.getMonth()
                break
              case 'a':
              case 'A':
                result += dayName[date.getDay() - 1]
                break
              case 'b':
              case 'B':
              case 'h':
                result += monthName[date.getMonth()]
                break
              case 'Y':
                result += date.getFullYear()
                break
              case 'd':
              case 'e':
                if (date.getDate() < 10) result += '0'
                result += date.getDate()
                break
              case 'w':
                result += date.getDay()
                break
              case 'p':
              case 'P':
                if (hour < 12) {
                  result += 'am'
                } else {
                  result += 'pm'
                }
                break
              case 'l':
              case 'I':
                if (hour % 12 < 10) result += '0'
                result += (hour % 12)
                break
            }
          }
          last = pos
        }
        return result
      }

      datestampStr = strptime('on %d %h %Y at %R', dateObj)
    }

    datestamp.innerHTML = datestampStr

    meta.insertAdjacentElement('beforeend', author)
    meta.insertAdjacentElement('beforeend', datestamp)

    const combody = document.createElement('div')
    combody.classList.add('body')
    combody.setAttribute('style', bodyStyle)
    combody.innerHTML = comment.body

    elem.insertAdjacentElement('beforeend', meta)
    elem.insertAdjacentElement('beforeend', combody)

    const flagUrl = '{{ .Endpoint }}' + '/flag/' + comment.id
    const actions = '<div class="actions" style="' + actionStyle + '"><form id="flag-' + comment.id + '" method="post" action="' + flagUrl + '"> <input type="hidden" name="name" value="value" /> <a href="#" onclick="event.preventDefault();flagComment(\'' + comment.id + '\')">Flag</a> </form>'
    elem.insertAdjacentHTML('beforeend', actions)
    return elem
  }

  function pluraliseComments (count) {
    if (count === 1) {
      return 'comment'
    }
    return 'comments'
  }

  function insertSingleComment (comment) {
    const ncc = document.getElementById('ncc')
    const elem = commentify(comment)
    ncc.insertAdjacentElement('beforeend', elem)
  }

  function urlToPostref (postref) {
    if (postref === '') {
      postref = window.location.pathname.replaceAll('/', '-').substring(1).slice(0, 40)
    }
    return postref
  }

  function showComments () {
    const commentURL = '{{ .Endpoint }}' + '/comments/' + urlToPostref('{{ .PostRef }}')
    fetch(commentURL, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json'
      }
    }).then((response) => {
      if (response.status !== 200) {
        return
      }
      return response.json()
    }).then((commentInfo) => {
      const comments = commentInfo.comments
      const ncc = document.getElementById('ncc')
      if (ncc === null) {
        console.log('ncc - could not find a div with id ncc')
      } else {
        const elem = document.createElement('hr')
        ncc.insertAdjacentElement('afterbegin', elem)
        ncc.insertAdjacentHTML('beforeend', '<div class="ncc-banner">' + commentInfo.count + ' ' + pluraliseComments(commentInfo.count) + '</div>')
        for (let com = 0; com < comments.length; com++) {
          insertSingleComment(comments[com])
        }
      }
    })
  }

  function submitComment (form) {
    const commentURL = '{{ .Endpoint }}' + '/comments/' + urlToPostref('{{ .PostRef }}')
    const commentData = new FormData(form)

    const c = Object.fromEntries(commentData)

    fetch(commentURL, {

      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(c)

    }).then((response) => {
      if (response.status !== 204) {
        document.querySelector('#ncc .error').innerHTML = 'ERROR'
        return
      }
      document.querySelector('#ncc .error').innerHTML = ''
      c.id = '00'
      insertSingleComment(c)
      const elem = document.getElementById('ncc-comment-00')
      elem.setAttribute('style', elem.getAttribute('style') + 'border: 1px dotted black;')
      elem.scrollIntoView()
      document.getElementById('submit_comment').reset()
    }).catch(function (err) {
      if (err instanceof TypeError) {
        document.querySelector('#ncc .error').innerHTML = 'ERROR - Is ncc offline?'
      } else {
        document.querySelector('#ncc .error').innerHTML = 'Unknown error'
      }
    })
  }

  function showCommentForm () {
    const ncc = document.getElementById('ncc')
    const commentForm = '<form id="submit_comment">' +
      '<div class="field" style="padding-bottom: 1em;"><label for="display_name" style="display: block;">Name</label><input type="text" size=30 name="display_name" /></div>' +
      '<div class="field" style="padding-bottom: 0.6em;"><label for="body" style="display: block;">Comment</label><textarea cols=45 rows=6 name="body" placeholder="Enter comment here…"></textarea></div>' +
      '<div class="field" style="padding-bottom: 1em;"><input type="submit" value="submit" name="submit"><span class="error" style="color: red; text-size: 0.6em"></span></div>' +
      '</form>'
    ncc.insertAdjacentHTML('beforeend', commentForm)
  }
  showComments()
  showCommentForm()
  document.getElementById('submit_comment').addEventListener('submit', function (e) {
    e.preventDefault()
    submitComment(e.target)
  })
})
