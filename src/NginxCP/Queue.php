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
		return $this->redis->lpop('purge_list');
	}
	
}
