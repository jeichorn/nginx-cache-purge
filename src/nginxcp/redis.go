package nginxcp

import(
    "gopkg.in/redis.v3"
    "time"
)

type JobBag struct {
    Bag []string
    Count int
}
type RedisQueue struct {
    client *redis.Client
    Jobs chan JobBag
}

func NewRedisQueue() *RedisQueue {
    client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

    return &RedisQueue{client, make(chan JobBag, 1)}
}

func (queue *RedisQueue) getJob() string {
    job := queue.client.RPop("purge_list")

    return job.Val()
}

func (queue *RedisQueue) completeJob(job string) {
    queue.client.HDel("in_purge_list", job)
}

func (queue *RedisQueue) clearInPurgeList() {
    queue.client.Del("in_purge_list")
}

func (queue *RedisQueue) Run() {
    bag := JobBag{make([]string, 0), 0}
    for {
        job := queue.getJob()
        if (job != "") {
            queue.completeJob(job)
            bag.Bag = append(bag.Bag, job)
            bag.Count++
        }

        if (job == "" || bag.Count > 10) {
            if (bag.Count > 0) {
                PrintTrace1("Adding %d jobs to the channel: %#v\n", bag.Count, bag)
                queue.Jobs <- bag
                bag.Bag = make([]string, 0)
                bag.Count = 0
            } else {
                time.Sleep(1 * time.Second)
            }
        }
    }
}
