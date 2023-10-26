# No-cookie comments

Users often hate cookies due to the pop-ups. This is a comment system designed
to work without the use of cookies.


It's very much an experiment, but you're welcome to use it if you want.


# How it works

Website readers post comments (via a JavaScript snippet) which calls to an API.
When subsequent users reload the page, the Javascript will pull down the
comments for the post and insert them into a div element with the id ncc.

# Setup

The main engine is an API which runs perpetually. It receives comments
submitted by users, and also requests to read comments for a given post, and
comment counts (for post index pages). This part works without cookies or
authentication - hence the name ncc.

# Configuration

In the spirit of 12factor.net, this application uses environment variables for
configuration. The possible settings are:

* `PORT` - which local port to run on
* `DATABASE_URL` - how to connect to the database
* `EXT_ENDOINT` - the prefix of the public-facing URL of this service (i.e.
  outside of any loadbalancer or container which might be in use)

# Some caveats

Although this system does not generate, store, or make use of cookies for
public users. it does still use cookies for website admins and moderators, in
order that they can perform authenticated operations such as moderation.

But of note, because personal names are still collected then GDPR and related
legislation will still apply.

Lastly, I am not a legal expert so you should not make inferences about your
obligations based on the information I provide here.
