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
			$this->keys[$domain][$key] = (string)$file;
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
                $domain = 'unknown';
                $key = substr($line, 5);
                if (preg_match('@--([^/]+)/@', $key, $match))
                    $domain = $match[1];

                return array($domain, $key);
            }
		}
		echo date('Y-m-d H:i:s')." - did't find a key in $file\n";
	}

	public function update($updates)
	{
		$updated = false;
		foreach($updates as $file => $status)
		{
            if (is_file($file))
            {
                list($domain, $key) = $this->keyFromFile($file);
                $this->keys[$domain][$key] = (string)$file;
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
		$regex = "|^([a-zA-Z0-9]+--)?(https?)?$host$regex$|";

		echo date('Y-m-d H:i:s')." - checking $rule with $regex\n";
        $count = 0;
        $unlink = 0;
        $s = microtime(true);
        if (isset($this->keys[$host]))
        {
            foreach($this->keys[$host] as $key => $file)
            {
                if (preg_match($regex, $key))
                {
                    echo date('Y-m-d H:i:s')." - Found a match $key\n";
                    $t = microtime(true);
                    @unlink($file);
                    unset($this->keys[$host][$key]);
                    $unlink += (microtime(true)-$t);
                    $count++;
                }
                else
                {
                    //echo date('Y-m-d H:i:s')." - Miss on $key\n";
                }
            }
        }
        else
        {
            echo date('Y-m-d H:i:s')." - No keys for $host\n";
        }
        $total = microtime(true)-$s;

        echo date('Y-m-d H:i:s')." - $rule took $total unlink took $unlink\n";

        return $count;
	}
}
