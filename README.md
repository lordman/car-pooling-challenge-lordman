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

## Language choice

I have chosen Go as the programming language to perform this challenge. I have opted in for
Go instead of other languages like Java, PHP or Python for several reasons.

First and foremost, Go uses static linking combining all dependency libraries
and modules into one single binary file based on OS type and architecture.
This means once the backend application is compiled, it is possible to just upload
compiled binary into server and it will work, without installing any dependencies there.

Note this lack of dependencies is specially important in this challenge, because the app
is going to be deployed in a Docker container: if I have chosen another language like, e.g.,
Python, PHP or Java, I would had needed to install in the container the Python or PHP interpreters
or the Java runtime environment, respectively. That would imply adding complexity to the deployment
and also a worse use of the resources.

Besides, Go is a statically typed language. A statically typed language is one where
variable types are declared explicitly for the compiler so even trivial bugs are caught
really easily, while in a dynamically typed language type inference is implemented by
the interpreter. Hence, some bugs may remain, due to the interpreter interpreting something
incorrectly. Using a statically typed language eliminates these interpretation issues.

The third factor considered to choose Go is performance. Go is really fast and generally
its performance is similar to that of Java or C++. Of course, this would depend on
the specific application, the load, and so on. I would need to perform some benchmarks
in order to guarantee Go will perform better that other languages, but this would be
beyond the scope of this challenge. To reinforce this statement about performance,
serve this benchmark in which it is compared with other extended languages for programming
web servers as a reference:

[Server-side I/O Performance: Node vs. PHP vs. Java vs. Go](https://www.toptal.com/back-end/server-side-io-performance-node-php-java-go)

These are the three factors that have made me opt for using Go. I could do a more detailed
analysis, but that would be out of the scope of this challenge.

## Third party libraries

### Gorilla/Mux
Package _gorilla/mux_ implements a request router and dispatcher for matching
incoming requests to their respective handler.

The name mux stands for "HTTP request multiplexer". Like the standard ```http.ServeMux```,
```mux.Router``` matches incoming requests against a list of registered routes
and calls a handler for the route that matches the URL or other conditions.

I used _gorilla/mux_ because that's what I have seen being used in several places.
It's very capable and makes it easier to specify which methods should be allowed
for specific routes. A good thing about _gorilla/mux_ is that it's compatible with
the built-in ```http.Handler``` API.

More detailed information about this library can be found on its GitHub website:

[https://github.com/gorilla/mux]

### Go Playground Validator

Go Playground Validator implements value validations for structs and individual fields
based on tags.

Input validation prevents improperly formed data from entering an information system.
Because it is difficult to detect a malicious user who is trying to attack software,
applications should check and validate all input entered into a system.

I used Go Playground Validator because it is, as far as I know, the most famous
and has the most stars on GitHub among the existing libraries.

You can find more detailed information about this library on its GitHub website:

[https://github.com/go-playground/validator]

## To-do list

There are several features/improvements I do not have currently implemented,
but would be a nice-to-have. In fact, they would be a must if this were going
to be released in production:

* **Automated unit testing**\
  I have manually tested all the considered use cases and also have trusted in
  the acceptance tests to verify everything works as intended. Nevertheless,
  it is a good practice to write and run automated unit tests to ensure that
  every _unit_ meets its design and behaves as intended. This would allow
  to catch bugs very early on in the development cycle and minimize the cost
  of fixing them. Those tests should be verified before committing the code,
  being complementary to any integration and acceptance testing.
* **Database**\
  Currently I have used slices to store the information about the cars and the journeys.
  In order to have some persistence in the data (not limited only to the program
  execution life) I would need to use a database to recover the status in case
  of crashing.
* **HTTPS**\
  The REST API is only provided via HTTP. This implies the traffic is sent in clear,
  which makes it susceptible to eavesdropping, man-in-the-middle attacks, and so son.
  The service should also be offered via HTTPS (two open ports should be needed).
  Ideally, HTTP traffic should be redirected to HTTPS.
* **Authentication/Authorization**\
  No authentication/authorization system has been implemented, so anyone who can reach
  the service would be able to use it, without any verification about if he/she has
  permissions to do so.
* **Non predictable ID for Journeys**\
  We are currently using a number to identify the journeys. In fact, I have assumed
  they arrive ordered by id (note the use of ```sort.Search``` to find a specific journey).
  Using a predictable identifier implies security risks (e.g., anyone could drop of a journey, 
  even if he/she is not the _owner_ of the journey). Having a strong
  authentication/authorization mechanism would mitigate those risks. On the other hand,
  using a non predictable identifier could have some performance implications.
* **Appropriate logging**\
  The application is currently logging writing to standard error. A better logging system
  should be used.
* **Audit third party libraries**\
  Trusting a third party library based only on its reputation (as I have done
  in this challenge) implies inherently assuming security risks. Also, choosing
  a specific library could lead to performance issues if we want to scale.
  Before proceeding with any final implementation work, an audit of these libraries
  to discover potential problems should carried out.
