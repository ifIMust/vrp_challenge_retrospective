## Description
This project is the result of a time-limited coding challenge.
The program solves a specific variant of a Vehicle Routing Problem.
Each load must be picked up and dropped off at specific locations, while minimizing overall cost.

The solution employed is Tabu Search. An initial valid solution is found using a greedy algorithm. Then the neighboring solution space is explored by testing similar solutions.
The final version uses 72 iterations of Tabu Search with a Tabu list size of 20. These settings allow it to process problems with up to 200 loads in 26 seconds or less.

## Build
A Go compiler is required.
```
go build -o vrp main.go
```

### Usage
The executable reads a text file containing a VRP, and writes a solution to stdout.

Each line of the file contains Cartesian coordinates for a pickup and dropoff location.

```
./vrp file_name
```

## Design and Evolution
### Research
I searched the Internet for background information, hoping to find inspiration for a creative solution.

My best sources were the Wikipedia pages on the [VRP](https://en.wikipedia.org/wiki/Vehicle_routing_problem), [Branch and Bound](https://en.wikipedia.org/wiki/Branch_and_bound), and [Tabu Search](https://en.wikipedia.org/wiki/Tabu_search).
The core algorithm used for Tabu search is based on the description found there.

### Design
Since all location data is available from the outset, the time costs of traveling to/from the depot, completing a load, and travelling between loads, are precomputed.

#### Greedy Algorithm
Each iteration, the collection of remaining loads is sorted by proximity to the driver's current location. In this way, each driver is assigned as many loads as possible.

Although my branch and bound algorithm produced a lower cost than this greedy algorithm, there was no difference after applying Tabu Search to the solutions. The greedy approach was faster and more robust.

#### Tabu Search
The cost function heavily penalizes additional drivers, so this Tabu Search explores removing loads from routes below a certain size (tabu.maxSourceRouteSize), and placing those loads in other routes. Larger values for maxSourceRouteSize can increase the number of neighboring candidates per iteration dramatically. Smaller values are much faster, but are less likely to find optimal solutions. Testing led to setting this number to 5 to achieve quality solutions.

Performing a deep compare of a candidate to the Tabu list is expensive, so that step is delayed until necessary. A larger Tabu list size slows time performance, but too small of a list impacts solution quality.

Overall, this is a time-consuming approach, but it provides quality solutions if enough iterations and time are permitted.

#### Further Work
With more time available, further improvements would include:
- Search additional neighbors for solutions. The current approach only moves loads to other driver's routes, and does not test permutations within a route.
- Decompose nested loops in tabu.getNeighbors for improved readability.
- Use a semaphore to control the number of goroutines actively processing Tabu Search candidates. Initial tests hurt performance, but it might need tuning.
- Increase test coverage.
- Rewrite LoadDistance to use a 2D array instead of its map of maps. The hash times are slowing performance; using the maps might even be worse than just performing the computations every time.

### Development Process
Here's some more detail of how the project progressed over 48 hours.
- Addressed project I/O needs and basic structures like Locations and Loads.
- Came up with a working greedy algorithm.
- Tried to improve the cost result by implementing a Branch and Bound approach. I got one working, but it produced only slightly better outcomes than the greedy approach. It was also brittle; small changes to the bounding function caused it to degenerate into long run times.
- Considered Integer Linear Programming approach. I was unsure if I could formulate the problem correctly with limited time remaining.
- Tried a shallow Tabu search that sought to move loads from routes of length 1 to other routes. This was successful, so I expanded it to a deeper search.
- Profiled Tabu search. Comparisons to the Tabu list are expensive, so I minimized those. Used goroutines when processing the candidate queue, to relieve CPU bottleneck. This reached the desired low cost level.
