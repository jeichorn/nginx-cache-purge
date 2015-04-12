<?php
namespace NginxCP;

class Inotify 
{
	protected $proc;
    protected $path;
    protected $cmd = 'inotifywait -e moved_to -e close_write -e create -rm';
    protected $inPing = false;
    protected $lastPing = 0;
    protected $debug = false;

	public function __construct($path)
	{
        $this->path = $path;

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
                if ($file == 'ping')
                {
                    $this->inPing = false;
                    if ($this->debug)
                        echo date('Y-m-d H:i:s')." - PONG\n";
                    continue;
                }
				$updates[$path.$file] = $event;
			}
		}
		while($line !== false);

        $this->ping();
		return $updates;
	}

    public function killStaleInotify()
    {
        exec('ps aux | grep '.escapeshellarg($this->cmd).' | grep -v grep', $output);
        foreach($output as $line)
        {
            if (preg_match('/^[a-zA-Z-]+\s+([0-9]+) /', $line, $match))
            {
                $pid = $match[1];
                posix_kill($pid, SIGKILL);
            }
        }
    }

    public function ping()
    {
        $now = microtime(true);
        if ($this->inPing)
        {
            if ($now - $this->lastPing > 3)
            {
                echo date('Y-m-d H:i:s')." - Ping Failed\n";
                throw new \Exception('Ping Failed', 10);
                return false;
            }
        }
        else if ($now - $this->lastPing > 10)
        {
            if ($this->debug)
                echo date('Y-m-d H:i:s')." - PING\n";
            $this->inPing = true;
            $this->lastPing = microtime(true);
            file_put_contents($this->path."/ping", $this->lastPing);
        }
    }
}
