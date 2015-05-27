package nginxcp

import(
    "gopkg.in/redis.v3"
    "time"
)

type RedisQueue struct {
    client *redis.Client
    Jobs chan string
}

func NewRedisQueue() *RedisQueue {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    return &RedisQueue{client, make(chan string, 10)}
}

func (queue *RedisQueue) getJob() string {
    job := queue.client.RPop("purge_list")

    return job.Val()
}

func (queue *RedisQueue) completeJob(job string) {
    queue.client.HDel("in_purge_list", job)
}

func (queue *RedisQueue) Run() {
    for {
        job := queue.getJob()

        if (job == "") {
            time.Sleep(1 * time.Second)
        } else {
            PrintTrace2("Adding a job to the channel: %#v\n", job)
            queue.Jobs <- job
        }
    }
}
