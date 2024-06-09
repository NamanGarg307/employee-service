## employee-records-service
Golang Microservice for managing employees

## Endpoint Introductions 

### Create Employee (POST : /api/v1/employee)

Params Used - 
~~~
- Name
- Position
- Salary
~~~

This function does the following - 
- Creates new records for Employees.
- Creates Employees in bulk.

### Update Employee (PUT : /api/v1/employee)

Params Used - 
~~~
- ID - Unique ID For Employee.
- Name
- Position
- Salary
~~~

This function does the following -
- Updates Employee Record corresponding to given ID.

### Delete Employee (DELETE : /api/v1/employee/:id)

Params Used - 
~~~ 
- ID - Unique ID For Employee.
~~~

This function does the following -
- Deletes Record for Employee corresponding to given ID.

### Get Employees (GET : /api/v1/employee)

Params Used - 
~~~
- page - page user want see.
- per_page - number of records to be displayed per page.
~~~

This function does the following -
- Fetches Records for all Employees in a paginated format.


### Get Employees By ID (GET : /api/v1/employee/:id)

Params Used - 
~~~
- ID - Unique ID For Employee.
~~~

This function does the following -
- Fetches Records for Employees corresponding to given ID.
