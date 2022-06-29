# Car Pooling Service Challenge

Design/implement a system to manage car pooling.

At Cabify we provide the service of taking people from point A to point B.
So far we have done it without sharing cars with multiple groups of people.
This is an opportunity to optimize the use of resources by introducing car
pooling.

You have been assigned to build the car availability service that will be used
to track the available seats in cars.

Cars have a different amount of seats available, they can accommodate groups of
up to 4, 5 or 6 people.

People requests cars in groups of 1 to 6. People in the same group want to ride
on the same car. You can take any group at any car that has enough empty seats
for them. If it's not possible to accommodate them, they're willing to wait until 
there's a car available for them. Once a car is available for a group
that is waiting, they should ride. 

Once they get a car assigned, they will journey until the drop off, you cannot
ask them to take another car (i.e. you cannot swap them to another car to
make space for another group).

In terms of fairness of trip order: groups should be served as fast as possible,
but the arrival order should be kept when possible.
If group B arrives later than group A, it can only be served before group A
if no car can serve group A.

For example: a group of 6 is waiting for a car and there are 4 empty seats at
a car for 6; if a group of 2 requests a car you may take them in the car.
This may mean that the group of 6 waits a long time,
possibly until they become frustrated and leave.

## Evaluation rules

This challenge has a partially automated scoring system. This means that before
it is seen by the evaluators, it needs to pass a series of automated checks
and scoring.

### Checks

All checks need to pass in order for the challenge to be reviewed.

- The `acceptance` test step in the `.gitlab-ci.yml` must pass in master before you
submit your solution. We will not accept any solutions that do not pass or omit
this step. This is a public check that can be used to assert that other tests 
will run successfully on your solution. **This step needs to run without 
modification**
- _"further tests"_ will be used to prove that the solution works correctly. 
These are not visible to you as a candidate and will be run once you submit 
the solution

### Scoring

There is a number of scoring systems being run on your solution after it is 
submitted. It is ok if these do not pass, but they add information for the
reviewers.

## API

To simplify the challenge and remove language restrictions, this service must
provide a REST API which will be used to interact with it.

This API must comply with the following contract:

### GET /status

Indicate the service has started up correctly and is ready to accept requests.

Responses:

* **200 OK** When the service is ready to receive requests.

### PUT /cars

Load the list of available cars in the service and remove all previous data
(existing journeys and cars). This method may be called more than once during 
the life cycle of the service.

**Body** _required_ The list of cars to load.

**Content Type** `application/json`

Sample:

```json
[
  {
    "id": 1,
    "seats": 4
  },
  {
    "id": 2,
    "seats": 6
  }
]
```

Responses:

* **200 OK** When the list is registered correctly.
* **400 Bad Request** When there is a failure in the request format, expected
  headers, or the payload can't be unmarshalled.

### POST /journey

A group of people requests to perform a journey.

**Body** _required_ The group of people that wants to perform the journey

**Content Type** `application/json`

Sample:

```json
{
  "id": 1,
  "people": 4
}
```

Responses:

* **200 OK** or **202 Accepted** When the group is registered correctly
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

### POST /dropoff

A group of people requests to be dropped off. Whether they traveled or not.

**Body** _required_ A form with the group ID, such that `ID=X`

**Content Type** `application/x-www-form-urlencoded`

Responses:

* **200 OK** or **204 No Content** When the group is unregistered correctly.
* **404 Not Found** When the group is not to be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

### POST /locate

Given a group ID such that `ID=X`, return the car the group is traveling
with, or no car if they are still waiting to be served.

**Body** _required_ A url encoded form with the group ID such that `ID=X`

**Content Type** `application/x-www-form-urlencoded`

**Accept** `application/json`

Responses:

* **200 OK** With the car as the payload when the group is assigned to a car. See below for the expected car representation 
```json
  {
    "id": 1,
    "seats": 4
  }
```

* **204 No Content** When the group is waiting to be assigned to a car.
* **404 Not Found** When the group is not to be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshalled.

## Tooling

At Cabify, we use Gitlab and Gitlab CI for our backend development work. 
In this repo you may find a [.gitlab-ci.yml](./.gitlab-ci.yml) file which
contains some tooling that would simplify the setup and testing of the
deliverable. This testing can be enabled by simply uncommenting the final
acceptance stage. Note that the image build should be reproducible within
the CI environment.

