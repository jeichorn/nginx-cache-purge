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
		$iterator = new RecursiveIteratorIterator(new RecursiveDirectoryIterator($path, FilesystemIterator::SKIP_DOTS));
		foreach($iterator as $file)
		{
			$this->keys[$this->keyFromFile($file)] = (string)$file;
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
				return substr($line, 5);
		}
		echo date('Y-m-d H:i:s')." - did't find a key in $file\n";
	}

	public function update($updates)
	{
		$updated = false;
		foreach($updates as $file => $status)
		{
			if ($status == 'DELETE')
			{
				$success = false;
				foreach($this->keys as $key => $f)
				{
					if ($file == $f)
					{
						unset($this->keys[$key]);
						$success = true;
						break;
					}
				}
				if (!$success)
				{
					echo date('Y-m-d H:i:s')." didn't find $file in the cache to delete\n";
				}
				else
				{
					$updated = true;
				}
			}
			else
			{
				if (is_file($file))
				{
					$this->keys[$this->keyFromFile($file)] = $file;
					$updated = true;
				}
			}
		}

		if ($updated)
		{
			$newcount = count($this->keys);
			echo date('Y-m-d H:i:s')." - update keys now have $newcount\n";
		}
	}

	public function purge($rule)
	{
		list($host, $path) = explode('::', $rule);

		$regex = preg_quote(str_replace('(.*)', '@@@', $path), '|');
		$regex = str_replace('@@@', '(.*)', $regex);	
		$regex = "|^https?$host$regex|";

		echo date('Y-m-d H:i:s')." - checking $rule with $regex\n";
		foreach($this->keys as $key => $file)
		{
			if (preg_match($regex, $key))
			{
				echo date('Y-m-d H:i:s')." - Found a match $key\n";
				unlink($file);
				//unset($this->keys[$key]); inotify will tell us to remove the key
			}
			else
			{
				echo date('Y-m-d H:i:s')." - Miss on $key\n";
			}
		}
	}
}
