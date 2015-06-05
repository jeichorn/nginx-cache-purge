<?php

$tmp = "/tmp/nginxcp-test";
$files = 10000;
$fileHeader = file_get_contents(__DIR__.'/cache-file-header.txt');
$domain = "test.com";
$fakedomain = "test2.com";
$iterations = 5;
$redis = new Redis();
$redis->connect('127.0.0.1');
$loadOnly = false;
if (!empty($argv[1]) && $argv[1] == 'loadonly')
    $loadOnly = true;

if (!file_exists($tmp))
    mkdir($tmp);
passthru("rm -rf $tmp/*");
$redis->del('purge_list');


passthru("killall -9 nginx-cache-purge");
$spec = [
        0 => array("pipe", "r"),
        1 => array("file", "/dev/stdout", 'a'),
        2 => array("pipe", "w"),
];
$cwd = getcwd();
$cmd = "./bin/nginx-cache-purge --debug=1 --path=$tmp";
if (!$loadOnly)
    $process = proc_open($cmd, $spec, $pipes, $cwd, []);

$total = 0;
for($i = 0; $i < $iterations; $i++)
{
    echo "Round $i\n";

    buildCache($tmp, $files, $fileHeader, $domain);
    $ti = $t1 = microtime(true);
    for($c = 0; $c < 10; $c++)
    {
        $file = str_pad(dechex($c), 10, 'a', STR_PAD_LEFT);
        $redis->lpush("purge_list", "$domain::/$file");
    }
    for($c = 0; $c < 10; $c++)
    {
        $file = str_pad(dechex($c), 10, 'a', STR_PAD_LEFT);
        $redis->lpush("purge_list", "$fakedomain::/$file");
    }
    if ($loadOnly)
        exit;

    $count = 20;
    do
    {
        $line = fgets($pipes[2]);
        if (strstr($line, 'Tested'))
        {
            $t = round(microtime(true)-$ti,2);
            $ti = microtime(true);
            echo "$t,";
            $count--;
        }
    }
    while($count > 0);
    echo "\n";

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


