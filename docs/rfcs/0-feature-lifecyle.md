- Author: xxchan <xxchan22f@gmail.com>
- Start Date: 2021-07-15
- RFC PR: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)
- Tracking Issue: [beyondstorage/go-storage#0](https://github.com/beyondstorage/go-storage/issues/0)

# GSP-0: Feature Lifecycle

## Background

Previous discussions: 
- https://forum.beyondstorage.io/t/how-to-maintain-go-service-more-easily/137
- https://xuanwo.io/reports/2021-28/

When we implemented new features, the process is:
1. [`go-storage`] An RFC is approved according to the RFC process.
2. [`go-storage`] A tracking issue is created.
3. [`go-storage`] The feature is implemented.
4. [`go-service-*`] The feature is implemented. Cases are:
   - Update `go-storage` and generate new code.
   - Add some boilerplate code.
   - Need special treatment for each service.
5. [`go-storage`] The tracking issue is closed.

After a few cycles of 1-5，
6. [`go-storage`] Release a new version.
7. [`go-service-*`] Release a new version.

Problems are:
- **Repetitive work**: The first case of step 4 is highly repetitive, and can be easily automated with dependabot and GitHub actions. And we don't need to update for each feature. We can update only once to include several features when `go-storage` releases a new version.

- **Ineffective tracking**: Step 5 is blocked by step 4. 
  
  There were statements like "make sure all services are updated". Such obscure tasks should be avoided.

  But it is also improper to specify explicitly the services to be updated. First, We separated `go-service-*` for low coupling, but if every feature should be immediately synced for all the services, it is a very burden for maintainers. `go-service-*` cannot implement several features together. Second, too many tracking issues won't be closed, which will confuse maintainers, and `go-storage` cannot release new versions quickly.

  This does not mean that we shouldn't track the implementation status of `go-service-*`. Instead, we can create a tracking issue in the service repo.

Before, we only specified the RFC process, but other parts of feature lifecycle (or release cycle) are not specified and thus the current practice leads to problems above.

## Proposal

So I propose to specify the feature lifecycle: 

1. [`go-storage`] An RFC SHOULD be approved according to the RFC process.
2. [`go-storage`] A tracking issue SHOULD be created, including what `go-service-*` should do.
3. [`go-storage`] The feature is implemented.
4. [`go-storage`] The tracking issue is closed. At the same time, create tracking issues for all the services (can be automated using `go-community`). 

Now the feature is viewed as implemented, and `go-storage` CAN release a new version.

Regarding as the service tracking issues:
- If it is only required to update `go-storage`, the task CAN be delayed when `go-storage` releases a new version and automatically done by dependabot.
- If manual work is required, implement it correspondingly. This CAN be done immediately or when `go-storage` releases a new version.
  
  It is possible that some other tasks (only need to update `go-storage`) are done by the way, and their tracking issues can be closed.

## Rationale

Improve maintainability.

## Compatibility

N/A

## Implementation

- Implement an automated tool to create tracking issues for all the services.
- Implement CI to let dependabot run `go-storage` code generator automatically.