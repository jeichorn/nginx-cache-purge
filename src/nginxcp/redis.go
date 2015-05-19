package nginxcp

import(
    "gopkg.in/redis.v3"
    "time"
    "fmt"
)

type RedisQueue struct {
    client *redis.Client
    Active bool
    Jobs chan string
}

func NewRedisQueue() *RedisQueue {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    // load initial needs to set active true, instead of setting it here
    return &RedisQueue{client, true, make(chan string, 10)}
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
        if (queue.Active) {
            job := queue.getJob()

            if (job == "") {
                time.Sleep(1 * time.Second)
            } else {
                DebugMessage(fmt.Sprintf("Adding a job to the channel: %#v\n", job))
                queue.Jobs <- job
            }
        } else {
            time.Sleep(1 * time.Second)
        }
    }
}