Additionally, you will find a basic Dockerfile which you could use a
baseline, be sure to modify it as much as needed, but keep the exposed port
as is to simplify the testing.

:warning: Avoid dependencies and tools that would require changes to the 
`acceptance` step of [.gitlab-ci.yml](./.gitlab-ci.yml), such as 
`docker-compose`

:warning: The challenge needs to be self-contained so we can evaluate it. 
If the language you are using has limitations that block you from solving this 
challenge without using a database, please document your reasoning in the 
readme and use an embedded one such as sqlite.

You are free to use whatever programming language you deem is best to solve the
problem but please bear in mind we want to see your best!

You can ignore the Gitlab warning "Cabify Challenge has exceeded its pipeline 
minutes quota," it will not affect your test or the ability to run pipelines on
Gitlab.

## Requirements

- The service should be as efficient as possible.
  It should be able to work reasonably well with at least $`10^4`$ / $`10^5`$ cars / waiting groups.
  Explain how you did achieve this requirement.
- You are free to modify the repository as much as necessary to include or remove
  dependencies, subject to tooling limitations above.
- Document your decisions using MRs or in this very README adding sections to it,
  the same way you would be generating documentation for any other deliverable.
  We want to see how you operate in a quasi real work environment.

## Feedback

In Cabify, we really appreciate your interest and your time. We are highly 
interested on improving our Challenge and the way we evaluate our candidates. 
Hence, we would like to beg five more minutes of your time to fill the 
following survey:

- https://forms.gle/EzPeURspTCLG1q9T7

Your participation is really important. Thanks for your contribution!

# Design and implementation decisions
#### Programming Language
* For the implementation, I have decided to use Go. The reasons are that 
I have interest to learn it and that it is a very popular language.

* This is the first time using go, so it is possible that some idiomatic style 
is missing sometimes and the code could be optimized easily.

#### Testing
* Several unit tests have been created after the implementation, but to keep 
them agnostic of the implementation, they will work simulating calls to the API. 
I will only check if the bodies and responses of the calls matched the expected 
ones.

#### Data Storage
* I have decided not to use a database, so I will have to use several variables 
to store the relation between 

* The service will be use several HashTables (in go this means the Data 
Structure Map<key,value>) to represent the relationship between groups of People, 
cars, free seats in cars and the cars assigned to each group.

* A dinamic array/vector/slice will be used to store the waiting list.

* When restarting the info of the cars and journey because of the `PUT /cars` 
call, we delete the keys of all the maps/hashtables. Here I assume go Maps 
manage the amount of memory they need when we add or delete values. For the 
waiting list, we free the memory and start a new one. This way, the 
allocated memory reflects the amount of data we are currently using.

* Since the values of a map can be the key of another map, we can say that 
we have redundant data in memory and that this solution takes more memory.
In exchange it should help to speed up the processing of the requests for 
`POST /journey`, `POST /dropoff` and `POST /locate`.

#### Input related decisions
* The format of the requests bodies must match the few samples inputs 
provided. This means:
    1. For `PUT /cars` each object has `id` and `seats` are required. If 
    one of them is missing in an object, it will be considered a **Bad 
    request**.
    2. For `PUT /cars`, a repeated Id won't be allowed. If an Id 
    is repeated, it will return a **400 Bad Request** response.
    3. For `POST /journey` the fields `id` and `people` are required. If 
    one of them is missing, it will be considered a **Bad request**.
    4. For `POST /dropoff` and `POST /locate` the only Pattern accepted 
    is `ID=X`, being X a single positive number. Otherwise, it is a 
    **Bad Request**.
    5. For cars, the `ID 0` can not be used in the input. Internally
    if a group is assigned to the car with ID 0, it means that the 
    group is waiting to be assigned to a car. For groups the id 0 is 
    allowed.
    6. Negative values are not allowed it doesn't make sense to use 
    them for the amount of seats or the people in a group. And in 
    samples IDs start at 1, we also have the 0 reserved for cars.
    So I consider a good decision to cap the negative values in the 
    inputs. A negative value will return a bad request.


