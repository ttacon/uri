#! /usr/bin/python
from rfc3986 import urlparse
#uri = urlparse('oasis:names:specification:docbook:dtd:xml:4.1.2')
#uri = urlparse('http://httpbin.org/get?utf8=\xe2\x98\x83')
#uri = urlparse('mailto:user@domain.com')
uri = urlparse('mailto://user@domain.com/abc#a?b')
print uri
