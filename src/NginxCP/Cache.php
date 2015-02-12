<?php
namespace NginxCP;
use RecursiveIteratorIterator;
use RecursiveDirectoryIterator;
use FilesystemIterator;


class Cache
{
	public static $CACHE_PATH = '/mnt/cache';
	public $keys = [];

	public function scan($path)
	{
        $this->keys = [];
		$iterator = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($path, FilesystemIterator::SKIP_DOTS));
		foreach($iterator as $file)
		{
            list($domain, $key) = $this->keyFromFile($file);

            $sfile = (string)$file;

			$this->keys[$domain][$key][$sfile] = $sfile;
		}
        echo date('Y-m-d H:i:s')." Key counts by domain: \n";
		foreach($this->keys as $domain => $keys)
		{
			echo date('Y-m-d H:i:s')." - $domain - ".count($keys)."\n";
		}
	}

	public function keyFromFile($file)
	{
		$fp = fopen($file, 'r');
		$top = fread($fp, 512);
		fclose($fp);
		$lines = explode("\n", $top);
		foreach($lines as $line)
		{
			if (substr($line, 0, 4) === "KEY:")
            {
                return $this->parseKey($line);
            }
		}
		echo date('Y-m-d H:i:s')." - did't find a key in $file\n";
	}

    public function parseKey($line)
    {
        $domain = 'unknown';
        $key = substr($line, 5); // skip over 'KEY: '
        if (preg_match('@--([^/]+)/@', $key, $match))
            $domain = $match[1];

        return array($domain, $key);
    }

	public function update($updates)
	{
		$updated = false;
		foreach($updates as $file => $status)
		{
            if (is_file($file))
            {
                list($domain, $key) = $this->keyFromFile($file);
                $sfile = (string)$file;

                $this->keys[$domain][$key][$sfile] = $sfile;
            }
		}

		if ($updated)
		{
			$newcount = count($this->keys);
			echo date('Y-m-d H:i:s')." - update keys now have $newcount domains\n";
		}
	}

	public function purge($rule)
	{
		list($host, $path) = explode('::', $rule);

		$regex = preg_quote(str_replace('(.*)', '@@@', $path), '|');
		$regex = str_replace('@@@', '(.*)', $regex);	

        // this assumes you have cache keys like
        // normalizedua--httpHostnamePath
        // https urls also match and normalizedua-- is optional
		$regex = "|^([a-zA-Z0-9]+--)?(https?)?$host$regex\??$|";

        $count = 0;
        $unlink = 0;
        $s = microtime(true);
        $possible = array($host, 'http'.$host, 'https'.$host);
        $found = false;
        foreach($possible as $index)
        {
            if (isset($this->keys[$index]))
            {
                $found = true;
                $count = count($this->keys[$index]);
                echo date('Y-m-d H:i:s')." - $host has $count keys checking $rule with $regex\n";
                foreach($this->keys[$index] as $key => $files)
                {
                    if (preg_match($regex, $key))
                    {
                        echo date('Y-m-d H:i:s')." - Found a match $key\n";
                        foreach($files as $file)
                        {
                            @unlink($file);
                        }
                        unset($this->keys[$index][$key]);
                        $unlink++;
                    }
                    else
                    {
                        //echo date('Y-m-d H:i:s')." - Miss on $key\n";
                    }
                }
                $total = round(microtime(true)-$s,4);

                echo date('Y-m-d H:i:s')." - $unlink key(s) killed in $total $rule\n";
            }
        }

        if (!$found)
        {
            echo date('Y-m-d H:i:s')." - No keys for $host\n";
        }

        return $unlink;
	}
}
