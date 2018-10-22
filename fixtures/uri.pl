#! /usr/bin/perl
use strict;
use Data::Validate::URI qw(is_uri);
use Data::Dumper;

#if (is_uri('http://c:/directory/file')) {
if (is_uri('http://httpbin.org/get?utf8=\xe2\x98\x83')) {
    print "ok\n";
} else {
    print "error\n";
};
