nginx cache purge [![Build Status](https://travis-ci.org/jeichorn/nginx-cache-purge.svg?branch=master)](https://travis-ci.org/jeichorn/nginx-cache-purge)
===================

Watches nginx cache dirs, keeps an index of keys, and regex purges items from a redis queue

Expects the cache-key to be in the format of: standard--httpsexample.com/

    {normalized-useragent}--{scheme}{url}{path}

All paths can have an optional ? at the end, even if there is no query string which can make building your cache key easier

Needs php 5.4+ and inotifytools installed
Expects a local redis server
