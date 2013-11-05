infoauth
========

OAuth gateway and API for infolab


Installation
------------

Obtain and install infoauth with `go get`:

	go get github.com/axelmagn/infoauth

There are a few variables that need to be set in the configuration file.  See example in `config/settings.ecfg`.  This is also the default config file path, which can be specified with a command line option.  Infoauth uses [envcfg](http://github.com/axelmagn/envcfg) as a settings parser, if you need to write your own settings.

once you have your configuration file set, run the app with:

	infoauth

or for a custom config path:

	infoauth -config="<CONFIG_FILE>"

This will begin serving the infoauth endpoints on the specified port.


Redirect API
------------

Infoauth's primary function is to expose both REST API and Redirection endpoints for obtaining access tokens via oauth.  It exposes the following endpoints:

	/google/url		# Create and print an oauthauthorization url 
					# which a user can follow to authorize our app via
					# their google account

	/linkedin/url	# Create and print an oauth authorization url 
					# which a user can follow to authorize our app via
					# their linkedin account

	/google/		# Same as above, but automatically redirects to
					# the url

	/linkedin/		# Same as above, but automatically redirects to
					# the url

	/				# Where the user is redirected to after
					# authenticating via openauth.  Takes
					# authorization code and exchanges it for an
					# access token.  If a redirect url is specified,
					# it redirects after this stage

When requesting a url, a redirect can be specified with the `redirect=<URL>` query parameter. The user will then be redirected to the url after the handshake is complete, with the access token, success status, and service name as query parameters. In cases of unsuccessful handshakes, the user is not redirected. At the time of writing, redirects do not occur in the case of bad inputs or server errors, so outgoing redirects will always have success=true.

Redirect query parameters:

	access_token	# the access token
	service			# service name.  Takes values [ GOOGLE | LINKEDIN ]
	success			# whether the exchange was successful
	error (planned) # in cases of unsuccessful queries, error description
