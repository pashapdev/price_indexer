# Fair pricer

> Library for building a "fair price" based on the stream of price data from different sources as "index price"  
"Index price" is list of pair (timestamp, price)  
"Fair price" - is median of the price  streams

"Index price" examples
```
1577836800, 100.1
1577836860, 102
```

<br/>

## Median calculation algorithm
we will use 2 heaps:  
minHeap (minimal element in root)  
maxHeap (maximal element in root)  

Step 1: Add next item to one of the heaps  
if next item is smaller than maxHeap root add it to maxHeap,
else add it to minHeap

Step 2: Balance the heaps (after this step heaps will be either balanced or
one of them will contain 1 more item)

if number of elements in one of the heaps is greater than the other by
more than 1, remove the root element from the one containing more elements and add to the other one

Step 3:  
If the heaps contain equal amount of elements  
median = (root of maxHeap + root of minHeap)/2  
Else  
median = root of the heap with more elements


## Merge price streams algorithm
Gorutine is created to get data from price stream  
Data is read from this channel  
Error from data is checked structure  
Price data type is checked  
Data is writeen in output channel  


## Building "Index price" algorithm
Price is read from aggregating price stream  
Timestamp of price is being checked  
Price is added to calculator  
When the current minute ends, the current median value is taken from the calculator  
Calculator is being cleared  


## API documentation

**Median**

> Calculate median value

- New() - create new instance
- Add(val decimal) - add val into calculator
- Clear() - clear values
- Calculate() - calculate current median value

**Aggregator**

> Aggregator of price streams

- New(ticker entities.Ticker) - create new instance
- AggregatedChannel() - get output channel
- Add(ctx context.Context, subscribers ...priceStreamSubscriber) - add input price streams
  - priceStreamSubscriber - is interface with only one function **SubscribePriceStream() chan entities.TickerPriceWithErr**  


TickePriceWithError
```
type TickerPrice struct {
	Ticker Ticker
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
}

type TickerPriceWithErr struct {
	TickerPrice
	Err error
}
```

fill Err field if you want to send an error  
fill TickerPrice field if you want to send data
- Close() - close all subscriptions and dispose data

**Indexer**

> Builder of "Index price"

- New(source, calculator, saver) - create new instance
  - source is input price data
  - calculator is entity to calculate "fair price"
  - saver is entity to save result
- Run() - start processing


## Examples
in command line execute
```
make run-example
```


**Endpoint to get list of subscribers**
```
curl --location --request GET 'http://localhost:8080/api/v1/subscribers'
```

**Endpoint to get delete subscriber by id**
```
curl --location --request DELETE 'http://localhost:8080/api/v1/subscriber?id=<id>'
```

**Endpoint to create new subscriber**
```
curl --location --request POST 'http://localhost:8080/api/v1/subscriber' \
--header 'Content-Type: application/json' \
--data-raw '{
    "index": <unique identificator of subscriber>,
    "price": <the value that is sent once per minute>
}'
```

**First output**
```
2022/04/18 03:28:00 1650241640, no data to make decision
```

**Create streams with prices 1,2,3,4,5**
```
2022/04/18 03:28:00 1650241640, no data to make decision
2022/04/18 03:29:00 1650241680, no data to make decision
2022/04/18 03:30:00 1650241740, 3
```

**Create streams with price 100000**
```
2022/04/18 03:28:00 1650241640, no data to make decision
2022/04/18 03:29:00 1650241680, no data to make decision
2022/04/18 03:30:00 1650241740, 3
2022/04/18 03:31:00 1650241800, 3
2022/04/18 03:32:00 1650241860, 3.5
```

**Look at your streams**
```
[
    {
        "index": 1,
        "price": "1"
    },
    {
        "index": 2,
        "price": "2"
    },
    {
        "index": 3,
        "price": "3"
    },
    {
        "index": 4,
        "price": "4"
    },
    {
        "index": 5,
        "price": "5"
    },
    {
        "index": 6,
        "price": "100000"
    }
]
```

**Remove stream with id=3**
```
2022/04/18 03:27:40 reading config...
2022/04/18 03:27:40 starting server...
2022/04/18 03:28:00 1650241640, no data to make decision
2022/04/18 03:29:00 1650241680, no data to make decision
2022/04/18 03:30:00 1650241740, 3
2022/04/18 03:31:00 1650241800, 3
2022/04/18 03:32:00 1650241860, 3.5
2022/04/18 03:33:00 1650241920, 3.5
2022/04/18 03:34:00 1650241980, 3.5
2022/04/18 03:35:00 1650242040, 4
```
