# SimpleGolangLoadBalancer
This a Simple Golang loadBalancer

## What is a Load Balancer?
A Load Balancer is a type of reverse proxy that distributes incoming requests among multiple servers. It acts as an intermediary between clients and servers, ensuring that the workload is evenly balanced across the server pool.

A Proxy is a software component that sits in front of the backend servers, processing incoming requests and distributing them among the servers to achieve improved performance.

Load Balancers are often integrated into proxies to monitor the state of server or service endpoints, such as http://localhost:8081.

## Types of Load Balancers:

Round Robin: Requests are distributed equally among the servers in a cyclical manner, ensuring a balanced workload distribution.

Least Connections: Requests are routed to the server with the fewest active connections, allowing for optimal resource utilization.
IP Hash: This load balancing method identifies the client's IP address and uses hashing to determine which server should handle the request. It ensures that requests from the same client IP are consistently directed to the same server.

For further information, you can visit Wikipedia to explore and learn more about Load Balancers..
