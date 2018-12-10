#! /usr/bin/python
from rfc3986 import urlparse, uri_reference
#uri = urlparse('oasis:names:specification:docbook:dtd:xml:4.1.2')
#uri = urlparse('http://httpbin.org/get?utf8=\xe2\x98\x83')
#uri = urlparse('mailto:user@domain.com')
uri = urlparse('mailto://user@domain.com/abc#a?b')
print uri
uri = uri_reference('//abc#a?b')
print uri
uri = uri_reference('')
print uri
uri = uri_reference('../dir')
print uri
uri = uri_reference('file.html')
print uri
