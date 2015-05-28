<?php

$tmp = "/tmp/nginxcp-test";
$files = 1000;
$fileHeader = file_get_contents(__DIR__.'/cache-file-header.txt');
$domain = "test.com";
$iterations = 5;
$redis = new Redis();
$redis->connect('127.0.0.1');

if (!file_exists($tmp))
    mkdir($tmp);
passthru("rm -rf $tmp/*");
$redis->del('purge_list');


$spec = [
        0 => array("pipe", "r"),
        1 => array("file", "/dev/stdout", 'a'),
        2 => array("file", "/dev/null", 'a'),
];
$cwd = getcwd();
$cmd = "./bin/nginx-cache-purge --debug=1 --path=$tmp";
$process = proc_open($cmd, $spec, $pipes, $cwd, []);

$total = 0;
for($i = 0; $i < $iterations; $i++)
{
    echo "Round $i\n";
    buildCache($tmp, $files, $fileHeader, $domain);
    $t1 = microtime(true);
    for($c = 0; $c < ($files/10); $c++)
    {
        $file = str_pad(dechex($c), 10, 'a', STR_PAD_LEFT);
        $redis->lpush("purge_list", "$domain::/$file");
    }

    do
    {
        $count = $redis->llen('purge_list');
        if ($count > 0)
            usleep(50000);
    }
    while($count > 0);

    $round = microtime(true)-$t1;
    echo "Round Complete in $round\n";
    $total += $round;
}
echo "$i iterations complete in $total\n";


function buildCache($tmp, $files, $fileHeader, $domain)
{
    for($i = 0; $i < $files; $i++)
    {
        $file = str_pad(dechex($i), 10, 'a', STR_PAD_LEFT);
        $dir = substr($file, -2);
        if (!file_exists("$tmp/$dir"))
            mkdir("$tmp/$dir");

        file_put_contents("$tmp/$dir/$file", $fileHeader."KEY: standard--https$domain/$file\n");
    }
}


