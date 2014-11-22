<?php

require __DIR__.'/../vendor/autoload.php';

use NginxCP\Cache;

Cache::$CACHE_PATH = __DIR__;
$cache = new Cache();

$keys = [
    "standard--httpsexample.com/wp-content/themes/twentytwelve/style.css??3.9.1",
    "standard--example.com/wp-content/themes/twentytwelve/style.css??3.9.1",
    "standard--httpexample.com/wp-content/themes/twentytwelve/style.css??3.9.1"
];

foreach ($keys as $k)
{
    list($dom, $key) = $cache->parseKey($k);
    assert('$dom == "example.com" && $key = $k');
}
