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
for them. If it's not possible to accommodate them, they're willing to wait.

Once they get a car assigned, they will journey until the drop off, you cannot
ask them to take another car (i.e. you cannot swap them to another car to
make space for another group). In terms of fairness of trip order: groups are
served in the order they arrive, but they ride opportunistically.

For example: a group of 6 is waiting for a car and there are 4 empty seats at
a car for 6; if a group of 2 requests a car you may take them in the car for
6 but only if you have nowhere else to make them ride. This may mean that the
group of 6 waits a long time, possibly until they become frustrated and
leave.

## Acceptance

The acceptance test step in the `.gitlab-ci.yml` must pass before you submit
your solution. We will not accept any solutions that do not pass or omit this
step.

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
    "seats": 7
  }
]
```

Responses:

* **200 OK** When the list is registered correctly.
* **400 Bad Request** When there is a failure in the request format, expected
  headers, or the payload can't be unmarshaled.

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
  payload can't be unmarshaled.

### POST /dropoff

A group of people requests to be dropped off. Wether they traveled or not.

**Body** _required_ A form with the group ID, such that `ID=X`

**Content Type** `application/x-www-form-urlencoded`

Responses:

* **200 OK** or **204 No Content** When the group is unregistered correctly.
* **404 Not Found** When the group is not to be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshaled.

### POST /locate

Given a group ID such that `ID=X`, return the car the group is traveling
with, or no car if they are still waiting to be served.

**Body** _required_ A url encoded form with the group ID such that `ID=X`

**Content Type** `application/x-www-form-urlencoded`

**Accept** `application/json`

Responses:

* **200 OK** With the car as the payload when the group is assigned to a car.
* **204 No Content** When the group is waiting to be assigned to a car.
* **404 Not Found** When the group is not to be found.
* **400 Bad Request** When there is a failure in the request format or the
  payload can't be unmarshaled.

## Tooling

In this repo you may find a [.gitlab-ci.yml](./.gitlab-ci.yml) file which contains
contains some tooling that would simplify the setup and testing of the
deliverable. This testing can be enabled by simply uncommenting the final
acceptance stage.

Additionally, you will find a basic Dockerfile which you could use a
baseline, be sure to modify it as much as needed, but keep the exposed port
as is to simplify the testing.

You are free to modify the repository as much as necessary to include or remove
dependencies, but please document these decisions using MRs or in this very
README adding sections to it, the same way you would be generating
documentation for any other deliverable. We want to see how you operate in a
quasi real work environment.

## Third party libraries

### Gorilla/Mux
Package _gorilla/mux_ implements a request router and dispatcher for matching
incoming requests to their respective handler.

The name mux stands for "HTTP request multiplexer". Like the standard http.ServeMux,
mux.Router matches incoming requests against a list of registered routes
and calls a handler for the route that matches the URL or other conditions.

We use _gorilla/mux_ because that's what I saw being used in several places. It's very capable and
makes it easier to specify which methods should be allowed for specific routes. A good thing about
_gorilla/mux_ is that it's compatible with the built-in http.Handler API.

You can find more detailed information about this library on its GitHub website:

[https://github.com/gorilla/mux]

### Go Playground Validator

Go Playground Validator implements value validations for structs and individual fields based on tags.

Input validation prevents improperly formed data from entering an information system.
Because it is difficult to detect a malicious user who is trying to attack software,
applications should check and validate all input entered into a system.

We use Go Playground Validator because it is, as far as we know, the most famous
and has the most stars on GitHub among the existing libraries.

You can find more detailed information about this library on its GitHub website:

[https://github.com/go-playground/validator]

## To-do list

There are several features/improvements we do not have currently implemented,
but would be a nice-to-have:

* **Database** Currently we use slices to store the information about the cars and the journeys.
  In order to have some persistance in the data (not limited only to the program execution life)
  we would need to use a database to recover the status in case of crashing.
* **HTTPS** We only provide our REST API via HTTP. This implies the traffic is sent in clear,
  which makes it susceptible to eavesdropping, man-in-the-middle attacks, and so son.
  We should also offer the service via HTTPS. Ideally, we should redirect the HTTP traffic to HTTPS.
* **Authentication/Authorization** We have not implemented any authentication/authorization system,
  so anyone who can reach the service would be able to use it, without any verification about if he/she has permissions to do so.
* **Non predictable ID for Journeys** We are currently using a number to identify the journeys.
  Using a predictalbe identifier implies security risks (e.g., anyone could drop of a journey, 
  even if he/she is not the _owner_ of the journey). Having a strong authentication/authorization
  mechanism would mitigate those risks.
* **Appropriate logging** We are currently logging writing to standard error.
  We should use a better logging system.
