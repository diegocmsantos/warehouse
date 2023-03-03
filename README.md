# Code Test

## Decisions

### No frameworks
I have been work with some Go web "frameworks" like Gorilla Mux and Gin.  
But for this test I decided to not use them just to play with the guts of http package.  
Then, I learned to value them because although they are really simple they help a lot.  

### Loosely coupling
The Warehouse ElementGetter interface make it easy to change beetween data sources.  
For now, it is just csv, but anything else that implements the GetElements function will be  
possible to use as data source.  

Also, the csv service receives a list of validators so we can add more validations to it without change the csv service.  

## The Challenge
What nice problem to solve.  
I started thinking that a linked list would solve the problem.  
I built it but I couldn't print the hierarchy as expected.  
Then, I went for a tree. Same problem.  
Then I used a raw approach with map of maps. It worked, but it was ugly and seemed fragile.  
Then, after going through all this I realized a map-linked-listsh would solve the problem.  
Really nice!!!

## Tests
Main keys:
- Good test coverage
- Table drive tests
- API tests using httptest package

I didn't grasp the requirement of having the tests into the io_test.go.  
So, I have create them in the places I think is the correct place and copy them to the required file.  

## Tradeoffs
I believe the biggest tradeoff I have made in this test was where to put the business logic about creating the list and validating things.  
For correctness, I believe they should be in the warehouse service, but thinking of performance I have put them in the csv, so I don't have to go over the list twice.  
In real world situation, I would go the opposite way first, then if after going to production and check if we have performance issues I would change for this approach.  
I value correctness and simplicity first then premature performance.  
About the csv header can be "out of order" I am assuming the item can be in the head or tail of the line.  
So, I wrote a check and a function to fix the order if it is the case.  
If the item can come in any position it's just a matter of enhance the `fixLineOrder` function to find the position of the item and put it in the tail.  