# Influx DB Hardware Requirement & Capacity
- System Requirement: serve 12,997,774,152 SMS/year
- Single Node: 8-core CPU, 32GB RAM, 1TB SSD.
- CPU: 5–10% at peak.
- RAM: 10–20 GB for indexes, rest for caching.
- Disk: 260 GB/year, no issue.
- No scaling needed for 13B/year.



# Clean Cache
go clean -cache 
