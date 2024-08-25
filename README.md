## Description
The [vrp_challenge](https://github.com/ifIMust/vrp_challenge) project was limited to 48 hours of time, so it cannot receive further updates.
This repository is based on that project, and serves to test additional improvements or variations without impacting the "hackathon" version.
See [vrp_challenge](https://github.com/ifIMust/vrp_challenge) for build instructions and more information.

The original solution produced these results (using given test data):
```
mean cost: 45744.88601842658
mean run time: 9035.528087615967ms
```

This version, using the same data:
```
mean cost: 45302.30033360918
mean run time: 18394.224894046783ms
```

The following improvements were made:
- **Enforce time limit.** 30s to complete a problem with 200 loads was a *requirement*. A timer is used to stop the Tabu search when time is running low.
- **Explore a larger neighbor space**
  - Move loads within the same route.
  - Remove limit on size of source routes for load reassignment.
- **Improve parallelism:** Use an additional goroutine to read and process results alongside queue workers.

By relying on the time limit to control iterations, the number of Tabu Search iterations was increased.
Using longer run times to optimize cost produces a very high quality solution.
The neighbor space is explored much more deeply, especially by removing the size limit of routes that might have loads moved from them.
