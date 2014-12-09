# .hiv Domain Status

[![Travis](https://travis-ci.org/dothiv/hiv-domain-status.svg?branch=master)](https://travis-ci.org/dothiv/hiv-domain-status/)

This package is a web crawler which determines the status of a registered .hiv 
domain.

It is set up as a microservice which is managed by a RESTful API.

The server component manages the domains to check:

 - Domains to check can be added and removed
 - The status of a domain can be queried

The crawler component tries to determine the status of each domain by crawling 
the webpage and analysing the response.

 - is the domain resolving
 - can the website be accessed
 - does the returned website (after following redirects) contain the 
   click-counter snippet
 - does the redirect target (if an iframe is used) work?

## Testing

Create a databse to run the tests on:

    CREATE USER hivdomainstatus;
	CREATE DATABASE hivdomainstatus;
	GRANT ALL PRIVILEGES ON DATABASE hivdomainstatus TO hivdomainstatus;
	
	psql -H localhost -U hivdomainstatus -d hivdomainstatus < sql/domain.sql
