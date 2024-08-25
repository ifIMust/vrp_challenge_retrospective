## Description
The [vrp_challenge](https://github.com/ifIMust/vrp_challenge) project was limited to 48 hours of time, so it cannot receive further updates.
This repository is based on that project, and serves to test additional improvements or variations without impacting the "hackathon" version.

With the benefit of sleep, here are some retrospective thoughts about that project:
- 30s to complete a problem with 200 loads was a *requirement*. A timer should be used to stop the Tabu search when time is running low. (Done)
- The exisiting search should explore a larger neighbor space to find better solutions.
  - Allow moving loads within the same route. (Done)

See [vrp_challenge](https://github.com/ifIMust/vrp_challenge) for build instructions and more information.
