# Software Requirement
The requirment of the system is to build a High Available and Scallable SMS Gateway System 

# Tech Stack
Backend: Gin, GO
DB: PostGreSQL, Influxdb
Frontend: Next.js

# The data flow of the system will be like,
- Data will be sent from DFS system and will be received by the SMS GW API
- Upon receiving, the data will be placed to a clustered rabbitmq queue(quoram) and will be saved to influxdb as well with a status of pending.
- There will be several consumers running per queue. After processing a message from any queue, the influxdb record (status) will also be updated to success.

# Service List
- SMSGW Service
- Redis Service
- RabbitMQ Service
- Consumer Services
- Database Service

# Message Format
{"mno":"Robi","msg_id":"2025032102343877835","msisdn":"01814266295","status":"queued","text":"","type":"general"}

# Influx DB Hardware Requirement & Capacity
- System Requirement: serve 12,997,774,152 SMS/year
- Single Node: 8-core CPU, 32GB RAM, 1TB SSD.
- CPU: 5–10% at peak.
- RAM: 10–20 GB for indexes, rest for caching.
- Disk: 260 GB/year, no issue.
- No scaling needed for 13B/year.

# Clean Cache
go clean -cache


# Questions
- Purging Mechanism?
- 