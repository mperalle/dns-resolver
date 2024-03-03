# Basic DNS resolver

Learning more about DNS by building my own DNS resolver in Go.

Features: 
- Resolve hostname to IP address
    - Query root nameserver
    - Query TLD nameserver
    - Query nameserver 

- Start a dns resolver server listening for incoming dns queries and responding back with answer record

- Cache records which were already queried for with expiration time

- Remove expired records based on TTL provided


Possible improvement: 
- Handle other types of record and not just A records
- Return all the A records available and not just one
- Query a second time for errors occuring when querying a dns server (timeout)
- Try another nameserver when one is not working

