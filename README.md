## Description
The [vrp_challenge](https://github.com/ifIMust/vrp_challenge) project was limited to 48 hours of time, so it cannot receive further updates.
This repository is based on that project, serves to test additional improvements or variations without impacting the "hackathon" version.

With the benefit of sleep, here are some retrospective thoughts about that project:
- 30s to complete a problem with 200 loads was a *requirement*.
  - Not enforcing the time limit in the program is a dangerous omission. A timer should be used to stop the Tabu search when time is running low.
  - Furthermore, the limit can be leveraged to search for a better solution for as long as possible, possibly achieving an even better solution.
- The exisiting search needs to explore a larger neighbor space to find better solutions more directly.
  - The most obvious omission to add is to permute each route in place.
- Creating hundreds of goroutines in a loop is not ideal. Control goroutine creation rate for possible performance gains.
- The algorithm has several parameters that tune the behavior. These should be in a config file to avoid recompiling when tuning.

See [vrp_challenge](https://github.com/ifIMust/vrp_challenge) for build instructions and more information.