#### Other decisions
* Some diagrams may be attached since they helped to think about how to 
implement some of the more complex operations.


### Performance thoughts(based on implementation details)
* For the management of the data, I'm using several Map Object to represent the
relations.
* Usually in all languages the Map equivalents promises in average a constant 
O(1), or O(n) worst case for insertion, search and delete.
* Since we can implement each of the requests as a bunch of access, insert and 
delete in the different maps the performance should be good without requiring 
to do black magic.
* The most time consuming operation, and bottleneck becomes the management of 
the waiting list, since there is no easy way to search and delete in an ordered 
collection in constant time.

* Since here is mentioned that the service should work with 100k cars and groups
a performance test has been created, if you execute stressTest.go it will create 
a request that will add $$10^5$$ cars throught the `PUT /cars` method. 
After that it will execute:
    1. $15 \times `10^4`$ `POST /journey` requests 
    2. $15 \times `10^4`$ `POST /locate` requests(for random groups ID)
    3. $15 \times `10^4`$ `POST /dropoff` requests(it should delete all groups)
    added for the test
* Next you can see the results of one of the executions of the test:
![image](./test_times.jpg)

#### Used tools
* VS Code was used to implement the API, witht go extensions

* For API testing during implementation, I used postman to create the requests, 
and picked some code from Postman tool to prepare the calls in the unit tests.

#### Future Work

* For future work, I could change the implementation to achieve a domain based
architecture, using a layer for data Access, a layer for business logic and a 
layer for service logic.

### GET /status
* When server is ready, a GET method will return a **200 OK** response 
    without any message. 
* A call with another method will return an **405 Method 
    Not Allowed** response. 
* Since there is no logic behind the call, only the 
    HTTP method is checked.

### PUT /cars
* Only the PUT method is allowed. Another method will return an **405 
    Method Not Allowed** response.
* The PUT method will return a **400 Bad Request** if one of the following 
conditions is met: 
    1. Request Body is empty 
    2. **Content Type** is not `application/json`
    3. The json is bad formatted
    4. The json does not match the expected pattern (e.g the id/seats 
    are missing, id = 0 or repeated, seats below minimun, etc)
* Otherwise it will return a **200 OK** response.

### POST /journey
* Only the POST method is allowed. Another method will return an **405 
    Method Not Allowed** response.
* The POST method will return a **400 Bad Request** if one of the following 
conditions is met: 
    1. Request Body is empty 
    2. **Content Type** is not `application/json`
    3. The json is bad formatted
    4. The json does not match the expected pattern (e.g the id/people 
    are missing,  people is a negative number, people below minimun, etc)
* The POST method will return a **500 Internal Server Error** if there is
already a group with the same id in the service
* Otherwise it will return a **200 OK** response.

### POST /dropoff
* Only the POST method is allowed. Another method will return an **405 
    Method Not Allowed** response.
* The POST method will return a **400 Bad Request** if one of the following 
conditions is met: 
    1. Request Body is empty 
    2. **Content Type** is not `application/x-www-form-urlencoded`
    3. There is more than 1 different key.
    4. The unique key is not ID
    5. There are several values for the ID key
    6. The value for ID is not an 
* The POST method will return a **404 Not Found** if there is no group with the 
request's ID
* The POST method will return a **204 No Content** if the group has no car 
assigned. If the dropped group is in a car, it will return a **200 OK** 
response.

### POST /locate

* Only the POST method is allowed. Another method will return an **405 
    Method Not Allowed** response.
* The **Content Type** of the response will be always `application/json`
* The POST method will return a **400 Bad Request** if one of the following 
conditions is met: 
    1. Request Body is empty 
    2. **Content Type** is not `application/x-www-form-urlencoded`
    3. There is more than 1 different key.
    4. The unique key is not ID
    5. There are several values for the ID key
    6. The value for ID is not an 
* The POST method will return a **404 Not Found** if there is no group with the 
request's ID
* The POST method will return a **204 No Content** if the dropped group had 
no car assigned. If the dropped group was assigned to a car, it will return 
a **200 OK** response. The body will be a json with the car data, matching the 
provided sample's pattern.
