<?php
namespace NginxCP;

class Inotify 
{
	protected $proc;
    protected $cmd = 'inotifywait -e delete -e moved_to -e close_write -e create -e delete -rm';

	public function __construct($path)
	{
        // we can leak these guys so cleanup on start
        $this->killStaleInotify();

		// start a inotifywait process
		$this->proc = popen("$this->cmd $path", "r");
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

    public function killStaleInotify()
    {
        exec('ps aux | grep '.escapeshellarg($this->cmd).' | grep -v grep', $output);
        foreach($output as $line)
        {
            if (preg_match('/^[a-zA-Z-]+ ([0-9]+) /', $line, $match))
            {
                $pid = $match[1];
                posix_kill($pid, SIGKILL);
            }
        }
    }
}
