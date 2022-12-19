# Task Description

Our goal is to propose the right partner (a craftsman) to a customer based
on their project requirements. Matching of customers and partners is a crucial part in our
product. It determines how happy our customers will be with our partners and our partners
with the quality of the customer we connect them with.

The last product category that we reworked was flooring. The goal is to propose the right
partner based on the details of a customer's flooring project.

Your Task

Your task is to write an API that offers the following functionality:

● Based on a customer-request, return a list of partners that offer the service. The list
should be sorted by “best match”. The quality of the match is determined first on
average rating and second by distance to the customer.

● For a specific partner, return the detailed partner data.
Matching a customer and partner should happen on the following criteria:

● The partner should be experienced with the materials the customer requests for the
project.

● The customer is within the operating radius of the partner.The data in the request from the customer is:

   ● Material for the floor (wood, carpet, tiles)
   
   ● Address (assume that this is the lat/long of the house)
   
   ● Square meters of the floor
   
   ● Phone number (for the partner to contact the customer)

The structure of the partner data is as follows:

● Experienced in flooring materials (wood, carpet, tiles or any combination)

● Address (assume that this is the lat/long of the office)

● Operating radius (consider the beeline from the address)

● Rating (for this assignment you can assume that you already have a rating for a
partner)


## Prerequisites

You should have the following go version installed as the project was developed using the following version.

* go version go1.14 darwin/amd64
* Postman

## Things to Note

I tried to import some seed data via script in mysql database but as there were some problems faced so decided to
go via json file "`partners.json`" placed in main project directory that contains some partners sample data to run and
test the two required API endpoints.

As it was my first ever project in Go so please do expect certain bad coding styles or slow approaches.
The file `main.go` contains all the relevant code for handling the requests.

Some Resources I used to get idea from, to write the code:

https://tutorialedge.net/golang/parsing-json-with-golang/

https://andela.com/insights/using-golang-to-create-a-restful-json-api/



## How to run

Open the terminal and change to location where you want to clone the project. Clone the project using ssh via command:

`git clone git@github.com:mubashirkamran4/Mubashir_AroundHome_Challenge.git`

Go to project directory:

`cd Mubashir_AroundHome_Challenge`

Fire up the go server via:

`go run main.go`

This should start the server listening on 8080.

### Get Matching Partners based on Customer Request

Open Postman, select `GET` request in the Request mehtod and Enter the following URL:

[http://localhost:8080/match_partners](http://localhost:8080/match_partners)

Under the Request URL, Select `Body` tab and click "`raw`" and select `JSON` in the drop down for the datatype:
Paste the following Parmeters:

```
{
   "Clatitude": 24.5,
   "Material": "wood",
   "Clatitude": 24.5,
   "Clongitude": 21.5,
   "PhoneNumber": "+49122225869",
   "SquareMetres": 12.3
}
```
Hit Send and you should get the following repsonse as JSON:

```
[
   {
      "name":"partner4",
      "experienced_in":"wood,carpet",
      "latitude":68.9,
      "longitude":2.2,
      "operating_radius_latitude":125.5,
      "operating_radius_longitude":55.5,
      "rating":4
   },
   {
      "name":"partner3",
      "experienced_in":"tiles,wood",
      "latitude":67.9,
      "longitude":1.6,
      "operating_radius_latitude":105.5,
      "operating_radius_longitude":28.5,
      "rating":3
   }
]
```
We should be able to see that only the partners with "wood" in their experiences list and customer being within their operating
radius are displayed. For calculating the operating radius, I am assuming the following formula:
  
 `(Partners's Latitude and Longitude)` `-` `(Customer's Latitude and Longitude)` should be <=
   ` Operating Radius Latitude and Longitude` of Partner.


The partners are first sorted based on rating and then based on their distance from customer's address. For sorting based
on distance I am assuming the following formula:

`(Partners's Latitude and Longitude)` `-` `(Customer's Latitude and Longitude)`. The least values of both latitude and longitude
in this difference should come first in the list.



### Get Specific Partner based on Name

Open Postman, select `GET` request in the Request mehtod and Enter the following URL:

[http://localhost:8080/get_partner/partner3](http://localhost:8080/get_partner/partner3)

Here `partner3` implies the name of any partner from the partners list we got as a response to previous request.

Hit Enter and you should get the following response mentioning the details of specific partner:
```
{
    "name":"partner3",
    "experienced_in":"tiles,wood",
    "latitude":67.9,
    "longitude":1.6,
    "operating_radius_latitude":105.5,
    "operating_radius_longitude":28.5,
    "rating":3
}
```


