<?php
namespace NginxCP;

class Queue
{
	protected $redis;

	public function __construct($redis)
	{
		$this->redis = $redis;
	}

	public function getJob()
	{
		return $this->redis->rpop('purge_list');
	}

    public function completeJob($job)
    {
        $this->redis->hdel("in_purge_list", $job);
    }
}
