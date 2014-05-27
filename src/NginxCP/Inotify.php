<?php
namespace NginxCP;

class Inotify 
{
	protected $proc;

	public function __construct($path)
	{
		// start a inotifywait process
		$this->proc = popen("inotifywait -e delete -e moved_to -e close_write -e create -e delete -rm $path", "r");
        if (!$this->proc)
        {
            throw new \Exception("Runing inotifywait failed");
        }
		stream_set_blocking($this->proc, 0);
	}

	public function getUpdates()
	{
		$updates = [];
		do
		{
			$line = fgets($this->proc);
			if ($line !== false)
			{
				list($path, $event, $file) = explode(" ", rtrim($line));
				$updates[$path.$file] = $event;
			}
		}
		while($line !== false);

		return $updates;
	}
}
