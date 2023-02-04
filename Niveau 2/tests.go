

jobs := make(chan int, numJobs)
// w est le nombre de jobs à faire, avec la cmde de PFR
for w := 1; w <= 3; w++ {
	go box_blur(w, jobs, results)
}

//ça c'est le nombre total de jobs qui devront être faits au final
for j := 1; j <= numJobs; j++ {
	jobs <- j
}
close(jobs)

func box_blur(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("worker", id, "started  job", j)
        time.Sleep(time.Second)
        fmt.Println("worker", id, "finished job", j)
        //results <- j * 2
    }
}