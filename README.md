## Description
This project is the result of a time-limited coding challenge.
The program solves a specific variant of a Vehicle Routing Problem.
Each load must be picked up and dropped off at specific locations, while minimizing overall cost.

The solution employed is Tabu Search. An initial valid solution is found using a greedy algorithm. Then the neighboring solution space is explored by testing similar solutions.

## Build
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

### Development Process
- Started with addressing project I/O needs and basic structures like Locations and Loads.
- Came up with a working greedy algorithm.
- Tried to improve the cost result by implementing a Branch and Bound approach. I was able to get one working, but it produced only slightly better outcomes than the greedy approach. It was also brittle. Small changes to the bounding function caused it to degenerate into long run times.
- Considered Integer Linear Programming approach. I was unsure if I could formulate the problem correctly with limited time remaining.
- Tried a shallow Tabu search that sought to move loads from routes of length 1 to other routes. This was successful, so I expanded it to a deep search.
- Profiled Tabu search. Comparisons to the Tabu list are expensive, so I minimized those. Used goroutines to process the candidate queue, to relieve CPU bottleneck.

### Design
Since all location data is available from the outset, the time costs of traveling to/from the depot, completing a load, and travelling between loads, are precomputed.

#### Greedy Algorithm
Each iteration, the collection of remaining loads is sorted by proximity to the driver's current location. In this way, each driver is assigned as many loads as possible.

Although my branch and bound algorithm produced a lower cost than this greedy algorithm, there was no difference after applying Tabu Search to the solutions. The greedy one is faster and more robust.

#### Tabu Search
The cost function heavily penalizes additional drivers, so this Tabu Search explores removing loads from routes below a certain size (tabu.maxSourceRouteSize), and placing those loads in other routes. Larger values for maxSourceRouteSize can increase the number of neighboring candidates per iteration dramatically. Smaller values are much faster, but are less likely to find optimal solutions. Testing led to setting this number to 5 to achieve low cost results.

Performing a deep compare of a candidate to the Tabu list is expensive, so that step is delayed until necessary. A larger Tabu list size slows time performance, but too small a list harms the cost performance.
