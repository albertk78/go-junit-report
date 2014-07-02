go-junit-report
===============

Converts `go test` output to an xml report, suitable for applications that
expect junit xml reports (e.g. [Jenkins](http://jenkins-ci.org)).

This is a modificated version of [jstemmer's work](https://github.com/jstemmer/go-junit-report), to support nested test cases. You may need this version if your write your tests by [testify `suite` package](https://github.com/stretchr/testify#suite-package).

Installation
------------

	go get github.com/wancw/go-junit-report

Usage
-----

	go test -v | go-junit-report > report.xml

