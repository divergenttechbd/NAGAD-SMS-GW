# Swagger
http://localhost:8090/swagger/index.html

# RabbitMQ
http://192.168.7.172:15672

# Search Influx DB by msg_id
from(bucket: "dtl-bucket")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["_measurement"] == "sms_delivery")
  |> filter(fn: (r) => r.msg_id == "2025040614145198531")

	rm -rf $(go env GOCACHE) $(go env GOMODCACHE)